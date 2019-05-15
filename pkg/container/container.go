package container

//import (
//	"bytes"
//	"github.com/luxas/ignite/pkg/constants"
//	"github.com/luxas/ignite/pkg/image"
//	"github.com/luxas/ignite/pkg/util"
//	"html/template"
//	"strings"
//)
//
//const dockerfileTpl = `
//FROM alpine:latest
//
//VOLUME {{.DataDir}}
//
//CP {{.Image}} /
//CP {{.Firecracker}} /
//
//ENTRYPOINT ["/ignite", "container", "{{.Image}}"]
//`
//
//type dockerImageOptions struct {
//	Image       string
//	Firecracker string
//	DataDir     string
//}
//
//func ExportToDocker(i *build.Image) error {
//	tpl, _ := template.New("dockerfile").Parse(dockerfileTpl)
//
//	var tplBuffer bytes.Buffer
//	if err := tpl.Execute(&tplBuffer, &dockerImageOptions{
//		Image: i.Path,
//		Firecracker: "/usr/bin/firecracker", // TODO: Temporary
//		DataDir: constants.DATA_DIR,
//	}); err != nil {
//		return err
//	}
//
//	args := []string{"build", "-f-"}
//
//	if _, err := util.ExecuteCommand("docker", append(args, strings.Split(tplBuffer.String(), "\n")...)...); err != nil {
//		return err
//	}
//
//	return nil
//}
