# Ignite - the GitOps VM

Ignite is a “GitOps-first” project, GitOps is supported out of the box using the `ignited gitops` command.
Previously this was integrated as `ignite gitops`, but this functionality has now moved to `ignited`,
Ignite's upcoming daemon binary.

In Git you declaratively store the desired state of a set of VMs you want to manage.
`ignited gitops` reconciles the state from Git, and applies the desired changes as state is updated in the repo.
It also commits and pushes any local changes/additions to the managed VMs back to the repository.

This can then be automated, tracked for correctness, and managed at scale - [just some of the benefits of GitOps](https://www.weave.works/technologies/gitops/).

The workflow is simply this:

- Run `ignited gitops [repo]`, where repo is an **SSH url** to your Git repo
- Create a file with the VM specification, specifying how much vCPUs, RAM, disk, etc. you’d like for the VM
- Run `git push` and see your VM start on the host

See it in action! (Note: The screencast is from an older version which differs somewhat)

[![asciicast](https://asciinema.org/a/255797.svg)](https://asciinema.org/a/255797)

## Try it out

Go ahead and create a Git repository.

**NOTE:** You need an SSH key for **root** that has push access to your repository. `ignited` will commit and push changes
back to the repository using the default key for it. To edit your root's git configuration, run
`sudo gitconfig --global --edit`. The root requirement will be removed in a future release.

 Here's a sample configuration you can push to it (my-vm.yaml):

```yaml
apiVersion: ignite.weave.works/v1alpha3
kind: VM
metadata:
  name: my-vm
  uid: 599615df99804ae8
spec:
  image:
    oci: weaveworks/ignite-ubuntu
  cpus: 2
  diskSize: 3GB
  memory: 800MB
  ssh: true
status:
  running: true
```

For a more complete example repository configuration, see [luxas/ignite-gitops](https://github.com/luxas/ignite-gitops)

After you have [installed Ignite](installation.md), you can do the following:

```console
ignited gitops git@github.com:<user>/<repository>.git
```

**NOTE:** HTTPS doesn't preserve authentication information for `ignited` to push changes,
you need to set up SSH authentication and use the SSH clone URL for now.

Ignite will now search that repo for suitable JSON/YAML files, and apply their state locally.
You should see `my-vm` starting up in `ignite ps`. To enter the VM, run `ignite ssh my-vm`.

Please refer to [docs/declarative-config.md](declarative-config.md) for the full API reference.
