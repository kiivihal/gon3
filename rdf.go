package gon3

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Term interface {
	String() string
	Equals(Term) bool
	RawValue() string
}

// This must be a full (i.e. not relative IRI)
type IRI struct {
	url string
}

func NewIRI(iri string) *IRI {
	return &IRI{
		url: iri,
	}
}

func (i *IRI) String() string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("<%s>", i.url)
}

func (i *IRI) Equals(other Term) bool {
	switch other.(type) {
	case *IRI:
		break
	default:
		return false
	}
	return i.String() == other.String()
}

func (i *IRI) RawValue() string {
	return fmt.Sprintf("%s", i.url)
}

func newIRIFromString(s string) *IRI {
	url := iriRefToURL(s)
	return &IRI{url}
}

func iriRefToURL(s string) string {
	// strip <>, unescape, parse into url
	if strings.HasPrefix(s, "<") {
		s = s[1 : len(s)-1]
	}
	return unescapeUChar(s)
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-blank-node
type BlankNode struct {
	ID    int
	Label string
}

func NewBlankNode(label string) *BlankNode {
	return &BlankNode{
		Label: label,
	}
}

func (b *BlankNode) String() string {
	return fmt.Sprintf("_:%s", b.Label)
}

func (b *BlankNode) Equals(other Term) bool {
	switch other.(type) {
	case *BlankNode:
		return true
	default:
		return false
	}
	panic("unreachable")
}

func (b *BlankNode) RawValue() string {
	return b.Label
}

func isBlankNode(t Term) bool {
	switch t.(type) {
	case *BlankNode:
		return true
	}
	return false
}

// see http://www.w3.org/TR/rdf11-concepts/#dfn-literal
type Literal struct {
	LexicalForm string
	DatatypeIRI *IRI
	LanguageTag string
}

func NewLiteral(label string) *Literal {
	return &Literal{
		LexicalForm: label,
	}
}

func NewLiteralWithDataType(label string, dtype *IRI) *Literal {
	return &Literal{
		LexicalForm: label,
		DatatypeIRI: dtype,
	}
}

func NewLiteralWithLanguage(label string, lang string) *Literal {
	return &Literal{
		LexicalForm: label,
		LanguageTag: lang,
	}
}

func (l *Literal) String() string {
	if l.LanguageTag != "" {
		return fmt.Sprintf("%q@%s", l.LexicalForm, l.LanguageTag)
	}
	return fmt.Sprintf("%q^^%s", l.LexicalForm, l.DatatypeIRI)
}

func (l *Literal) Equals(other Term) bool {
	switch other.(type) {
	case *Literal:
		break
	default:
		return false
	}
	otherLit := other.(*Literal)
	return l.LexicalForm == otherLit.LexicalForm && l.DatatypeIRI.Equals(otherLit.DatatypeIRI) && l.LanguageTag == otherLit.LanguageTag
}

func (l *Literal) RawValue() string {
	return l.LexicalForm
}

func lexicalForm(s string) string {
	var unquoted string
	if strings.HasPrefix(s, `"""`) || strings.HasPrefix(s, `'''`) {
		unquoted = s[3 : len(s)-3]
	} else {
		unquoted = s[1 : len(s)-1]
	}
	u := unescapeUChar(unquoted)
	ret := unescapeEChar(u)
	return ret
}

func unescapeEChar(s string) string {
	replacements := []struct {
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
		s = strings.ReplaceAll(s, r.old, r.new)
	}
	return s
}

func getIndexOfEscape(s string, substr string) int {
	for {
		idx := strings.Index(s, substr)
		if idx < 0 {
			return idx
		}

		// search through runes backward from the index to ensure the escape isn't escaped
		var size int
		count := 1
		escapeRune := []rune(substr)[0]
		var r rune
		for i := idx; i > 0; i -= size {
			r, size = utf8.DecodeLastRuneInString(s[:i])
			if r != escapeRune {
				break
			}
			count++
		}

		// an odd number of escape characters indicates it's not escaped
		if count%2 == 1 {
			return idx
		}

		// skip that false match and check the rest of the string
		idx += len(substr)
		s = s[idx:]
	}
}

func unescapeUChar(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); {
		// Not an escape sequence
		if i >= len(s)-1 || s[i] != '\\' {
			result.WriteByte(s[i])
			i++
			continue
		}

		// Handle escaped backslash
		if s[i+1] == '\\' {
			result.WriteByte('\\')
			i += 2
			continue
		}

		// Check for Unicode escape
		if i+1 < len(s) && (s[i+1] == 'u' || s[i+1] == 'U') {
			seqLen := 4
			if s[i+1] == 'U' {
				seqLen = 8
			}

			// Check if we have enough characters for a complete sequence
			if i+2+seqLen > len(s) {
				result.WriteString(s[i:])
				break
			}

			// Extract the hex digits
			hexDigits := s[i+2 : i+2+seqLen]

			// Verify all characters are valid hex digits
			validHex := true
			for _, c := range hexDigits {
				if !unicode.Is(unicode.ASCII_Hex_Digit, rune(c)) {
					validHex = false
					break
				}
			}

			if validHex {
				// Try to unquote the sequence
				if unquoted, err := strconv.Unquote(`"` + s[i:i+2+seqLen] + `"`); err == nil {
					result.WriteString(unquoted)
					i += 2 + seqLen
					continue
				}
			}

			// For invalid or incomplete sequences, write just the first part
			// and continue processing from the next character
			result.WriteString(s[i : i+2]) // Write \u or \U
			i += 2
			continue
		}

		// Not a recognized escape sequence
		result.WriteByte(s[i])
		i++
	}
	return result.String()
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
	ch := make(chan Term, 3)
	ch <- t.Subject
	ch <- t.Predicate
	ch <- t.Object
	close(ch)
	return ch
}

