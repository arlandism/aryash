package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var exitPrompt = "Valar Dohaeris"
var prompt = string([]byte{0xF0, 0x9F, 0x92, 0x80})

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Valar Morghulis")

	for {
		fmt.Printf("%s ", prompt)
		text, err := reader.ReadString('\n')
		if text == "exit\n" {
			os.Exit(0)
		}
		if err == io.EOF {
			fmt.Println(exitPrompt)
			os.Exit(0)
		}
		fmt.Printf(text)
	}
}
