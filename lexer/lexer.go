package lexer

import (
	"fmt"
	"strconv"
	"unicode"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

// Lexer holds the input and current scan position
type Lexer struct {
	source  string
	start   int
	current int
	line    int
	tokens  []Token
	Errors  []LexError
}

func New(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
		tokens: []Token{},
		Errors: []LexError{},
	}
}

func (l *Lexer) Lex() []Token {
	for !l.isAtEnd() {
		l.start = l.current
		ch := l.advance()

		switch {
		case ch != '\n' && isWhitespace(ch):
		case ch == '\n':
			l.line++
		case ch == '(':
			l.addToken(LEFT_PAREN, nil)
		case ch == ')':
			l.addToken(RIGHT_PAREN, nil)
		case ch == '{':
			l.addToken(LEFT_BRACE, nil)
		case ch == '}':
			l.addToken(RIGHT_BRACE, nil)
		case ch == ',':
			l.addToken(COMMA, nil)
		case ch == '.':
			l.addToken(DOT, nil)
		case ch == '+':
			l.addToken(PLUS, nil)
		case ch == '-':
			l.addToken(MINUS, nil)
		case ch == ';':
			l.addToken(SEMICOLON, nil)
		case ch == '*':
			l.addToken(STAR, nil)
		case ch == '!':
			if l.match('=') {
				l.addToken(BANG_EQUAL, nil)
			} else {
				l.addToken(BANG, nil)
			}
		case ch == '=':
			if l.match('=') {
				l.addToken(EQUAL_EQUAL, nil)
			} else {
				l.addToken(EQUAL, nil)
			}
		case ch == '<':
			if l.match('=') {
				l.addToken(LESS_EQUAL, nil)
			} else {
				l.addToken(LESS, nil)
			}
		case ch == '>':
			if l.match('=') {
				l.addToken(GREATER_EQUAL, nil)
			} else {
				l.addToken(GREATER, nil)
			}
		case ch == '/':
			if l.match('/') {
				l.skipComment()
			} else {
				l.addToken(SLASH, nil)
			}
		case ch == '"':
			l.lexString()
		case isDigit(ch):
			l.lexNumber()
		case isAlpha(ch):
			l.lexIdentifier()
		default:
			// throw unknown token error and break
			l.logError(fmt.Sprintf("Unexpected character: %c", ch))
		}
	}

	l.tokens = append(l.tokens, Token{Type: EOF, Lexeme: "", Literal: nil})
	return l.tokens
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() rune {
	ch := l.source[l.current]
	l.current++
	return rune(ch)
}

func (l *Lexer) addToken(t TokenType, literal any) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, Token{
		Type:    t,
		Lexeme:  text,
		Literal: literal,
	})
}

func (l *Lexer) match(expected rune) bool {
	if l.isAtEnd() || l.peek() != expected {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return 0
	}
	return rune(l.source[l.current])
}

func (l *Lexer) peekNext() rune {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return rune(l.source[l.current+1])
}

func (l *Lexer) skipComment() {
	for l.peek() != '\n' && !l.isAtEnd() {
		l.advance()
	}
}

func (l *Lexer) lexString() {
	for l.peek() != '"' && !l.isAtEnd() {
		if l.peek() == '\n' {
			// allow multiline string, inc line count
			l.line++
		}
		l.advance()
	}

	if l.isAtEnd() {
		l.logError("Unterminated string.")
		return
	}

	// closing "
	l.advance()
	value := l.source[l.start+1 : l.current-1]
	l.addToken(STRING, value)
}

func (l *Lexer) lexNumber() {
	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance()
		for !l.isAtEnd() && isDigit(l.peek()) {
			l.advance()
		}
	}

	value, err := strconv.ParseFloat(l.source[l.start:l.current], 64)
	if err != nil {
		l.logError("Number cannot be parsed to float64")
		return
	}
	l.addToken(NUMBER, value)
}

func (l *Lexer) lexIdentifier() {
	for !l.isAtEnd() && isAlphaNumeric(l.peek()) {
		l.advance()
	}
	text := l.source[l.start:l.current]
	tokenType := IDENTIFIER
	if reserved, ok := keywords[text]; ok {
		tokenType = reserved
	}
	l.addToken(tokenType, nil)
}

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isAlpha(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isAlphaNumeric(ch rune) bool {
	return isAlpha(ch) || isDigit(ch)
}

func (l *Lexer) logError(message string) {
	l.Errors = append(l.Errors, LexError{
		line:    l.line,
		message: message,
	})
}
