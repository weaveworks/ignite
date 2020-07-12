package filter

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
)

const (
	filterSeparator = ","
	regexString     = `^(?P<key>{{(?:\.|[a-zA-Z0-9]+)+}})(?P<operator>(?:=|==|!=|=~|!~))(?P<value>[a-zA-Z0-9-_:/\.\s]+)$`
)

type metaFilter struct {
	identifier    string
	expectedValue string
	operator      string
}

func (mf metaFilter) isExpected(object *api.VM) (bool, error) {
	w := &bytes.Buffer{}
	tm, err := template.New("generic-filtering-vm").Parse(mf.identifier)
	if err != nil {
		return false, fmt.Errorf("failed to configure filtering with following template: %s", mf.identifier)
	}
	err = tm.Execute(w, object)
	if err != nil {
		return false, fmt.Errorf("failed to apply filtering on VM, the filter might be incorrect")
	}
	res := w.String()
	switch mf.operator {
	case "==":
		return mf.isEqual(res), nil
	case "=":
		return mf.isEqual(res), nil
	case "!=":
		return !mf.isEqual(res), nil
	case "=~":
		return mf.contains(res), nil
	case "!~":
		return !mf.contains(res), nil
	default:
		return false, fmt.Errorf("Unexpected operator: %s", mf.operator)
	}
}

func (mf metaFilter) isEqual(value string) bool {
	return mf.expectedValue == value
}

func (mf metaFilter) contains(value string) bool {
	return strings.Contains(value, mf.expectedValue)
}

// MultipleMetaFilter stores multiples metaFilter rule
type MultipleMetaFilter struct {
	filters []metaFilter
}

// AreExpected checks fileting rules are expected, an AND logical condition is applid between the underlying filters
func (mmf *MultipleMetaFilter) AreExpected(object *api.VM) (bool, error) {
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
func extractKeyValueFiltering(str string) (string, string, string, error) {
	reg, err := regexp.Compile(regexString)
	if err != nil {
		return "", "", "", err
	}
	matches := reg.FindAllStringSubmatch(str, -1)
	if len(matches) != 1 {
		return "", "", "", fmt.Errorf("failed to generate filter")
	}
	match := matches[0]
	if len(match) != 4 {
		return "", "", "", fmt.Errorf("failed to generate filter")
	}
	return match[1], match[3], match[2], nil
}

// extractMultipleKeyValueFiltering extracts all the keys and values to filter
func extractMultipleKeyValueFiltering(f string) ([]metaFilter, error) {
	filterList := strings.Split(f, filterSeparator)
	captureList := make([]metaFilter, 0, len(filterList))
	for _, filter := range filterList {
		key, value, op, err := extractKeyValueFiltering(filter)
		if err != nil {
			return nil, fmt.Errorf("failed to extract keys-values from filter list %s", filterList)
		}
		captureList = append(captureList, metaFilter{identifier: key, expectedValue: value, operator: op})
	}
	return captureList, nil
}

// GenerateMultipleMetadataFiltering extract filterings and generates MultipleMetadataFiltering
func GenerateMultipleMetadataFiltering(str string) (*MultipleMetaFilter, error) {
	metaFilterList, err := extractMultipleKeyValueFiltering(str)
	if err != nil {
		return nil, err
	}
	return &MultipleMetaFilter{
		filters: metaFilterList,
	}, nil
}
