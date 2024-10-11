package itop

import (
	"bufio"
	"fmt"
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lxr := lexer.New(line)
		parsr := parser.New(lxr)
		program := parsr.ParseProgram()

		if len(parsr.Errors()) != 0 {
			printParserErrors(out, parsr.Errors())
			continue
		}

		fmt.Fprintf(out, "%+v\n", program.String())

		//for tkn := lxr.NextToken(); tkn.Type != token.EOF; tkn = lxr.NextToken() {
		//	fmt.Fprintf(out, "%+v\n", tkn)
		//}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	fmt.Println(" parser errors:")
	for _, msg := range errors {
		fmt.Fprintln(out, "\t"+msg)
	}
}
