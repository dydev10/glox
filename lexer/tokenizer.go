package lexer

import (
	"bufio"
	"io"
	"os"
)

func Tokenize(filename string) ([]Token, error) {
	result := make([]Token, 0)
	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				return result, err
			}
			result = append(result, Token{EOF, "", ""})
			break
		}
		switch b {
		case '(':
			result = append(result, Token{LEFT_PAREN, "(", ""})
		case ')':
			result = append(result, Token{RIGHT_PAREN, ")", ""})
		}
	}
	return result, nil
}
