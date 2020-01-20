# Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# 	http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

# Set this to pass additional commandline flags to the go compiler, e.g. "make test EXTRAGOARGS=-v"
EXTRAGOARGS:=
DISABLE_ROOT_TESTS?=1

# The below files are needed and can be downloaded from the internet
testdata_objects = testdata/vmlinux testdata/root-drive.img testdata/firecracker

# --location is needed to follow redirects on github.com
curl = curl --location

all: build

test: all-tests

unit-tests: $(testdata_objects) check-kvm
	DISABLE_ROOT_TESTS=$(DISABLE_ROOT_TESTS) go test -short ./... $(EXTRAGOARGS)

all-tests: $(testdata_objects) check-kvm
	DISABLE_ROOT_TESTS=$(DISABLE_ROOT_TESTS) go test ./... $(EXTRAGOARGS)

check-kvm:
	@test -w /dev/kvm || \
		(echo "In order to run firecracker, $(shell whoami) must have write permission to /dev/kvm"; false)

generate build clean:
	go $@ $(EXTRAGOARGS)

distclean: clean
	rm -rf $(testdata_objects)

testdata/vmlinux:
	$(curl) -o $@ https://s3.amazonaws.com/spec.ccfc.min/img/hello/kernel/hello-vmlinux.bin

testdata/firecracker:
	$(curl) -o $@ https://github.com/firecracker-microvm/firecracker/releases/download/v0.18.0/firecracker-v0.18.0
	chmod +x $@

testdata/root-drive.img:
	$(curl) -o $@ https://s3.amazonaws.com/spec.ccfc.min/img/hello/fsfiles/hello-rootfs.ext4

.PHONY: all generate clean distclean build test unit-tests all-tests check-kvm
