package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

// prompts
const (
	PromptEntry = "Valar Morghulis\nType 'exit' at any time to quit the shell."
	PromptExit  = "Valar Dohaeris"
)

// built-ins
const (
	CommandExit      = "exit"
	CommandChangeDir = "cd"
)

// currentCmd is set by 'handleCommand' and is leveraged by the signal handler
// to determine whether the signal should be propagated to a subprocess or
// the foreground shell.
// I don't love this shared state approach, so may rethink it.
var currentCmd *exec.Cmd = nil

func signalHandler(sigChan <-chan os.Signal) {
	for {
		sig := <-sigChan
		if currentCmd != nil { // kill the currently executing subprocess
			currentCmd.Process.Signal(sig)
		} else { // kill the foreground
			os.Exit(1)
		}
	}
}

func handleCommand(s string) error {
	parts := strings.Split(s, " ")
	switch parts[0] {
	case "":
		return nil
	case CommandChangeDir:
		err := os.Chdir(parts[1])
		return err
	default:
		cmd := exec.Command(parts[0], parts[1:]...)
		currentCmd = cmd
		out, err := cmd.Output()
		currentCmd = nil
		// -1 means either the user sent an interrupt signal or the child process isn't done yet
		// In our case it *must* be the former since using 'Output' waits. So if we reach
		// this line then execution has stopped.
		// There's a defensive cond check just in case there's something I'm missing.
		if cmd.ProcessState.ExitCode() == -1 {
			if cmd.ProcessState == nil {
				fmt.Println("Something went wrong. This shouldn't happen.")
			}
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf(string(out))
		return nil
	}
}

var icon = string([]byte{0xF0, 0x9F, 0x92, 0x80})
var commandFlag = flag.String("c", "", "Run a command in a shell subprocess and then exit")

func main() {
	flag.Parse()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go signalHandler(sigChan)
	if *commandFlag != "" {
		handleCommand(*commandFlag)
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(PromptEntry)
		for {
			fmt.Printf("%s ", icon)
			raw, err := reader.ReadString('\n')
			text := strings.TrimRight(raw, "\n")
			if err == io.EOF || text == CommandExit {
				fmt.Println(PromptExit)
				os.Exit(0)
			}
			err = handleCommand(text)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("can't process command: %s\n", text)
			}
		}
	}
}
