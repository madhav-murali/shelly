package trie

type Node struct {
	isEnd bool
	child map[rune]*Node
}

func NewNode() *Node {
	return &Node{child: make(map[rune]*Node)}
}

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{
		root: &Node{child: make(map[rune]*Node)},
	}
}

func (t *Trie) Add(cmd string) {
	cur := t.root
	for _, c := range cmd {
		if cur.child[rune(c)] == nil {
			cur.child[rune(c)] = NewNode()
		}
		cur = cur.child[rune(c)]
	}
	cur.isEnd = true
}

func (t *Trie) AddAll(cmds []string) {
	for _, cmd := range cmds {
		t.Add(cmd)
	}
}

func (t *Trie) Find(cmd string, isPrefix bool) bool {
	cur := t.root
	for _, c := range cmd {
		if cur.child[rune(c)] == nil {
			return false
		}
		cur = cur.child[rune(c)]
	}
	if !isPrefix {
		return cur.isEnd
	}
	return true
}

func (t *Trie) IsPrefix(prefix string) bool {
	return t.Find(prefix, true)
}

// Returns the best match from the given prefix, if a prefix is not found returns false
func (t *Trie) HasPrefix(prefix string) ([]string, bool) {
	cur := t.root
	for _, runeVal := range prefix {
		if cur.child[runeVal] == nil {
			return nil, false
		}
		cur = cur.child[runeVal]
	}
	var resutls []string
	t.collect(cur, prefix, &resutls)
	return resutls, len(resutls) > 0
}

// Internal function to collect all the matches for a prefix in trie
func (t *Trie) collect(node *Node, currentWord string, results *[]string) {
	if node == nil {
		return
	}

	if node.isEnd {
		*results = append(*results, currentWord)
	}
	for key, childNode := range node.child {
		nextWord := currentWord + string(key)
		t.collect(childNode, nextWord, results)
	}
}
