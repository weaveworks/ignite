package client

import (
	"testing"
)

var c = NewClient("/tmp/foo")

func TestClientGet(t *testing.T) {
	vm, err := c.VMs().Get("1234")
	t.Error(vm, err)

	t.Error(c.VMs().Delete("1234"))

	vms, err := c.VMs().List()
	t.Error(len(vms), err)
}

func TestClientDynamic(t *testing.T) {
	vm, err := c.Dynamic("VM").Get("12")
	t.Error(vm, err)

	vms, err := c.Dynamic("VM").List()
	t.Error(len(vms), err)
}
