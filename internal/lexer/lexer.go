package lexer

import (
	"io"
	"iter"
	"slices"
	"strings"
	"unicode"
)

type LexemReader interface {
	AddLexem(Lexem) bool
	Seq() iter.Seq[Lexem]
}

type Lexer struct {
	LexemReader
	reader *strings.Reader
}

func NewLexer(reader *strings.Reader) *Lexer {
	return &Lexer{reader: reader, LexemReader: &LexemReaderList{}}
}

func (l *Lexer) Read() {
	for {
		r, _, err := l.reader.ReadRune()
		if err == io.EOF {
			break
		}
		switch {
		case slices.Contains([]rune{' ', '\t', '\n'}, r):
			continue
		case r == '(':
			l.AddLexem(
				Lexem{Kind: ListOpen, Content: "(", Location: Location{}},
			)
		case r == ')':
			l.AddLexem(
				Lexem{Kind: ListClose, Content: ")", Location: Location{}},
			)
		case slices.Contains([]rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}, r):
			l.reader.UnreadRune()
			l.ReadNumber()
		case r == '"':
			l.ReadString()
		case unicode.IsLetter(r):
			l.reader.UnreadRune()
			l.ReadSymbol()
		default:
			l.AddLexem(Lexem{
				Kind:     Unknown,
				Content:  string(r),
				Location: Location{},
			})
		}
	}
}

func (l *Lexer) ReadSymbol() {
	sym := ""
	for {
		r, _, err := l.reader.ReadRune()
		// fmt.Printf("wtf %c\n", r)
		switch {
		case err == io.EOF:
			return
		case unicode.IsLetter(r):
			sym += string(r)
		default:
			l.AddLexem(Lexem{
				Kind:     Symbol,
				Content:  sym,
				Location: Location{},
			})
			l.reader.UnreadRune()
			return
		}
	}
}

func (l *Lexer) ReadString() {
	str := ""
	for {
		r, _, err := l.reader.ReadRune()
		if err == io.EOF {
			return
		}
		switch r {
		case '"':
			l.AddLexem(Lexem{
				Kind:     String,
				Content:  str,
				Location: Location{},
			})
			return
		default:
			str += string(r)
		}
	}
}

func (l *Lexer) ReadNumber() {
	num := ""
	for {
		r, _, err := l.reader.ReadRune()
		if err == io.EOF {
			return
		}
		switch r {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			num += string(r)
		default:
			l.AddLexem(Lexem{
				Kind:     Number,
				Content:  num,
				Location: Location{},
			})
			l.reader.UnreadRune()
			return
		}
	}
}
