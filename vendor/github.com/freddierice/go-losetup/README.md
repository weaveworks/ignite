# go-losetup
A losetup implementation for go-lang

## how to get it
```
go get gopkg.in/freddierice/go-losetup.v1
```

## example usage
```go
// attach a raw file to a loop device
dev, err := losetup.Attach("rawfile.img", 0, false)
if err != nil {
	// error checking
}

fmt.Printf("attached rawfile.img to %v\n", dev.Path())

err := dev.Detach()
if err != nil {
	// error checking
}
```

