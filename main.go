package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"wapuugotchicli/quiz"
	"os/exec"
	"runtime"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	clearTerminal()
	fmt.Println("Welcome to the WapuuGotchi Quiz!")
	fmt.Println("Answer the questions by typing the number of your choice and pressing Enter.")
	fmt.Println("Press Enter without typing anything to exit.\n")
	fmt.Println("To start the game, press Enter to start the quiz.")
	fmt.Println()
	
	reader.ReadString('\n')

	var correctAnswers int
	var totalQuestions int

	for {
		clearTerminal()

		fmt.Printf("Score: %d/%d\n\n", correctAnswers, totalQuestions)

		totalQuestions++
		
		var quest = quiz.GetRandomQuest()
		fmt.Println(quest.Question)

		for i, option := range quest.Options {
			fmt.Printf("%d) %s\n", i+1, option)
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Check if input is a valid number
		if input == "" {

			fmt.Println("Good bye.")
			return
		}

		switch input {
		case fmt.Sprintf("%d", quest.Correct+1):
			fmt.Println("You are correct! Yeah, great job! 🎉")
			correctAnswers++
		default:
			fmt.Println("Sorry, that's incorrect. That is sad. 😢")
		}

		time.Sleep(1 * time.Second)
	}
}

func clearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default: // linux, macOS, etc.
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
