package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

type data struct {
	v1alpha1.ImageData
}

func (i *data) GetID() string {
	return i.Spec.Source.Digest
}

func (i *data) Path() string {
	return path.Join(constants.IMAGE_DIR, i.GetID())
}

func (i *data) Save() error {
	f, err := os.Create(path.Join(i.Path(), constants.METADATA))
	if err != nil {
		return err
	}
	defer f.Close()

	y, err := json.MarshalIndent(&i, "", "    ")
	if err != nil {
		return err
	}

	if _, err := f.Write(append(y, '\n')); err != nil {
		return err
	}

	return nil
}

func (i *data) Load() error {
	//c := util.NewChainError()
	//defer c.Recover()

	if !util.DirExists(i.Path()) {
		return fmt.Errorf("nonexistent data: %v", i)
	}

	f := path.Join(i.Path(), constants.METADATA)

	if !util.FileExists(f) {
		return fmt.Errorf("metadata file missing for data: %v", i)
	}

	d, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}

	//return c.Result()
	return nil
}
