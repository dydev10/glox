package parser

import (
	"fmt"

	"github.com/dydev10/glox/ast"
	"github.com/dydev10/glox/lexer"
)

type Parser struct {
	tokens  []*lexer.Token
	current int
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ast.Expr[string] {
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

func (p *Parser) consume(t lexer.TokenType, m string) lexer.Token {
	if p.check(t) {
		return p.advance()
	}

	//
	panic(fmt.Sprint(p.peek(), m))
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

func (p *Parser) expression() ast.Expr[string] {
	return p.equality()
}

func (p *Parser) equality() ast.Expr[string] {
	expr := p.comparison()

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr[string] {
	expr := p.term()

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) term() ast.Expr[string] {
	expr := p.factor()

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) factor() ast.Expr[string] {
	expr := p.unary()

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.Binary[string]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) unary() ast.Expr[string] {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right := p.unary()

		return &ast.Unary[string]{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr[string] {
	if p.match(lexer.FALSE) {
		return &ast.Literal[string]{Value: "false"}
	}
	if p.match(lexer.TRUE) {
		return &ast.Literal[string]{Value: "true"}
	}
	if p.match(lexer.NIL) {
		return &ast.Literal[string]{}
	}

	if p.match(lexer.NUMBER, lexer.STRING) {
		return &ast.Literal[string]{Value: p.previous().Literal}
	}

	if p.match(lexer.LEFT_PAREN) {
		expr := p.expression()
		p.consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.Grouping[string]{
			Expression: expr,
		}
	}

	panic("Unknown primary rule!!")
}
