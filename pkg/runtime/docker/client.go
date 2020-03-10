package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/preflight"
	"github.com/weaveworks/ignite/pkg/preflight/checkers"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	dcSocket = "/var/run/docker.sock"

	// streamLoadedImage is a substring of the stream logs from docker server
	// response when importing a local docker image from tar file.
	streamLoadedImage = "Loaded image"
)

// dockerClient is a runtime.Interface
// implementation serving the Docker client
type dockerClient struct {
	client *client.Client
}

// dockerStreamOutput is the structure of a docker stream output.
type dockerStreamOutput struct {
	Stream string `json:"stream"`
}

var _ runtime.Interface = &dockerClient{}

// GetDockerClient builds a client for talking to docker
func GetDockerClient() (*dockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.35"))
	if err != nil {
		return nil, err
	}

	return &dockerClient{
		client: cli,
	}, nil
}

func (dc *dockerClient) PullImage(image meta.OCIImageRef) (err error) {
	var rc io.ReadCloser
	if rc, err = dc.client.ImagePull(context.Background(), image.Normalized(), types.ImagePullOptions{}); err == nil {
		// Don't output the pull command
		util.DeferErr(&err, rc.Close)
		_, err = io.Copy(ioutil.Discard, rc)
	}

	return
}

// ImportImage imports a local image from a given tar file and returns a list of
// OCI image refs of all the imported images.
func (dc *dockerClient) ImportImage(imageFilePath string) ([]meta.OCIImageRef, error) {
	log.Debugf("docker: Importing image from %s", imageFilePath)
	r, err := os.Open(imageFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	resp, err := dc.client.ImageLoad(context.Background(), r, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imgRefs := []meta.OCIImageRef{}

	// Maybe open the tar file and read manifest.json to get the image refs?

	// Response body content looks like:
	// {"stream":"Loaded image: hello-world:latest\n"}
	// {"stream":"Loaded image: busybox:latest\n"}
	//
	// Parse the docker steam output and get the image refs.
	rd := bufio.NewScanner(resp.Body)
	for rd.Scan() {
		log.Debugf("docker: %s", rd.Text())

		// Parse the json output to get the stream message.
		var stream dockerStreamOutput
		if err := json.Unmarshal(rd.Bytes(), &stream); err != nil {
			return imgRefs, err
		}

		// Split the message at ":" and check if the first part contains "Loaded
		// image".
		message := strings.SplitN(stream.Stream, ":", 2)
		if strings.Contains(message[0], streamLoadedImage) {
			// Trim the second part of the string and form an OCI image ref.
			image, err := meta.NewOCIImageRef(strings.TrimSpace(message[1]))
			if err != nil {
				return imgRefs, err
			}
			imgRefs = append(imgRefs, image)
		}
	}

	return imgRefs, nil
}

func (dc *dockerClient) InspectImage(image meta.OCIImageRef) (*runtime.ImageInspectResult, error) {
	res, _, err := dc.client.ImageInspectWithRaw(context.Background(), image.Normalized())
	if err != nil {
		return nil, err
	}

	// By default parse the OCI content ID from the Docker image ID
	contentRef := res.ID
	if len(res.RepoDigests) > 0 {
		// If the image has Repo digests, use the first one of them to parse
		// the fully qualified OCI image name and digest. All of the digests
		// point to the same contents, so it doesn't matter which one we use
		// here. It will be translated to the right content by the runtime.
		contentRef = res.RepoDigests[0]
	}

	// Parse the OCI content ID based on the available reference
	id, err := meta.ParseOCIContentID(contentRef)
	if err != nil {
		return nil, err
	}

	r := &runtime.ImageInspectResult{
		ID:   id,
		Size: res.Size,
	}

	return r, nil
}

func (dc *dockerClient) ExportImage(image meta.OCIImageRef) (r io.ReadCloser, cleanup func() error, err error) {
	config, err := dc.client.ContainerCreate(context.Background(), &container.Config{
		Cmd:   []string{"sh"}, // We need a temporary command, this doesn't need to exist in the image
		Image: image.Normalized(),
	}, nil, nil, "")
	if err != nil {
		return
	}

	if r, err = dc.client.ContainerExport(context.Background(), config.ID); err == nil {
		cleanup = func() error { return dc.RemoveContainer(config.ID) }
	}

	return
}

func (dc *dockerClient) InspectContainer(container string) (*runtime.ContainerInspectResult, error) {
	res, _, err := dc.client.ContainerInspectWithRaw(context.Background(), container, false)
	if err != nil {
		return nil, err
	}

	return &runtime.ContainerInspectResult{
		ID:        res.ID,
		Image:     res.Image,
		Status:    res.State.Status,
		IPAddress: net.ParseIP(res.NetworkSettings.IPAddress),
		PID:       uint32(res.State.Pid),
	}, nil
}

func (dc *dockerClient) AttachContainer(container string) (err error) {
	// TODO: Rework to perform the attach via the Docker client,
	// this will require manual TTY and signal emulation/handling.
	// Implement the pseudo-TTY and remove this call, see
	// https://github.com/weaveworks/ignite/pull/211#issuecomment-512809841
	var ec int
	if ec, err = util.ExecForeground("docker", "attach", container); err != nil {
		if ec == 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			err = nil // Don't treat it as an error
		}
	}

	return
}

func (dc *dockerClient) RunContainer(image meta.OCIImageRef, config *runtime.ContainerConfig, name, id string) (string, error) {
	binds := make([]string, 0, len(config.Binds))
	for _, bind := range config.Binds {
		binds = append(binds, fmt.Sprintf("%s:%s", bind.HostPath, bind.ContainerPath))
	}

	devices := make([]container.DeviceMapping, 0, len(config.Devices))
	for _, device := range config.Devices {
		devices = append(devices, container.DeviceMapping{
			PathOnHost:        device.HostPath,
			PathInContainer:   device.ContainerPath,
			CgroupPermissions: "rwm",
		})
	}

	stopTimeout := int(config.StopTimeout)
	bindings, exposed := portBindingsToDocker(config.PortBindings)

	c, err := dc.client.ContainerCreate(context.Background(), &container.Config{
		Hostname:     config.Hostname,
		ExposedPorts: exposed,
		Tty:          true, // --tty
		OpenStdin:    true, // --interactive
		Cmd:          config.Cmd,
		Image:        image.Normalized(),
		Labels:       config.Labels,
		StopTimeout:  &stopTimeout,
	}, &container.HostConfig{
		Binds:        binds,
		NetworkMode:  container.NetworkMode(config.NetworkMode),
		PortBindings: bindings,
		AutoRemove:   config.AutoRemove,
		CapAdd:       config.CapAdds,
		Resources: container.Resources{
			Devices: devices,
		},
	}, nil, name)
	if err != nil {
		return "", err
	}

	return c.ID, dc.client.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{})
}

