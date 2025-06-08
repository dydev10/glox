package parser

import (
	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/lexer"
)

type Parser struct {
	tokens  []*lexer.Token
	current int
	Errors  []*ParseError
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (ast.Expr[string], error) {
	return p.expression()
}

func (p *Parser) peek() *lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *lexer.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == lexer.EOF
}

func (p *Parser) advance() *lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) match(tokenTypes ...lexer.TokenType) bool {
	// if slices.ContainsFunc(tokenTypes, p.check) {
	// 	p.advance()
	// 	return true
	// }
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t lexer.TokenType, m string) (*lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	err := p.logError(m)
	return nil, err
}

/*
Language grammar rule functions:
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
*/

func (p *Parser) expression() (ast.Expr[string], error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr[string], error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr[string], error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr[string], error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr[string], error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr[string], error) {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.Unary[string]{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr[string], error) {
	if p.match(lexer.FALSE) {
		return &ast.Literal[string]{Value: "false"}, nil
	}
	if p.match(lexer.TRUE) {
		return &ast.Literal[string]{Value: "true"}, nil
	}
	if p.match(lexer.NIL) {
		return &ast.Literal[string]{}, nil
	}

	if p.match(lexer.NUMBER, lexer.STRING) {
		return &ast.Literal[string]{Value: p.previous().Literal}, nil
	}

	if p.match(lexer.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, consumeErr := p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		if consumeErr != nil {
			return nil, consumeErr
		}
		return &ast.Grouping[string]{
			Expression: expr,
		}, nil
	}

	err := p.logError("Expect expression.")
	return nil, err
}

// save errors
func (p *Parser) logError(message string) *ParseError {
	err := &ParseError{
		token:   p.peek(),
		message: message,
	}
	p.Errors = append(p.Errors, err)
	return err
}
