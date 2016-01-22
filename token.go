package gon3

import (
	"fmt"
	"github.com/rychipman/easylex"
)

const (
	// tokens expressed as literal strings in http://www.w3.org/TR/turtle/#sec-grammar-grammar
	tokenAtPrefix easylex.Token = iota
	tokenAtBase
	tokenEndTriples
	tokenA
	tokenPredicateListSeparator
	tokenObjectListSeparator
	tokenStartBlankNodePropertyList
	tokenEndBlankNodePropertyList
	tokenStartCollection
	tokenEndCollection
	tokenLiteralDatatypeTag // TODO: rename
	tokenTrue
	tokenFalse

	// terminal tokens from http://www.w3.org/TR/turtle/#terminals
	tokenIRIRef
	tokenPNameNS
	tokenPNameLN
	tokenBlankNodeLabel
	tokenLangTag
	tokenInteger
	tokenDecimal
	tokenDouble
	tokenExponent
	tokenStringLiteralQuote
	tokenStringLiteralSingleQuote
	tokenStringLiteralLongQuote
	tokenStringLiteralLongSingleQuote
	tokenAnon
)
