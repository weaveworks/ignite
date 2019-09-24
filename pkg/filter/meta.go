package filter

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
)

const (
	filterSeparator   = ","
	filterApplyFailed = "failed to apply filtering"
	regexString       = `^(?P<key>{{(?:\.|[a-zA-Z]+)+}})=(?P<value>[a-zA-Z0-9-_]+)$`
)

type metaFilter struct {
	identifier    string
	expectedValue string
}

func (mf metaFilter) isExpected(object *runtime.ObjectMeta) (bool, error) {
	w := &bytes.Buffer{}
	tm, err := template.New("generic-filtering-vm").Parse(mf.identifier)
	if err != nil {
		return false, fmt.Errorf("failed to configure filtering with following template: %s", mf.identifier)
	}
	err = tm.Execute(w, object)
	if err != nil {
		return false, fmt.Errorf("failed to apply filtering on VM metadata")
	}
	res := w.String()
	if res != mf.expectedValue {
		return false, nil
	}
	return true, nil
}

// MultipleMetaFilter stores multiples metaFilter rule
type MultipleMetaFilter struct {
	filters []metaFilter
}

// AreExpected checks fileting rules are expected, an AND logical condition is applid between the underlying filters
func (mmf *MultipleMetaFilter) AreExpected(object *runtime.ObjectMeta) (bool, error) {
	for _, mf := range mmf.filters {
		res, err := mf.isExpected(object)
		if err != nil {
			return false, err
		} else if !res {
			return false, nil
		}
	}
	return true, nil
}

// extractKeyValueFiltering extracts the key to search for and the expected value form a string
func extractKeyValueFiltering(str string) (string, string, error) {
	reg, err := regexp.Compile(regexString)
	if err != nil {
		return "", "", err
	}
	matches := reg.FindAllStringSubmatch(str, -1)
	if len(matches) != 1 {
		return "", "", fmt.Errorf("failed to generate filter")
	}
	match := matches[0]
	if len(match) != 3 {
		return "", "", fmt.Errorf("failed to generate filter")
	}
	return match[1], match[2], nil
}

// extractMultipleKeyValueFiltering extracts all he keys and values to filter
func extractMultipleKeyValueFiltering(f string) ([]map[string]string, error) {
	filterList := strings.Split(f, filterSeparator)
	captureList := make([]map[string]string, 0, len(filterList))
	for _, filter := range filterList {
		key, value, err := extractKeyValueFiltering(filter)
		if err != nil {
			return nil, err
		}
		captureList = append(captureList, map[string]string{"key": key, "value": value})
	}
	return captureList, nil
}

// GenerateMultipleMetadataFiltering extract filterings and generates MultipleMetadataFiltering
func GenerateMultipleMetadataFiltering(str string) (*MultipleMetaFilter, error) {
	filtersInfos, err := extractMultipleKeyValueFiltering(str)
	if err != nil {
		return nil, err
	}
	metaFilterList := make([]metaFilter, 0, len(filtersInfos))
	for _, fInfo := range filtersInfos {
		metaFilterList = append(metaFilterList, metaFilter{identifier: fInfo["key"], expectedValue: fInfo["value"]})
	}
	return &MultipleMetaFilter{
		filters: metaFilterList,
	}, nil
}
