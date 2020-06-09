package cmd

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/gitops"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/libgitops/pkg/gitdir"
)

type gitOpsFlags struct {
	branch   string
	interval time.Duration
	timeout  time.Duration

	identityFile string
	username     string
	password     string

	paths     []string
	hostsFile string
}

// NewCmdGitOps runs the GitOps functionality of Ignite
func NewCmdGitOps(out io.Writer) *cobra.Command {
	f := &gitOpsFlags{
		branch:   "master",
		interval: 30 * time.Second,
		timeout:  1 * time.Minute,

		identityFile: "",
		username:     "",
		password:     "",

		//paths:        []string{},
		//hostsFile:    "~/.ssh/known_hosts",
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
				opts.IdentityFileContent, err = ioutil.ReadFile(f.identityFile)
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
	fs.StringVar(&f.username, "https-username", f.username, "What username to use when authenticating with Git over HTTPS")
	fs.StringVar(&f.password, "https-password", f.password, "What password/access token to use when authenticating with Git over HTTPS")

	// TODO: We need to add path prefix support to the WatchStorage to support this
	// fs.StringSliceVarP(&f.paths, "paths", "p", f.paths, "What subdirectories to care about. Default the whole repository")

	// TODO: When https://github.com/fluxcd/toolkit/issues/2 is fixed and the
	// https://github.com/fluxcd/source-controller/tree/master/internal/crypto/ssh/knownhosts package
	// can be vendored, we'll add support for reading a hosts file. In the meantime, the
	// SSH_KNOWN_HOSTS variable can be used.
	// fs.StringVar(&f.hostsFile, "hosts-file", f.hostsFile, "What hosts file to use")
}
