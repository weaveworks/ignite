package storage

import (
	"fmt"
	"strings"

	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type Cache interface {
	ListMeta() ignitemeta.APITypeList
	Has(ref string) bool
	MatchOne(ref string) (*ignitemeta.APIType, error)
	MatchMany(refs ...string) ignitemeta.APITypeList
}

func NewCache(list ignitemeta.APITypeList) Cache {
	byID := map[string]*ignitemeta.APIType{}
	byName := map[string]*ignitemeta.APIType{}
	for _, item := range list {
		byID[string(item.UID)] = item
		byName[string(item.Name)] = item
	}
	return &cache{
		list:   list,
		byID:   byID,
		byName: byName,
	}
}

type cache struct {
	list   ignitemeta.APITypeList
	byID   map[string]*ignitemeta.APIType
	byName map[string]*ignitemeta.APIType
}

func (c *cache) ListMeta() ignitemeta.APITypeList {
	return c.list
}

func (c *cache) byIDOrName(ref string) *ignitemeta.APIType {
	if _, ok := c.byID[ref]; ok {
		return c.byID[ref]
	}
	if _, ok := c.byName[ref]; ok {
		return c.byName[ref]
	}
	return nil
}

func (c *cache) Has(ref string) bool {
	return c.byIDOrName(ref) != nil
}

func (c *cache) prefixFilter(ref string) *[]string {
	matches := []string{}
	for _, item := range c.list {
		if strings.HasPrefix(string(item.UID), ref) {
			matches = append(matches, string(item.UID))
			continue
		}
	}
	return &matches
}

func (c *cache) MatchOne(ref string) (*ignitemeta.APIType, error) {
	if match := c.byIDOrName(ref); match != nil {
		return match, nil
	}

	matches := c.prefixFilter(ref)
	if len(*matches) == 1 {
		return c.byID[(*matches)[0]], nil
	}
	if len(*matches) > 1 {
		return nil, fmt.Errorf("multiple matches: %v", *matches)
	}
	return nil, fmt.Errorf("no matches")
}

func (c *cache) MatchMany(refs ...string) ignitemeta.APITypeList {
	result := ignitemeta.APITypeList{}
	for _, ref := range refs {
		for _, match := range *c.prefixFilter(ref) {
			result = append(result, c.byIDOrName(match))
		}
	}
	return result
}
