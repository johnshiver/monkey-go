package monkey_interpreter

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := NewLexer(line)
		for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
