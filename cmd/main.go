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

	sb := &strings.Builder{}
	listIdxs := make([]int, 0)
	currIdx := -1

	next, _ := iter.Pull(lxr.Seq())

	for {
		l, haveNext := next()
		if !haveNext {
			break
		}
		if l.Kind == lexer.ListOpen {
			listIdxs = append(listIdxs, currIdx)
			currIdx = 0
		} else if l.Kind == lexer.ListClose {
			lastIdx := len(listIdxs) - 1
			currIdx = listIdxs[lastIdx]
			listIdxs = listIdxs[:lastIdx]
		} else if currIdx == 0 && l.Kind == lexer.Symbol {
			// fmt.Fprintf(os.Stderr, "%++v\n", l)
			if l.Content == "package" {

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.Symbol) {
					fmt.Fprintf(os.Stderr, "package name extected")
					break
				}

				fmt.Fprintf(sb, "package %s\n\n", l.Content)

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.ListClose) {
					lastIdx := len(listIdxs) - 1
					currIdx = listIdxs[lastIdx]
					listIdxs = listIdxs[:lastIdx]
					break
				}

			} else if l.Content == "func" {

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.Symbol) {
					fmt.Fprintf(os.Stderr, "func name extected")
					break
				}
				fName := l.Content

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.ListOpen) {
					fmt.Fprintf(os.Stderr, "param list expected")
					break
				}

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.ListClose) {
					fmt.Fprintf(os.Stderr, "now only empty param list supported")
					break
				}

				fmt.Fprintf(sb, "func %s () {\n", fName)

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.ListOpen) {
					fmt.Fprintf(os.Stderr, "'statement' list expected")
					break
				}

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.Symbol) {
					fmt.Fprintf(os.Stderr, "function to call name extected")
					break
				}
				function := l.Content

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.String) {
					fmt.Fprintf(os.Stderr, "function argument extected")
					break
				}
				argument := l.Content

				l, haveNext = next()
				if !CheckKind(haveNext, l, lexer.ListClose) {
					fmt.Fprintf(os.Stderr, "end of 'statement' list expected")
					break
				}
				fmt.Fprintf(sb, "  %s(\"%s\")\n", function, argument)

				fmt.Fprintf(sb, "}\n")
			}
			currIdx++
		}
	}

	fmt.Println(sb.String())
}

func CheckKind(haveNext bool, have lexer.Lexem, expected lexer.Kind) bool {
	if !haveNext {
		fmt.Printf("unexpected end of file, %++v expected\n", expected)
		return false
	}
	if have.Kind != expected {
		fmt.Printf("got unexpected %++v when %++v expected\n", have, expected)
		return false
	}
	return true
}
