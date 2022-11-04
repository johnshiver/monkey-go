package monkey_interpreter

import (
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		_, err := fmt.Fprintf(out, PROMPT)
		if err != nil {
			return
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := NewLexer(line)
		p := NewParser(l)
		program := p.ParseProgram()
		if len(p.errors) != 0 {
			printParserErrors(out, p.errors)
			continue
		}
		evaluated := Eval(program)
		if evaluated != nil {
			_, err := io.WriteString(out, evaluated.Inspect())
			if err != nil {
				return
			}
			_, err = io.WriteString(out, "\n")
			if err != nil {
				return
			}
		}
	}
}
func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		_, err := io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			return
		}
	}
}
