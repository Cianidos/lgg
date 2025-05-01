package lexer

const (
	ListOpen  Kind = iota // "("
	ListClose             // ")"
	Symbol
	String
	Number

	Unknown
)

type Location struct {
	Offset uint32
}

type Kind int

type Lexem struct {
	Kind
	Content string
	Location
}
