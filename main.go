package main

import (
	"bufio"
	"fmt"
	"github.com/xdarksyderx/pokeshell/pokeshell"
	"os"
	"strings"
)

func main() {
	pS := pokeshell.CreatePokeshell()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("pokeshell> ")
		input, readingError := reader.ReadString('\n')
		if readingError != nil {
			fmt.Printf("Error reading command: %s", readingError)
			continue
		}
		parsedInput := strings.Split(strings.TrimSpace(input), " ")
		if cmd, exists := pS.CommandList[parsedInput[0]]; exists {
			commandError := cmd.Callback(&pS, parsedInput[1:])
			if commandError != nil {
				fmt.Println(commandError)
			}
		} else {
			fmt.Println("Unknown command:", parsedInput)
		}
	}
}
