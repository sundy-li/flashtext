package flashtext

type trie struct {
	key  rune
	next map[rune]*trie
	word string
}

func NewTrie(r rune) *trie {
	return &trie{
		key:  r,
		next: make(map[rune]*trie),
	}
}

func (t *trie) addKeyword(keyword string) {
	var pos = t
	for _, r := range keyword {
		pos = pos.getOrSet(r)
	}
	pos.word = keyword
}

func (t *trie) removeKeyword(keyword string) {
	var pos = t
	for _, r := range keyword {
		pos = pos.getOrSet(r)
	}
	//fake delete
	pos.word = ""
}

func (t *trie) getOrSet(r rune) *trie {
	if next, ok := t.next[r]; ok {
		return next
	}
	next := NewTrie(r)
	t.next[r] = next
	return next
}

func (t *trie) exists(keyword string) bool {
	var keytrie = t
	for _, r := range keyword {
		if next, ok := keytrie.next[r]; ok {
			keytrie = next
		} else {
			return false
		}
	}
	return keytrie.word != ""
}
