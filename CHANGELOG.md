# Changelog

## v0.3.0 (unreleased)

Major release with significant UX and internal improvements:

 - There is no longer a difference between an Ignite image and an OCI image, this is now the same thing. Ignite operates on OCI images
   directly, for both OS images and kernels.
 - Ignite now uses Devicemapper's `thin-provisioning` feature, for way better performance and added capabilities.
 - It is now possible to do `ignite run [OCI image]` directly, and everything (e.g. pulling the image) is handled automatically. e.g. `ignite run -i weaveworks/ignite-ubuntu`.
 - Now `ignite images` shows OCI images that are cached and ready to use, and `ignite kernels` the kernel images ready to use.
 - Added an example usage guide for running a Kubernetes cluster in HA mode using kubeadm and Ignite.
 - Removed `ignite build`, and `ignite image/kernel import`; as these are no longer needed
 - Importing an image from a tROADMAPar file is no longer possible, package the contents in an OCI image instead
 - Added a new command `ignite ssh [vm]` and flag: `ignite run --ssh`. This allows for automatic SSH logins.
 - Now Ignite logs user-friendly messages by default. To get machine-readable output, use the `--quiet` flag.
 - Ignite now requires the user to be `root`. This will be revisited later, when the architecture has changed.
 - The command outputs and structure is now more user-friendly.
 - Fixed several bugs both under the hood, and user-affecting ones

## v0.2.0

Major release with significant improvements

 - Ignite is now using `devicemapper` under the hood, for overlay snapshots for filesystem writes, allowing for image reuse, efficient use of space and way faster builds!
 - Added sample Ubuntu 18.04 and CentOS 7 OS images & a 4.19 kernel build
 - Automatic network configuration, now the OS image doesn't need to enable DHCP, as that is done in the kernel
 - Automatically populate `/etc/hosts` and `/etc/resolv.conf`, too
 - Add an option to bind a port exposed by the VM to a host port (`ignite run -p 80:80`)
 - Add an option for modifying the kernel command line (`ignite run --kernel-args`)
 - Add an option to copy files from the host into the VM (`ignite run --copy-files`)
 - Add an option to specify the amount of cores, RAM, and overlay size (`ignite run --cpus 2 --memory 1024 --size 4GB`)
 - Removed the need for the Ignite container to run with `--privileged`
 - Allow for force-deletions of images, kernels and vms.
 - Added documentation.
 - Moved repo from luxas/ignite to weaveworks/ignite

## v0.1.0

This is the first, proof-of-concept version of Ignite.
It has all the essential features, and a pretty complete implementation of the docker UX.
