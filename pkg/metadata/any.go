package metadata

import "io/ioutil"

func LoadAllMetadata(path string, loadFunc func(string) (Metadata, error)) ([]Metadata, error) {
	var mds []Metadata

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
