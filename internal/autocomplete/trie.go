package auto

import (
	"os"
	"strings"
)

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

func (t *Trie) IndexSystemPath() {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return
	}

	dirs := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.Type().IsRegular() {
				t.Add(entry.Name())
			}
		}
	}
}

func (t *Trie) LCP(prefix string) string {
	cur := t.root
	for _, c := range prefix {
		next, ok := cur.child[c]
		if !ok {
			return ""
		}
		cur = next
	}
	lsp := prefix
	for len(cur.child) == 1 && !cur.isEnd {
		//var nextRune rune
		for r, next := range cur.child {
			lsp += string(r)
			cur = next
		}
	}
	return lsp
}

// func (t *Trie) PrintAllFilesAndDirs(term *term.Terminal) {
// 	var names []string
// 	t.printAllWords(t.root, "", &names)
// 	for _, word := range names {
// 		fmt.Fprintln(term, word)
// 	}
// }

// func (t *Trie) printAllWords(node *Node, prefix string, names *[]string) {
// 	if node.isEnd {
// 		*names = append(*names, prefix)
// 	}
// 	keys := make([]rune, 0, len(node.child))
// 	for char := range node.child {
// 		keys = append(keys, char)
// 	}
// 	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
// 	for _, char := range keys {
// 		t.printAllWords(node.child[char], prefix+string(char), names)
// 	}
// }

// func (t *Trie) GetCurrentDirFiles() {
// 	// files, err := os.ReadDir(".")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// for _, entry := range files {
// 	// 	//fmt.Println(entry.Name())
// 	// 	t.Add(entry.Name())
// 	// }
// 	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !d.IsDir() {
// 			t.Add(path)f
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
//try with filepath.Walkdir in the file completion later too.
