package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/lithammer/dedent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/gitops"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/libgitops/pkg/gitdir"
)

const defaultKnownHostsPath = "~/.ssh/known_hosts"

type gitOpsFlags struct {
	branch   string
	interval time.Duration
	timeout  time.Duration

	identityFile string
	hostsFile    string
	username     string
	password     string

	paths []string
}

// NewCmdGitOps runs the GitOps functionality of Ignite
func NewCmdGitOps(out io.Writer) *cobra.Command {
	f := &gitOpsFlags{
		branch:   "master",
		interval: 30 * time.Second,
		timeout:  1 * time.Minute,

		identityFile: "",
		hostsFile:    defaultKnownHostsPath,
		username:     "",
		password:     "",

		//paths:        []string{},
	}
	cmd := &cobra.Command{
		Use:   "gitops <repo-url>",
		Short: "Run the GitOps feature of Ignite",
		Long: dedent.Dedent(`
			Run Ignite in GitOps mode watching the given repository. The repository needs
			to be publicly cloneable. Ignite will watch for changes in the master branch
			by default, overridable with the branch flag (-b, --branch). If any new/changed
			VM specification files are found in the repo (in JSON/YAML format), their
			configuration will automatically be declaratively applied.

			To quit GitOps mode, use (Ctrl + C).
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			opts := gitdir.GitDirectoryOptions{
				Branch:   f.branch,
				Interval: f.interval,
				Timeout:  f.timeout,
			}
			if f.identityFile != "" {
				var err error
				// support ~ prefixes in the path
				f.identityFile, err = resolveHomeEnv(f.identityFile)
				log.Tracef("Parsed identity file path: %s", f.identityFile)
				util.GenericCheckErr(err)

				opts.IdentityFileContent, err = ioutil.ReadFile(f.identityFile)
				util.GenericCheckErr(err)
			}
			if f.hostsFile != "" {
				var err error
				// support ~ prefixes in the path
				f.hostsFile, err = resolveHomeEnv(f.hostsFile)
				log.Tracef("Parsed_known hosts file path: %s", f.hostsFile)
				util.GenericCheckErr(err)

				opts.KnownHostsFileContent, err = ioutil.ReadFile(f.hostsFile)
				util.GenericCheckErr(err)
			}
			if f.username != "" {
				opts.Username = &f.username
			}
			if f.password != "" {
				opts.Password = &f.password
			}

			util.GenericCheckErr(gitops.RunGitOps(args[0], opts))
		},
	}

	addGitOpsFlags(cmd.Flags(), f)
	return cmd
}

func addGitOpsFlags(fs *pflag.FlagSet, f *gitOpsFlags) {
	fs.StringVarP(&f.branch, "branch", "b", f.branch, "What branch to sync")
	fs.DurationVar(&f.interval, "interval", f.interval, "Sync interval for pushing to and pulling from the remote")
	fs.DurationVar(&f.timeout, "timeout", f.timeout, "Git operation (clone, push, pull) timeout")

	fs.StringVar(&f.identityFile, "identity-file", f.identityFile, "What SSH identity file to use for pushing")
	fs.StringVar(&f.hostsFile, "hosts-file", f.hostsFile, "What known_hosts file to use for remote verification")
	fs.StringVar(&f.username, "https-username", f.username, "What username to use when authenticating with Git over HTTPS")
	fs.StringVar(&f.password, "https-password", f.password, "What password/access token to use when authenticating with Git over HTTPS")

	// TODO: We need to add path prefix support to the WatchStorage to support this
	// fs.StringSliceVarP(&f.paths, "paths", "p", f.paths, "What subdirectories to care about. Default the whole repository")
}

// resolveHomeEnv resolves the ~ character in file paths at runtime
func resolveHomeEnv(path string) (string, error) {
	// exit quickly if possible
	if !strings.Contains(path, "~") {
		return path, nil
	}
	// get the HOME env var and fail if not found
	homeEnv := os.Getenv("HOME")
	if len(homeEnv) == 0 {
		return "", fmt.Errorf("Couldn't lookup home directory for user")
	}
	// replace the first tilde in the path with $HOME
	return strings.Replace(path, "~", homeEnv, 1), nil
}
