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
	// noboundaryWords default to a-zA-Z0-9_
	noboundaryWords map[rune]bool
	// lock for the map write
	sync.RWMutex
}

type ExtractResult struct {
	Keyword string
	// StartIndex is the keyword index in the sentences
	StartIndex int
}

type Option struct {
	// Longest set to true will just match the longest keyword,
	Longest  bool
	SpanInfo bool
}

var (
	defaultOption = &Option{
		Longest:  true,
		SpanInfo: false,
	}
)

func NewKeywordProcessor() *keywordProcessor {
	p := &keywordProcessor{
		dicts:           make(map[string]string),
		noboundaryWords: make(map[rune]bool),
		keytrie:         NewTrie('r'),
	}
	for i := 0; i < 26; i++ {
		p.AddNoBoundaryWords(rune('a' + i))
		p.AddNoBoundaryWords(rune('A' + i))
	}
	for i := 0; i < 10; i++ {
		p.AddNoBoundaryWords(rune('0' + i))
	}
	p.AddNoBoundaryWords('-')
	return p
}

func (p *keywordProcessor) SetCaseSenstive(caseSenstive bool) {
	p.caseSensitive = caseSenstive
}

func (p *keywordProcessor) AddNoBoundaryWords(noboundaryWords ...rune) {
	for _, w := range noboundaryWords {
		p.noboundaryWords[w] = true
	}
}

func (p *keywordProcessor) AddKeywords(keywords ...string) {
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

func (p *keywordProcessor) ExtractKeywords(sentence string, option ...*Option) (res []*ExtractResult) {
	extractOption := defaultOption
	if len(option) > 0 {
		extractOption = option[0]
	}
	res = make([]*ExtractResult, 0, 20)
	if !p.caseSensitive {
		sentence = strings.ToLower(sentence)
	}
	runes := []rune(sentence)
	size := len(runes)
	idx := 0
	begin := true
	var curTrie *trie
	for idx < size {
		curTrie = p.keytrie
		c := runes[idx]
		if _, ok := p.noboundaryWords[c]; !ok {
			idx++
			begin = true
		} else if !begin {
			idx++
		} else {
			var j = idx
			foundKeyword := ""
			for j = idx; j < size; j++ {
				c = runes[j]
				curTrie = curTrie.next[c]
				if curTrie == nil {
					break
				}
				if curTrie.word != "" && (j == size-1 || !p.noboundaryWords[runes[j+1]]) {
					foundKeyword = curTrie.word
					if !extractOption.Longest {
						res = append(res, &ExtractResult{p.dicts[foundKeyword], idx})
						idx = j
					}
				}
			}
			if foundKeyword == "" {
				idx++
			} else if extractOption.Longest {
				res = append(res, &ExtractResult{p.dicts[foundKeyword], idx})
				idx = j
			}
			begin = false
		}
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
