## Ignite - the GitOps VM

Ignite is a “GitOps-first” project, GitOps is supported out of the box using the `ignite gitops` command.

In Git you declaratively store the desired state of a set of VMs you want to manage.
`ignite gitops` reconciles the state from Git, and applies the desired changes as state is updated in the repo.

This can then be automated, tracked for correctness, and managed at scale - [just some of the benefits of GitOps](https://www.weave.works/technologies/gitops/).

The workflow is simply this:

 - Run `ignite gitops [repo]`, where repo points to your Git repo
 - Create a file with the VM specification, specifying how much vCPUs, RAM, disk, etc. you’d like from the VM
 - Run `git push` and see your VM start on the host

See it in action!

[![asciicast](https://asciinema.org/a/255797.svg)](https://asciinema.org/a/255797)

### Try it out

In this folder, there are two sample files [declaratively specifying how VMs should be run](../docs/declarative-config.md).
This means, you can try this feature out yourself!

After you have [installed Ignite](../docs/installation.md), you can do the following:

```console
$ ignite gitops https://github.com/luxas/ignite-gitops
```

Ignite will now search that repo for suitable JSON/YAML files, and apply their state locally.
(You can go and check the files out first, too, at: https://github.com/luxas/ignite-gitops)

To show how you could create your own repo, similar to `luxas/ignite-gitops`, refer to these two files:

 - [amazonlinux-vm.json](amazonlinux-vm.json)
 - [ubuntu-vm.yaml](ubuntu-vm.yaml)

You can use these files as the base for your own Git-managed VM-spawning flow.

Please refer to [docs/declarative-config.md](../docs/declarative-config.md) for the full API reference.
