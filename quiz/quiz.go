package quiz

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
)

// Quiz represents one quiz question
type Quiz struct {
	Question string
	Options  []string
	Correct  int
	FeedbackCorrect   string
	FeedbackIncorrect string
}

var Quests []Quiz = getQuiz()

func GetRandomQuest() Quiz {
	index := rand.Intn(len(Quests))
	return Quests[index]
}



func getQuiz() []Quiz {
	resp, err := http.Get("https://raw.githubusercontent.com/wapuugotchi/wapuugotchi/refs/heads/main/inc/games/quiz/data/QuizWordPress.php")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	quizzes := parsePHP(string(body))
	if len(quizzes) == 0 {
		return nil
	}
	return quizzes
}

func parsePHP(php string) []Quiz {
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

			correct := lines[i+4]
			correctStart := strings.Index(correct, "__( '")
			correctEnd := strings.Index(correct, "', 'wapuugotchi' )")
			if correctStart != -1 && correctEnd != -1 {
				correct = correct[correctStart+5 : correctEnd]
			}

			feedbackCorrect := lines[i+5]
			feedbackCorrectStart := strings.Index(feedbackCorrect, "__( '")
			feedbackCorrectEnd := strings.Index(feedbackCorrect, "', 'wapuugotchi' )")
			if feedbackCorrectStart != -1 && feedbackCorrectEnd != -1 {
				feedbackCorrect = feedbackCorrect[feedbackCorrectStart+5 : feedbackCorrectEnd]
			}

			feedbackIncorrect := lines[i+6]
			feedbackIncorrectStart := strings.Index(feedbackIncorrect, "__( '")
			feedbackIncorrectEnd := strings.Index(feedbackIncorrect, "', 'wapuugotchi' )")
			if feedbackIncorrectStart != -1 && feedbackIncorrectEnd != -1 {
				feedbackIncorrect = feedbackIncorrect[feedbackIncorrectStart+5 : feedbackIncorrectEnd]
			}

			var options []string

			optionLine := lines[i+3]

			for len(optionLine) > 10 {
				start := strings.Index(optionLine, "__( '")
				end := strings.Index(optionLine, "', 'wapuugotchi' )")

				if start != -1 && end != -1 {
					optionText := optionLine[start+5 : end]
					options = append(options, optionText)
					optionLine = optionLine[end+18:]
				} else {
					break // Exit if pattern not found
				}
			}


			// Insert correct answer at random position
			insertPos := rand.Intn(len(options) + 1)
			options = append(options[:insertPos], append([]string{correct}, options[insertPos:]...)...)

	
			quizzes = append(quizzes, Quiz{
				Question: question,
				Options:  options,
				Correct: insertPos,
				FeedbackCorrect:   feedbackCorrect,
				FeedbackIncorrect: feedbackIncorrect,
			})
		}

	}

	return quizzes
}
