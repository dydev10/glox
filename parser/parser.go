package parser

import (
	"fmt"

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

// main entry point to run glox statements
func (p *Parser) Parse() ([]ast.Stmt, error) {
	var statements []ast.Stmt

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			// TODO: probably not return here, return error after loop along with result, because synchronize should be used now
			// or call synchronize() here instead of inside declaration??
			// p.synchronize() // remove return err, just sync on error
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

// alternate entry point to parse single expression without statement
func (p *Parser) ParseExpression() (ast.Expr, error) {
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

/**
* expression parsing
*
* Language grammar expression rule functions:
* expression     → assignment ;
*	assignment     → ( call "." )? IDENTIFIER "=" assignment | logic_or ;
* logic_or       → logic_and ( "or" logic_and )* ;
* logic_and      → equality ( "and" equality )* ;
* equality       → comparison ( ( "!=" | "==" ) comparison )* ;
* comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
* term           → factor ( ( "-" | "+" ) factor )* ;
* factor         → unary ( ( "/" | "*" ) unary )* ;
* unary          → ( "!" | "-" ) unary | call ;
* call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
* arguments      → expression ( "," expression )* ;
* primary        → "true" | "false" | "nil" | NUMBER | STRING | "(" expression ")"| IDENTIFIER ;
 */

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(lexer.EQUAL) {
		equals := p.previous()
		value, valueErr := p.assignment()
		if valueErr != nil {
			return nil, valueErr
		}

		if varExpr, ok := expr.(*ast.Variable); ok {
			return &ast.Assign{
				Name:  varExpr.Name,
				Value: value,
			}, nil
		} else if getExpr, ok := expr.(*ast.Get); ok {
			return &ast.Set{
				Object: getExpr.Object,
				Name:   getExpr.Name,
				Value:  value,
			}, nil
		} else {
			// log invalid assignment target error, but don't return
			// TODO: log this error instead of just printing here
			assignTargetErr := &ParseError{token: equals, message: "Invalid assignment target."}
			fmt.Println(assignTargetErr)
		}

	}

	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.OR) {
		operator := p.previous()

		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Right:    right,
			Operator: operator,
		}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(lexer.AND) {
		operator := p.previous()

		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.Logical{
			Left:     expr,
			Right:    right,
			Operator: operator,
		}
	}

	return expr, nil
}

func (p *Parser) equality() (ast.Expr, error) {
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
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
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
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
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
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
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
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.call()
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := []ast.Expr{}

	if !p.check(lexer.RIGHT_PAREN) {
		// do-while loop
		for hasArg := true; hasArg; hasArg = p.match(lexer.COMMA) {
			// print too many arguments error but don't stop parsing
			if len(arguments) >= 255 {
				tooManyArgsErr := p.logError("Can't have more than 255 arguments.")
				fmt.Println(tooManyArgsErr)
			}

			arg, argErr := p.expression()
			if argErr != nil {
				return nil, argErr
			}
			arguments = append(arguments, arg)
		}
	}

	paren, parenErr := p.consume(lexer.RIGHT_PAREN, "Expect ')' after arguments.")
	if parenErr != nil {
		return nil, parenErr
	}

	return &ast.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	// unconditional loop, depends on break to exit loop
	for {
		if p.match(lexer.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				// break loop to return error
				return nil, err
			}
		} else if p.match(lexer.DOT) {
			name, nameErr := p.consume(lexer.IDENTIFIER, "Expect property name after '.'.")
			if nameErr != nil {
				return nil, nameErr
			}
			expr = &ast.Get{
				Object: expr,
				Name:   name,
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(lexer.FALSE) {
		return &ast.Literal{Value: false}, nil
	}
	if p.match(lexer.TRUE) {
		return &ast.Literal{Value: true}, nil
	}
	if p.match(lexer.NIL) {
		return &ast.Literal{}, nil
	}
	if p.match(lexer.NUMBER, lexer.STRING) {
		return &ast.Literal{Value: p.previous().Literal}, nil
	}
	if p.match(lexer.THIS) {
		return &ast.This{Keyword: p.previous()}, nil
	}
	if p.match(lexer.IDENTIFIER) {
		return &ast.Variable{Name: p.previous()}, nil
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
		return &ast.Grouping{
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

// synchronize on parsing errors
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == lexer.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case lexer.CLASS:
			fallthrough
		case lexer.FUN:
			fallthrough
		case lexer.VAR:
			fallthrough
		case lexer.FOR:
			fallthrough
		case lexer.IF:
			fallthrough
		case lexer.WHILE:
			fallthrough
		case lexer.PRINT:
			fallthrough
		case lexer.RETURN:
			return
		default:
			p.advance()
		}
	}
}

/**
* statements parsing
*
* Language grammar statements rule functions:
*	program        → declaration* EOF ;
*	declaration    → classDecl | funDecl | varDecl | statement ;
* classDecl      → "class" IDENTIFIER "{" function* "}" ;
*	funDecl        → "fun" function ;
*	function       → IDENTIFIER "(" parameters? ")" block ;
*	parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
*	varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
*	statement      → exprStmt | forStmt | ifStm | printStmt | returnStmt | whileStmt | block;
* forStmt        → "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
* whileStmt      → "while" "(" expression ")" statement ;
* ifStmt         → "if" "(" expression ")" statement ( "else" statement )? ;
* block          → "{" declaration* "}" ;
*	exprStmt       → expression ";" ;
*	printStmt      → "print" expression ";" ;
* returnStmt     → "return" expression? ";" ;
 */

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(lexer.CLASS) {
		return p.classDeclaration()
	}

	if p.match(lexer.FUN) {
		return p.function("function")
	}

	if p.match(lexer.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(lexer.FOR) {
		return p.forStatement()
	}

	if p.match(lexer.IF) {
		return p.ifStmt()
	}

	if p.match(lexer.PRINT) {
		return p.printStatement()
	}

	if p.match(lexer.RETURN) {
		return p.returnStatement()
	}

	if p.match(lexer.WHILE) {
		return p.whileStatement()
	}

	if p.match(lexer.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.Block{Statements: statements}, nil
	}

	return p.expressionStatement()
}

func (p *Parser) classDeclaration() (ast.Stmt, error) {
	name, nameErr := p.consume(lexer.IDENTIFIER, "Expect class name.")
	if nameErr != nil {
		return nil, nameErr
	}

	if _, err := p.consume(lexer.LEFT_BRACE, "Expect '{' before class body."); err != nil {
		return nil, err
	}

	methods := []*ast.Function{}
	for !p.check(lexer.RIGHT_BRACE) && !p.isAtEnd() {
		function, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, function)
	}

	if _, err := p.consume(lexer.RIGHT_BRACE, "Expect '}' after class body."); err != nil {
		return nil, err
	}

	return &ast.Class{
		Name:    name,
		Methods: methods,
	}, nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	_, lParenErr := p.consume(lexer.LEFT_PAREN, "Expect '(' after 'for'.")
	if lParenErr != nil {
		return nil, lParenErr
	}

	var initializer ast.Stmt
	var initializerErr error
	if p.match(lexer.SEMICOLON) {
		initializer = nil
	} else if p.match(lexer.VAR) {
		initializer, initializerErr = p.varDeclaration()
	} else {
		initializer, initializerErr = p.expressionStatement()
	}

	if initializerErr != nil {
		return nil, initializerErr
	}

	var condition ast.Expr
	var condErr error
	if !p.check(lexer.SEMICOLON) {
		condition, condErr = p.expression()
		if condErr != nil {
			return nil, condErr
		}
	}

	_, condCloseErr := p.consume(lexer.SEMICOLON, "Expect ';' after loop condition.")
	if condCloseErr != nil {
		return nil, condCloseErr
	}

	var increment ast.Expr
	var incErr error
	if !p.check(lexer.RIGHT_PAREN) {
		increment, incErr = p.expression()
		if incErr != nil {
			return nil, incErr
		}
	}

	_, rParenErr := p.consume(lexer.RIGHT_PAREN, "Expect ')' after for clauses.")
	if rParenErr != nil {
		return nil, rParenErr
	}

	body, bodyErr := p.statement()
	if bodyErr != nil {
		return nil, bodyErr
	}

	if increment != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				body,
				&ast.Expression{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (p *Parser) ifStmt() (ast.Stmt, error) {
	_, lParenErr := p.consume(lexer.LEFT_PAREN, "Expect '(' after 'if'.")
	if lParenErr != nil {
		return nil, lParenErr
	}

	condition, condErr := p.expression()
	if condErr != nil {
		return nil, condErr
	}

	_, rParenErr := p.consume(lexer.RIGHT_PAREN, "Expect ')' after if condition.")
	if rParenErr != nil {
		return nil, rParenErr
	}

	thenBranch, thenErr := p.statement()
	if thenErr != nil {
		return nil, thenErr
	}

	var elseBranch ast.Stmt
	var elseErr error
	if p.match(lexer.ELSE) {
		elseBranch, elseErr = p.statement()
		if elseErr != nil {
			return nil, elseErr
		}
	}

	return &ast.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, nameErr := p.consume(lexer.IDENTIFIER, "Expect variable name.")
	if nameErr != nil {
		return nil, nameErr
	}

	var initializer ast.Expr
	var initializerErr error
	if p.match(lexer.EQUAL) {
		initializer, initializerErr = p.expression()
		if initializerErr != nil {
			return nil, initializerErr
		}
	}

	_, semiErr := p.consume(lexer.SEMICOLON, "Expect ';' after variable declaration.")
	if semiErr != nil {
		return nil, semiErr
	}

	return &ast.Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
	_, lParenErr := p.consume(lexer.LEFT_PAREN, "Expect '(' after 'while'.")
	if lParenErr != nil {
		return nil, lParenErr
	}

	cond, condErr := p.expression()
	if condErr != nil {
		return nil, condErr
	}

	_, rParenErr := p.consume(lexer.RIGHT_PAREN, "Expect ')' after condition.")
	if rParenErr != nil {
		return nil, rParenErr
	}

	body, bodyErr := p.statement()
	if bodyErr != nil {
		return nil, bodyErr
	}

	return &ast.While{
		Condition: cond,
		Body:      body,
	}, nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, semiErr := p.consume(lexer.SEMICOLON, "Expect ';' after value.")
	if semiErr != nil {
		return nil, semiErr
	}

	return &ast.Print{Expression: value}, nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()

	var value ast.Expr
	var valErr error
	if !p.check(lexer.SEMICOLON) {
		value, valErr = p.expression()
		if valErr != nil {
			return nil, valErr
		}
	}

	if _, err := p.consume(lexer.SEMICOLON, "Expect ';' after return value."); err != nil {
		return nil, err
	}

	return &ast.Return{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, semiErr := p.consume(lexer.SEMICOLON, "Expect ';' after expression.")
	if semiErr != nil {
		return nil, semiErr
	}

	return &ast.Expression{Expression: expr}, nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !p.check(lexer.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	_, err := p.consume(lexer.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) function(kind string) (*ast.Function, error) {
	name, nameErr := p.consume(lexer.IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	if nameErr != nil {
		return nil, nameErr
	}

	if _, err := p.consume(lexer.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind)); err != nil {
		return nil, err
	}

	parameters := []*lexer.Token{}

	if !p.check(lexer.RIGHT_PAREN) {
		// do-while loop
		for hasArg := true; hasArg; hasArg = p.match(lexer.COMMA) {
			// print too many arguments error but don't stop parsing
			if len(parameters) >= 255 {
				tooManyArgsErr := p.logError("Can't have more than 255 parameters.")
				fmt.Println(tooManyArgsErr)
			}

			arg, argErr := p.consume(lexer.IDENTIFIER, "Expect parameter name.")
			if argErr != nil {
				return nil, argErr
			}
			parameters = append(parameters, arg)
		}
	}

	if _, err := p.consume(lexer.RIGHT_PAREN, "Expect ')' after parameters."); err != nil {
		return nil, err
	}

	if _, err := p.consume(lexer.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind)); err != nil {
		return nil, err
	}

	body, bodyErr := p.block()
	if bodyErr != nil {
		return nil, bodyErr
	}

	return &ast.Function{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}
