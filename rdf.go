package gon3

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Term interface {
	String() string
	Equals(Term) bool
}

// This must be a full (i.e. not relative IRI)
type IRI struct {
	url *url.URL
}

func (i IRI) String() string {
	return fmt.Sprintf("<%s>", i.url)
}

func (i IRI) Equals(other Term) bool {
	switch other.(type) {
	case IRI:
		break
	default:
		return false
	}
	return i.String() == other.String()
}

func newIRIFromString(s string) (IRI, error) {
	url, err := iriRefToURL(s)
	return IRI{url}, err
}

func iriRefToURL(s string) (*url.URL, error) {
	// strip <>, unescape, parse into url
	if strings.HasPrefix(s, "<") {
		s = s[1 : len(s)-1]
	}
	unescaped := unescapeUChar(s)
	return url.Parse(unescaped)
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-blank-node
type BlankNode struct {
	Id    int
	Label string
}

func (b BlankNode) String() string {
	return fmt.Sprintf("_:%s", b.Label)
}

func (b BlankNode) Equals(other Term) bool {
	panic("unimplemented")
}

func isBlankNode(t Term) bool {
	switch t.(type) {
	case BlankNode:
		return true
	}
	return false
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-literal
type Literal struct {
	LexicalForm string
	DatatypeIRI IRI
	LanguageTag string
}

func (l Literal) String() string {
	if l.LanguageTag != "" {
		return fmt.Sprintf("%q@%s", l.LexicalForm, l.LanguageTag)
	}
	return fmt.Sprintf("%q^^%s", l.LexicalForm, l.DatatypeIRI)
}

func (l Literal) Equals(other Term) bool {
	switch other.(type) {
	case Literal:
		break
	default:
		return false
	}
	otherLit := other.(Literal)
	return l.LexicalForm == otherLit.LexicalForm && l.DatatypeIRI.Equals(otherLit.DatatypeIRI) && l.LanguageTag == otherLit.LanguageTag
}

func lexicalForm(s string) string {
	var unquoted string
	if strings.HasPrefix(s, `"""`) || strings.HasPrefix(s, `'''`) {
		unquoted = s[3 : len(s)-3]
	} else {
		unquoted = s[1 : len(s)-1]
	}
	// TODO: resolve escapes
	u := unescapeUChar(unquoted)
	ret := unescapeEChar(u)
	return ret
}

func unescapeEChar(s string) string {
	var replacements = []struct {
		old string
		new string
	}{
		{`\t`, "\t"},
		{`\b`, "\b"},
		{`\n`, "\n"},
		{`\r`, "\r"},
		{`\f`, "\f"},
		{`\"`, `"`},
		{`\'`, `'`},
		{`\\`, `\`},
	}
	for _, r := range replacements {
		s = strings.Replace(s, r.old, r.new, -1)
	}
	return s
}

func unescapeUChar(s string) string {
	for {
		var start, hex, end string
		uIdx := strings.Index(s, `\u`)
		UIdx := strings.Index(s, `\U`)
		if uIdx >= 0 {
			start = s[:uIdx]
			if uIdx+6 > len(s) {
				hex = s[uIdx+2:]
				end = ""
			} else {
				hex = s[uIdx+2 : uIdx+6]
				end = s[uIdx+6:]
			}
		} else if UIdx >= 0 {
			start = s[:UIdx]
			if UIdx+10 > len(s) {
				hex = s[UIdx+2:]
				end = ""
			} else {
				hex = s[UIdx+2 : UIdx+10]
				end = s[UIdx+10:]
			}
		} else {
			break
		}
		num, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			panic(err) // TODO: this shouldn't happen
		}
		s = fmt.Sprintf("%s%s%s", start, string(rune(num)), end)
	}
	return s
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-rdf-triple
type Triple struct {
	Subject, Predicate, Object Term
}

func (t *Triple) String() string {
	return fmt.Sprintf("%s %s %s .", t.Subject, t.Predicate, t.Object)
}

func (t *Triple) includes(term Term) bool {
	for node := range t.IterNodes() {
		if node.Equals(term) {
			return true
		}
	}
	return false
}

func (t *Triple) IterNodes() <-chan Term {
	panic("unimplemented")
}

// An RDF graph is a set of RDF triples
type Graph struct {
	triples []*Triple
	uri     IRI
}

func (g Graph) String() string {
	str := ""
	i := -1
	for t := range g.IterTriples() {
		i += 1
		if i > 0 {
			str += "\n"
		}
		str = fmt.Sprintf("%s%s", str, t)
	}
	return str
}

func (g Graph) IterNodes() <-chan Term {
	panic("unimplemented")
}

func (g Graph) IterTriples() <-chan *Triple {
	panic("unimplemented")
}
