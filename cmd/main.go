package main

import (
	"flag"
	"fmt"
	"iter"
	"lgg/internal/lexer"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println(".lgo files must be provided as arguments")
		return
	}

	fileName := flag.Arg(0)

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("invalid file: %v", err)
		return
	}
	// fmt.Println("content of ", fileName, " is ---------------")
	// fmt.Println(string(content))
	// fmt.Println("---------------------------------------------")

	reader := strings.NewReader(string(content))

	lxr := lexer.NewLexer(reader)
	lxr.Read()

	gen := NewGenerator(lxr.Seq())

	fmt.Println(gen.Generate())
}

type Generator struct {
	sequence iter.Seq[lexer.Lexem]
	sb       *strings.Builder
}

func NewGenerator(sequence iter.Seq[lexer.Lexem]) *Generator {
	return &Generator{
		sequence: sequence,
		sb:       &strings.Builder{},
	}
}

func (gen *Generator) Generate() string {
	next, _ := iter.Pull(gen.sequence)

	for {
		l, haveNext := next()
		if !haveNext {
			break
		}
		if l.Kind == lexer.ListOpen {
		} else if l.Kind == lexer.ListClose {
		} else if l.Kind == lexer.Symbol {
			if l.Content == "package" {
				if Break := gen.genPackage(next); Break {
					break
				}
			} else if l.Content == "func" {
				if Break := gen.genFunction(next); Break {
					break
				}
			}
		}
	}

	return gen.sb.String()
}

func (gen *Generator) genFunction(
	next func() (lexer.Lexem, bool),
) (Break bool) {

	fName, err := gen.parseLexem(next, lexer.Symbol, "func name")
	if err != nil {
		return true
	}
	_, err = gen.parseLexem(next, lexer.ListOpen, "param list")
	if err != nil {
		return true
	}
	_, err = gen.parseLexem(
		next,
		lexer.ListClose,
		"now only empty param list supported, end of list",
	)
	if err != nil {
		return true
	}

	fmt.Fprintf(gen.sb, "func %s () {\n", *fName)

	for {
		l, haveNext := next()
		if l.Kind == lexer.ListClose || !haveNext {
			break
		} else if l.Kind == lexer.ListOpen {
			// start of 'statement' list
		} else {
			fmt.Fprintf(os.Stderr, "end of function or start of 'statement list expected'")
			return true
		}

		function, err := gen.parseLexem(
			next,
			lexer.Symbol,
			"function to call name",
		)
		if err != nil {
			return true
		}

		var argument *string
		l, haveNext = next()
		if l.Kind == lexer.ListClose {
		} else if l.Kind == lexer.String {
			content := l.Content
			argument = &content
		} else {
			fmt.Fprintf(os.Stderr, "end of function or start of 'statement list expected'")
			return true
		}

		if argument != nil {
			_, err = gen.parseLexem(
				next,
				lexer.ListClose,
				"end of 'statement' list",
			)
			if err != nil {
				return true
			}
		}

		if argument != nil {
			fmt.Fprintf(gen.sb, "  %s(\"%s\")\n", *function, *argument)
		} else {
			fmt.Fprintf(gen.sb, "  %s()\n", *function)
		}
	}

	fmt.Fprintf(gen.sb, "}\n")

	return false
}

func (gen *Generator) parseLexem(
	next func() (lexer.Lexem, bool),
	kind lexer.Kind,
	lName string,
) (*string, error) {
	l, haveNext := next()
	return gen.checkLexem(l, haveNext, kind, lName)
}

func (gen *Generator) checkLexem(
	l lexer.Lexem,
	haveNext bool,
	kind lexer.Kind,
	lName string,
) (*string, error) {

	if !CheckKind(haveNext, l, kind) {
		fmt.Fprintf(os.Stderr, "%s extected", lName)
		return nil, fmt.Errorf("%s expected", lName)
	}

	res := l.Content
	return &res, nil
}

func (gen *Generator) genPackage(next func() (lexer.Lexem, bool)) (Break bool) {
	l, haveNext := next()
	if !CheckKind(haveNext, l, lexer.Symbol) {
		fmt.Fprintf(os.Stderr, "package name extected")
		return true
	}

	fmt.Fprintf(gen.sb, "package %s\n\n", l.Content)

	l, haveNext = next()
	if !CheckKind(haveNext, l, lexer.ListClose) {
		return true
	}

	return false
}

func CheckKind(haveNext bool, have lexer.Lexem, expected lexer.Kind) bool {
	if !haveNext {
		fmt.Fprintf(
			os.Stderr,
			"unexpected end of file, %++v expected\n",
			expected,
		)
		return false
	}
	if have.Kind != expected {
		fmt.Fprintf(
			os.Stderr,
			"got unexpected %++v when %++v expected\n",
			have,
			expected,
		)
		return false
	}
	return true
}
