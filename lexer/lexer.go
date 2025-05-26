package lexer

import (
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
			l.addToken(LEFT_PAREN, "")
		case ch == ')':
			l.addToken(RIGHT_PAREN, "")
		case ch == '{':
			l.addToken(LEFT_BRACE, "")
		case ch == '}':
			l.addToken(RIGHT_BRACE, "")
		case ch == ',':
			l.addToken(COMMA, "")
		case ch == '.':
			l.addToken(DOT, "")
		case ch == '+':
			l.addToken(PLUS, "")
		case ch == '-':
			l.addToken(MINUS, "")
		case ch == ';':
			l.addToken(SEMICOLON, "")
		case ch == '*':
			l.addToken(STAR, "")
		case ch == '!':
			if l.match('=') {
				l.addToken(BANG_EQUAL, "")
			} else {
				l.addToken(BANG, "")
			}
		case ch == '=':
			if l.match('=') {
				l.addToken(EQUAL_EQUAL, "")
			} else {
				l.addToken(EQUAL, "")
			}
		case ch == '<':
			if l.match('=') {
				l.addToken(LESS_EQUAL, "")
			} else {
				l.addToken(LESS, "")
			}
		case ch == '>':
			if l.match('=') {
				l.addToken(GREATER_EQUAL, "")
			} else {
				l.addToken(GREATER, "")
			}
		case ch == '/':
			if l.match('/') {
				l.skipComment()
			} else {
				l.addToken(SLASH, "")
			}
		case ch == '"':
			l.lexString()
		case unicode.IsDigit(ch):
			l.lexNumber()
		case unicode.IsLetter(ch):
			l.lexIdentifier()
		default:
			// throw unknown token error and break
			l.tokenError(ch)
		}
	}

	l.tokens = append(l.tokens, Token{Type: EOF, Lexeme: "", Literal: ""})
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

func (l *Lexer) addToken(t TokenType, literal string) {
	text := l.source[l.start:l.current]
	l.tokens = append(l.tokens, Token{
		Type:    t,
		Lexeme:  text,
		Literal: literal,
	})
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
		// TODO: add error reporting for this error
		panic("Unterminated string")
		return
	}

	// closing "
	l.advance()
	value := l.source[l.start+1 : l.current-1]
	l.addToken(STRING, value)
}

func (l *Lexer) lexNumber() {
	for unicode.IsDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == '.' && unicode.IsDigit(l.peekNext()) {
		l.advance()
		for !l.isAtEnd() && unicode.IsDigit(l.peek()) {
			l.advance()
		}
	}

	value := l.source[l.start:l.current]
	// value, err := strconv.ParseFloat(l.source[l.start:l.current], 64)
	// if err != nil {
	// 	panic("Float parser error")
	// }
	l.addToken(NUMBER, value)
}

func (l *Lexer) lexIdentifier() {
	for !l.isAtEnd() && unicode.IsLetter(l.peek()) {
		l.advance()
	}
	text := l.source[l.start:l.current]
	tokenType := IDENTIFIER
	if reserved, ok := keywords[text]; ok {
		tokenType = reserved
	}
	l.addToken(tokenType, "")
}

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func (l *Lexer) tokenError(ch rune) {
	l.Errors = append(l.Errors, LexError{
		line: l.line,
		ch:   ch,
	})
}
