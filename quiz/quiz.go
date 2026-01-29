package quiz

import (
	"io"
	"net/http"
	"regexp"
	"fmt"
	"strings"
)

// Quiz represents one quiz question
type Quiz struct {
	ID            string
	Question      string
	Options       []string
	CorrectAnswer string
	CorrectText   string
	IncorrectText string
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

	quizzes := parsePHPQuiz(string(body))

	var result strings.Builder
	for _, quiz := range quizzes {
		result.WriteString("Question: " + quiz.Question + "\n")
		for i, option := range quiz.Options {
			result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, option))
		}
		result.WriteString("\n")
	}

	return result.String()
}

// parsePHPQuiz extracts quiz questions from a PHP string
func parsePHPQuiz(php string) []Quiz {
	var quizzes []Quiz

	// Regex to capture new Quiz(...) calls
	reQuiz := regexp.MustCompile(`new Quiz\(\s*'([^']+)',\s*__\(\s*'([^']+)'[^)]*\),\s*array\(([^)]*)\),\s*__\(\s*'([^']+)'[^)]*\),\s*__\(\s*'([^']+)'[^)]*\),\s*__\(\s*'([^']+)'[^)]*\)\s*\)`)

	matches := reQuiz.FindAllStringSubmatch(php, -1)

	for _, m := range matches {
		if len(m) != 7 {
			continue
		}

		// m[1]=ID, m[2]=question, m[3]=options, m[4]=correctAnswer, m[5]=correctText, m[6]=incorrectText
		opts := strings.Split(m[3], ",")
		for i := range opts {
			opts[i] = strings.TrimSpace(strings.Trim(opts[i], `__() '"))`))
		}

		quiz := Quiz{
			ID:            m[1],
			Question:      m[2],
			Options:       opts,
			CorrectAnswer: m[4],
			CorrectText:   m[5],
			IncorrectText: m[6],
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes
}