package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var exitPrompt = "Valar Dohaeris"
var entryPrompt = "Valar Morghulis\nType 'exit' at any time to quit the shell."
var icon = string([]byte{0xF0, 0x9F, 0x92, 0x80})
var commandFlag = flag.String("c", "", "Run a command in a shell subprocess and then exit")

func handleCommand(s string) error {
	parts := strings.Split(s, " ")
	out, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func main() {
	flag.Parse()
	if *commandFlag != "" {
		handleCommand(*commandFlag)
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(entryPrompt)
		for {
			fmt.Printf("%s ", icon)
			raw, err := reader.ReadString('\n')
			text := strings.TrimRight(raw, "\n")
			if err == io.EOF || text == "exit" {
				fmt.Println(exitPrompt)
				os.Exit(0)
			}
			if text == "" {
				continue
			}
			err = handleCommand(text)
			if err != nil {
				fmt.Printf("can't process command: %s\n", text)
			}
		}
	}
}
