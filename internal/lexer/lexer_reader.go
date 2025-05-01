package lexer

import "iter"

type LexemReaderList struct {
	lexems []Lexem
}

func (r *LexemReaderList) AddLexem(l Lexem) bool {
	r.lexems = append(r.lexems, l)
	return true
}

func (r *LexemReaderList) Seq() iter.Seq[Lexem] {
	return func(yield func(Lexem) bool) {
		for _, l := range r.lexems {
			if !yield(l) {
				break
			}
		}
	}
}
