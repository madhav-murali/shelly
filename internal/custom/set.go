package custom

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
