package parse

type parser struct {
	l        *lexer
	last     item
	rejected bool
	qCount   int
}

func (p *parser) scan() item {
	// there is one rejected item
	if p.rejected {
		p.rejected = false
		return p.last
	}

	p.last = p.l.nextItem()
	p.rejected = false
	if p.last.typ == ItemQuestionMark {
		p.qCount++
	}
	return p.last
}

func (p *parser) scanIgnoreWhiteSpace() item {
	t := p.scan()
	if t.typ == ItemWhiteSpace {
		t = p.scan()
	}
	return t
}

func (p *parser) reject() {
	p.rejected = true
}

// AST return the abstract source tree for given query, currently only select is supported
func AST(q string) (*Query, error) {
	p := &parser{
		l: lex(q),
	}
	s, err := newStatement(p)
	if err != nil {
		p.l.drain() // make sure the lexer is terminated
		return nil, err
	}

	return &Query{
		Statement: s,
	}, nil
}
