package lexer

import "fmt"

type ItemType string

type Item struct {
	Type  ItemType
	Value string
	End   int
}

// TODO: Write a description of what each item does.
const (
	ItemError ItemType = "ERR"

	ItemEOF       ItemType = ""
	ItemComma     ItemType = "COMMA"      // , used to specify a two-address addresser
	ItemDollar    ItemType = "DOLLAR"     // $ used to specify last line
	ItemBackslash ItemType = "BACK-SLASH" // \ used for escaping
	ItemSlash     ItemType = "SLASH"      // / for dividing address
	ItemInt       ItemType = "INT"        // [0-9]+ used to specify line number
	ItemLit       ItemType = "LIT"
	ItemCmd       ItemType = "CMD"
	ItemDiv       ItemType = "DIV"
	ItemIdent     ItemType = "IDENT"
	ItemLBrace    ItemType = "L-BRACE"
	ItemRBrace    ItemType = "R-BRACE"
	ItemExpMark   ItemType = "EXP-MARK"
	ItemSemicolon ItemType = "SEMICOLON"
	ItemNewline   ItemType = "NEW-LINE"

	// TODO: Are some of these even used?
	ItemLParen   ItemType = "L-PAREN"
	ItemRParen   ItemType = "R-PAREN"
	ItemQuestion ItemType = "QUESTION"
	ItemPeriod   ItemType = "PERIOD"
	ItemColon    ItemType = "COLON"
)

func (i Item) String() string {
	switch i.Type {
	case ItemEOF:
		return "EOF"
	case ItemError:
		return i.Value
	}
	if len(i.Value) > 10 {
		return fmt.Sprintf("%10q...", i.Value)
	}
	return fmt.Sprintf("%q", i.Value)
}
