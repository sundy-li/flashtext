package flashtext

import (
	"strings"
	"sync"
)

// keywordProcessor is the processor of keyword extract
type keywordProcessor struct {
	// dicts store the keyword => cleanName
	dicts map[string]string
	// keytrie is the trie struct
	keytrie *trie
	// caseSensitive or not
	caseSensitive bool
	// boundaryWords
	boundaryWords map[rune]bool
	// lock for the map write
	sync.RWMutex
}

func NewKeywordProcessor() *keywordProcessor {
	p := &keywordProcessor{
		dicts:         make(map[string]string),
		boundaryWords: make(map[rune]bool),
		keytrie:       NewTrie('r'),
	}
	p.AddBoundaryWords('.', '\t', '\n', '\a', ' ', ',', '/')
	return p
}

func (p *keywordProcessor) SetConfig(caseSenstive bool) {
	p.caseSensitive = caseSenstive
}

func (p *keywordProcessor) AddBoundaryWords(boundaryWords ...rune) {
	for _, w := range boundaryWords {
		p.boundaryWords[w] = true
	}
}

func (p *keywordProcessor) AddKeyword(keywords ...string) {
	for _, keyword := range keywords {
		p.AddKeywordAndName(keyword, keyword)
	}
}

func (p *keywordProcessor) AddKeywordAndName(keyword string, cleanName string) {
	p.Lock()
	defer p.Unlock()

	if !p.caseSensitive {
		keyword = strings.ToLower(keyword)
	}
	p.keytrie.addKeyword(keyword)
	p.dicts[keyword] = cleanName
}

func (p *keywordProcessor) Extracts(sentence string, longest bool) (res []string) {
	var pos = p.keytrie

	res = make([]string, 0, 20)
	if !p.caseSensitive {
		sentence = strings.ToLower(sentence)
	}
	runes := []rune(sentence)
	size := len(runes)
	for i := 0; i < size; {
		if i != 0 && len(p.boundaryWords) != 0 {
			if _, ok := p.boundaryWords[runes[i-1]]; !ok {
				i = i + 1
				continue
			}
		}
		pos = p.keytrie
		var matchedStr string
		var jump = i + 1
		for j := i; j < size; j++ {
			r := runes[j]
			if next, ok := pos.next[r]; ok {
				pos = next
				if pos.endpoint {
					matchedStr = p.dicts[string(runes[i:j+1])]
					if j+1 != size {
						if len(p.boundaryWords) != 0 {
							if _, ok := p.boundaryWords[runes[j+1]]; !ok {
								matchedStr = ""
								continue
							}
						}
					}
					if !longest {
						if jump == i+1 {
							jump = j + 1
						}
						res = append(res, matchedStr)
					} else {
						jump = j + 1
					}
				}
			} else {
				break
			}
		}
		if longest && matchedStr != "" {
			res = append(res, matchedStr)
		}
		i = jump
	}

	return res
}

func (p *keywordProcessor) RemoveKeywords(keywords ...string) {
	p.Lock()
	defer p.Unlock()
	for _, keyword := range keywords {
		if !p.caseSensitive {
			keyword = strings.ToLower(keyword)
		}
		p.keytrie.removeKeyword(keyword)
	}
}

func (p *keywordProcessor) Exists(keyword string) bool {
	return p.keytrie.exists(keyword)
}
