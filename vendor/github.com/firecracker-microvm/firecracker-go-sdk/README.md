A basic Go interface to the Firecracker API
====

[![Build status](https://badge.buildkite.com/de08ca676829bedbf6de040c2c2ba1a5d2892e220997c2abdd.svg)](https://buildkite.com/firecracker-microvm/firecracker-go-sdk)

This package is a Go library to interact with the Firecracker API. It
is designed as an abstraction of the OpenAPI-generated client that
allows for convenient manipulation of Firecracker VM from Go programs.

There are some Firecracker features that are not yet supported by the
SDK.  These are tracked as GitHub issues with the
[firecracker-feature](https://github.com/firecracker-microvm/firecracker-go-sdk/issues?q=is%3Aissue+is%3Aopen+label%3Afirecracker-feature)
label . Contributions to address missing features are welcomed.

Developing
---

Please see [HACKING](HACKING.md)

Building
---

This library requires Go 1.11 and Go modules to build.  A Makefile is provided
for convenience, but is not required.  When using the Makefile, you can pass
additional flags to the Go compiler via the `EXTRAGOARGS` make variable.

Tools
---

There's a [firectl](https://github.com/firecracker-microvm/firectl/)
tool that provides a simple command-line interface to launching a
firecracker VM. It also serves as an example client of this SDK.

Network configuration
---

Firecracker, by design, only supports Linux tap devices. The SDK
provides facilities to attach a tap device to the Firecracker VM, but
the client is responsible for further configuration.

License
====

This library is licensed under the Apache 2.0 License. 
