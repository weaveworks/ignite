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
apiVersion: ignite.weave.works/v1alpha4
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

**NOTE:** `uid` must be set. VM configurations without uid are ignored.

For a more complete example repository configuration, see [luxas/ignite-gitops](https://github.com/luxas/ignite-gitops)

After you have [installed Ignite](installation.md), you can do the following:

```console
ignited gitops git@github.com:<user>/<repository>.git
```

**NOTE:** HTTPS doesn't preserve authentication information for `ignited` to push changes,
you need to set up SSH authentication and use the SSH clone URL for now.

Ignite will now search that repo for suitable JSON/YAML files, and apply their state locally.
You should see `my-vm` starting up in `ignite ps`. To enter the VM, run `ignite ssh my-vm`.

### Using a local git repo for testing

Create a new directory and initialize it as a bare git repo. For ignited to
push updates to the git repo, the repo must be created as a bare git repo.

```console
$ mkdir ~/ignite-gitops
$ cd ~/ignite-gitops
$ git init --bare
$ ls
HEAD  branches  config  description  hooks  info  objects  refs
```

Clone this git repo and add the sample VM configuration:

```console
$ git clone file:///home/user/ignite-gitops ignite-gitops-clone
Cloning into 'ignite-gitops-clone'...
warning: You appear to have cloned an empty repository.
$ cd ignite-gitops-clone
$ git remote -v
origin	file:///home/user/ignite-gitops (fetch)
origin	file:///home/user/ignite-gitops (push)
$ # Create my-vm.yaml, commit and push.
```

Run ignited against the bare git repo:

```console
$ sudo ignited gitops file:///home/user/ignite-gitops
INFO[0000] Starting GitOps loop for repo at "file:///home/user/ignite-gitops"
INFO[0000] Whenever changes are pushed to the target branch, Ignite will apply the desired state locally
INFO[0000] Initializing the Git repo...
INFO[0000] Running in read-write mode, will commit back current status to the repo
INFO[0000] Starting the commit loop...
INFO[0000] Starting the checkout loop...
INFO[0000] Starting to clone the repository file:///home/user/ignite-gitops with timeout 1m0s
INFO[0000] New commit observed on branch "master": 7095d3603649792f44110cc5cf7deb0ded897e4b. User initiated: true
INFO[0003] Creating VM "599615df99804ae8" with name "my-vm"...
INFO[0004] Starting VM "599615df99804ae8" with name "my-vm"...
INFO[0004] Networking is handled by "cni"
INFO[0004] Started Firecracker VM "599615df99804ae8" in a container with ID "ignite-599615df99804ae8"
...
INFO[0030] A new commit with the actual state has been created and pushed to the origin: "186a8eb6eced5843480d4a10872db071ecc76439"
INFO[0030] New commit observed on branch "master": 186a8eb6eced5843480d4a10872db071ecc76439. User initiated: false
...
```

**NOTE:** Root permission may be required to push updates to the main repo
because ignited pushes git updates as root.

Please refer to [docs/declarative-config.md](declarative-config.md) for the full API reference.
