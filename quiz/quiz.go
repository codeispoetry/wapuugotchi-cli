package quiz

import (
	"io"
	"net/http"
	"math/rand"
	"strings"
)

// Quiz represents one quiz question
type Quiz struct {
	Question      string
	Options       []string
}

func Display() string {
	return "\n" + getData()
}

func getData() string {
	resp, err := http.Get("https://raw.githubusercontent.com/wapuugotchi/wapuugotchi/refs/heads/main/inc/games/quiz/data/QuizWordPress.php")
	if err != nil {
		return "Error fetching quiz data: " + err.Error()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error reading quiz data: " + err.Error()
	}

	quizzes := parseQuiz(string(body))

	var result strings.Builder

	randomIndex := rand.Intn(len(quizzes))
	randomIndex = 0 // for testing purposes, always pick the first quiz
	result.WriteString(quizzes[randomIndex].Question)
	for i, option := range quizzes[randomIndex].Options {
		result.WriteString("\n")
		result.WriteString(string('A' + i))
		result.WriteString(". ")
		result.WriteString(option)
	}


	result.WriteString("\n")
	
	return result.String()
}

// parsePHPQuiz extracts quiz questions from a PHP string
func parseQuiz(php string) []Quiz {
	var quizzes []Quiz

	lines := strings.Split(php, "\n")
	for i, line := range lines {
		if strings.Contains(line, "$quiz[] = new Quiz(") {
			// Simulate parsing quiz data
			// In a real implementation, you would extract the actual data here

			question := lines[i+2]
			// Extract question text between first and second single quote
			if strings.Count(question, "'") >= 2 {
				firstQuote := strings.Index(question, "'")
				secondQuote := strings.Index(question[firstQuote+1:], "'")
				if secondQuote != -1 {
					question = question[firstQuote+1 : firstQuote+1+secondQuote]
				}
			}


			var options []string
			optionLine := lines[i+3]
			
			
			for len(optionLine) > 60 {
				start := strings.Index(optionLine, "__( '")
				end := strings.Index(optionLine, "', 'wapuugotchi' )")
				
				if start != -1 && end != -1 {
					optionText := optionLine[start+5 : end]
					options = append(options, optionText)
					optionLine = optionLine[end+18:] // Skip past the closing pattern
				} else {
					break // Exit if pattern not found
				}
			}


			quizzes = append(quizzes, Quiz{
				Question:      question,
				Options:       options,
			})
		}
	
	}

		
	

	return quizzes
}