ignite - CLI reference
======================

ignite: easily run Firecracker VMs

Synopsis
--------

Ignite is a containerized Firecracker microVM administration tool.
It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

Administration is divided into three subcommands:

  image       Manage base images for VMs
  kernel      Manage VM kernels
  vm          Manage VMs

Ignite also supports the same commands as the Docker CLI.
Combining an Image and a Kernel gives you a runnable VM.

Example usage:

.. code-block:: shell

  $ ignite run centos:7 \
    --cpus 2 \
    --memory 2GB \
    --ssh \
    --name my-vm
  $ ignite images
  $ ignite kernels
  $ ignite ps
  $ ignite logs my-vm
  $ ignite ssh my-vm


Options
-------

.. code-block:: shell

  -h, --help                 help for ignite
      --log-level loglevel   Specify the loglevel for the program (default info)
  -q, --quiet                The quiet mode allows for machine-parsable output, by printing only IDs


In this folder you can read more about how to use Ignite, and how it works:

.. toctree::
  :glob:
  :titlesonly:
  :maxdepth: 0

  ignite_attach
  ignite_completion
  ignite_create
  ignite_exec
  ignite_gitops
  ignite_image
  ignite_inspect
  ignite_kernel
  ignite_kill
  ignite_logs
  ignite_ps
  ignite_rm
  ignite_rmi
  ignite_rmk
  ignite_run
  ignite_ssh
  ignite_start
  ignite_stop
  ignite_version
  ignite_vm
  *
