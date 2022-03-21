package monkey_interpreter

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := Lexer{input: input}
	l.readChar()
	return &l
}

// readChar sets char to current read position and advances lexer cursor
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// peekChar is like readChar, but doesnt increment reader
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	startPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.ch {
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '+':
		tok = newToken(PLUS, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '!':
		// could be != or !
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: NOT_EQ, Literal: literal}
		} else {
			tok = newToken(BANG, l.ch)
		}
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '<':
		tok = newToken(LT, l.ch)
	case '>':
		tok = newToken(GT, l.ch)
	case '=':
		// could be == or =
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: EQ, Literal: literal}
		} else {
			tok = newToken(ASSIGN, l.ch)
		}
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		// NOTE: read ident + read number dont need to readChar
		// so can return immediately

		// not at a token, so we need to do some read logic
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		}

		// TODO: how does integer overflow work?
		if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(ILLEGAL, l.ch)
	}
	// after token read is finished, advance cursor
	l.readChar()
	return tok
}