func (dc *dockerClient) StopContainer(container string, timeout *time.Duration) error {
	if err := dc.client.ContainerStop(context.Background(), container, timeout); err != nil {
		return err
	}

	// Wait for the container to be stopped
	return dc.waitForContainer(container, cont.WaitConditionNotRunning, nil)
}

func (dc *dockerClient) KillContainer(container, signal string) error {
	if err := dc.client.ContainerKill(context.Background(), container, signal); err != nil {
		return err
	}

	// Wait for the container to be killed
	return dc.waitForContainer(container, cont.WaitConditionNotRunning, nil)
}

func (dc *dockerClient) RemoveContainer(container string) error {
	// Container waiting can fail if the container gets removed before
	// we reach the waiting fence. Start the waiter in a goroutine before
	// performing the removal to ensure proper removal detection.
	errC := make(chan error)
	readyC := make(chan struct{})
	go func() {
		errC <- dc.waitForContainer(container, cont.WaitConditionRemoved, &readyC)
	}()

	<-readyC // The ready channel is used to wait until removal detection has started
	if err := dc.client.ContainerRemove(context.Background(), container, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	// Wait for the container to be removed
	return <-errC
}

func (dc *dockerClient) ContainerLogs(container string) (io.ReadCloser, error) {
	return dc.client.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
		ShowStdout: true, // We only need stdout, as TTY mode merges stderr into it
	})
}

func (cc *dockerClient) Name() runtime.Name {
	return runtime.RuntimeDocker
}

func (dc *dockerClient) RawClient() interface{} {
	return dc.client
}

func (dc *dockerClient) PreflightChecker() preflight.Checker {
	return checkers.NewExistingFileChecker(dcSocket)
}

func (dc *dockerClient) waitForContainer(container string, condition cont.WaitCondition, readyC *chan struct{}) error {
	resultC, errC := dc.client.ContainerWait(context.Background(), container, condition)

	// The ready channel can be used to wait until
	// the container wait request has been sent to
	// Docker to ensure proper ordering
	if readyC != nil {
		*readyC <- struct{}{}
	}

	select {
	case result := <-resultC:
		if result.Error != nil {
			return fmt.Errorf("failed to wait for container %q: %s", container, result.Error.Message)
		}
	case err := <-errC:
		return err
	}

	return nil
}