// An RDF graph is a set of RDF triples
type Graph struct {
	triples []*Triple
	uri     *IRI
}

func (g *Graph) IsomorphicTo(other *Graph) bool {
	cg1 := g.Canonicalize()
	cg2 := other.Canonicalize()
	return cg1.IsomorphicTo(cg2)
}

func (g *Graph) String() string {
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

func (g *Graph) IterTriples() <-chan *Triple {
	ch := make(chan *Triple)
	go func() {
		for _, trip := range g.triples {
			ch <- trip
		}
		close(ch)
	}()
	return ch
}

func (g *Graph) NodesSorted() []Term {
	set := make(map[string]Term)
	for t := range g.IterTriples() {
		for n := range t.IterNodes() {
			if _, has := set[n.String()]; !has {
				set[n.String()] = n
			}
		}
	}
	terms := make([]Term, 0)
	for _, t := range set {
		terms = append(terms, t)
	}
	termsSlice := TermSlice(terms)
	sort.Sort(termsSlice)
	fmt.Printf("nodes before sort: %+v\nafter:             %+v\n", terms, termsSlice)
	return termsSlice
}

type TermSlice []Term

func (ts TermSlice) Len() int {
	return len(ts)
}

func (ts TermSlice) Less(i, j int) bool {
	iNode := ts[i]
	jNode := ts[j]
	iPriority := 0
	jPriority := 0
	switch iNode.(type) {
	case *BlankNode:
		iPriority = 1
	case *Literal:
		iPriority = 2
	case *IRI:
		iPriority = 3
	}
	switch jNode.(type) {
	case *BlankNode:
		jPriority = 1
	case *Literal:
		jPriority = 2
	case *IRI:
		jPriority = 3
	}
	if iPriority > jPriority {
		return true
	} else if jPriority > iPriority {
		return false
	}
	return strings.Compare(iNode.String(), jNode.String()) > 0
}

func (ts TermSlice) Swap(i, j int) {
	iNode := ts[i]
	jNode := ts[j]
	ts[i] = jNode
	ts[j] = iNode
}
