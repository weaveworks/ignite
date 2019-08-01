# Developer documentation

Ignite is a Go project using well-known libraries like:

 - github.com/spf13/cobra
 - github.com/spf13/pflag
 - k8s.io/apimachinery
 - sigs.k8s.io/yaml
 - github.com/firecracker-microvm/firecracker-go-sdk

and so on.

It uses Go modules as the vendoring mechanism.

## Build from source

The only build requirement is Docker.

```console
make binary
make image
```

## Pre-commit tidying

Before committing, please run this make target to tidy your local environment:

```console
make tidy
```

## Building reference OS images

```console
make -C images WHAT=ubuntu
make -C images WHAT=centos
```
