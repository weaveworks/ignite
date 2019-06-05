package metadata

import "io/ioutil"

// This interface enables taking in any struct which embeds Metadata
type AnyMetadata interface {
	GetMD() *Metadata
}

// Verify that Metadata implements AnyMetadata
var _ AnyMetadata = &Metadata{}

func (md *Metadata) GetMD() *Metadata {
	return md
}

func LoadAllMetadata(path string, loadFunc func(string) (AnyMetadata, error)) ([]AnyMetadata, error) {
	var mds []AnyMetadata

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			md, err := loadFunc(entry.Name())
			if err != nil {
				return nil, err
			}

			mds = append(mds, md)
		}
	}

	return mds, nil
}
