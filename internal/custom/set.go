package custom

import (
	"maps"
	"slices"
	"strings"
)

type CmdSet struct {
	items map[string]struct{}
}

func NewCmdSet() *CmdSet {
	return &CmdSet{
		items: make(map[string]struct{}),
	}
}

func (c *CmdSet) Find(key string) bool {
	_, exists := c.items[key]
	return exists
}

func (c *CmdSet) Add(key string) {
	c.items[key] = struct{}{}
}

func (c *CmdSet) Delete(key string) {
	delete(c.items, key)
}

func (c *CmdSet) ReturnAll() []string {
	return slices.Collect(maps.Keys(c.items))
}

func (c *CmdSet) ReturnAllLower() []string {
	var res []string
	res = c.ReturnAll()
	for i, word := range res {
		res[i] = strings.ToLower(word)
	}
	return res
}
