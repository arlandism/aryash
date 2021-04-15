package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const exitPrompt = "Valar Dohaeris"
const entryPrompt = "Valar Morghulis\nType 'exit' at any time to quit the shell."
const CHANGE_DIR = "cd"
const EXIT = "exit"

func handleCommand(s string) error {
	parts := strings.Split(s, " ")
	switch parts[0] {
	case "":
		return nil
	case CHANGE_DIR:
		err := os.Chdir(parts[1])
		return err
	default:
		var stdout bytes.Buffer
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = &stdout
		err := cmd.Start()
		if err != nil {
			return err
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
		fmt.Printf(string(stdout.Bytes()))
		return nil
	}
}

var icon = string([]byte{0xF0, 0x9F, 0x92, 0x80})
var commandFlag = flag.String("c", "", "Run a command in a shell subprocess and then exit")

func main() {
	flag.Parse()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	if *commandFlag != "" {
		handleCommand(*commandFlag)
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(entryPrompt)
		for {
			fmt.Printf("%s ", icon)
			raw, err := reader.ReadString('\n')
			text := strings.TrimRight(raw, "\n")
			if err == io.EOF || text == EXIT {
				fmt.Println(exitPrompt)
				os.Exit(0)
			}
			err = handleCommand(text)
			if err != nil {
				fmt.Printf("can't process command: %s\n", text)
			}
		}
	}
}
