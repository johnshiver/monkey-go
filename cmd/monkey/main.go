package main

import (
	"fmt"
	monkey "monkey-interpreter"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey programming language!\n", usr.Username)
	fmt.Printf("Feel free to type in commands\n")
	monkey.Start(os.Stdin, os.Stdout)
}
