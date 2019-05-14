## Firecracker Ignite

Lets you start a micro-VM with Firecracker easily.
Integrates well with cloud-native projects like CNI, containerd and Docker.

### How to use

```console
$ ignite build luxas/ubuntu:18.10 my-vm-image
$ ignite images
$ ignite start my-vm-image
$ ignite attach my-vm-image
$ ignite ps
```

### Maintainers

- Lucas Käldström, @luxas
- Dennis Marttinen, @twelho

### License

Apache 2.0
