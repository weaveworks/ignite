package metadata

import (
	"io/ioutil"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

func LoadAllMetadata(path string, loadFunc func(meta.UID) (Metadata, error)) ([]Metadata, error) {
	var mds []Metadata

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			md, err := loadFunc(meta.UID(entry.Name()))
			if err != nil {
				return nil, err
			}

			mds = append(mds, md)
		}
	}

	return mds, nil
}
