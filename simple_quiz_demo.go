package main

import (
	"fmt"
	"wapuugotchi_cli/quiz"
)

func main() {
	fmt.Println("=== Go Modul 'quiz' - Exportierte Elemente ===\n")
	
	// Erstelle ein neues Quiz
	q := quiz.NewQuiz("Beispiel Quiz")
	
	// Zeige die Struktur
	fmt.Println("🔹 Hauptstrukturen:")
	fmt.Println("   • Quiz - Repräsentiert ein vollständiges Quiz")
	fmt.Println("   • Question - Repräsentiert eine einzelne Frage")
	
	// Zeige die exportierten Funktionen
	fmt.Println("\n🔹 Konstruktoren:")
	fmt.Println("   • NewQuiz(title string) *Quiz")
	fmt.Println("   • CreateSampleQuiz() *Quiz")
	
	fmt.Println("\n🔹 Quiz-Methoden:")
	fmt.Println("   • AddQuestion(text, options, correct) - Fügt eine Frage hinzu")
	fmt.Println("   • GetCurrentQuestion() - Gibt aktuelle Frage zurück")
	fmt.Println("   • AnswerQuestion(answer) - Beantwortet eine Frage")
	fmt.Println("   • IsFinished() - Prüft ob Quiz beendet ist")
	fmt.Println("   • GetScore() - Gibt Punktestand zurück")
	fmt.Println("   • GetProgress() - Gibt Fortschritt zurück")
	fmt.Println("   • Shuffle() - Mischt Fragenreihenfolge")
	fmt.Println("   • Reset() - Setzt Quiz zurück")
	fmt.Println("   • GetSummary() - Gibt Zusammenfassung zurück")
	
	// Demonstration der Verwendung
	fmt.Println("\n🔹 Beispiel-Verwendung:")
	
	// Eine Frage hinzufügen
	q.AddQuestion(
		"Was ist 2 + 2?",
		[]string{"3", "4", "5", "6"},
		1, // Index der richtigen Antwort (4)
	)
	
	fmt.Printf("Quiz erstellt: '%s'\n", q.Title)
	fmt.Printf("Anzahl Fragen: %d\n", len(q.Questions))
	
	// Aktuelle Frage abrufen
	if question, ok := q.GetCurrentQuestion(); ok {
		fmt.Printf("Aktuelle Frage: %s\n", question.Text)
		fmt.Printf("Antwortmöglichkeiten: %v\n", question.Options)
	}
	
	// Antwort geben
	correct := q.AnswerQuestion(1) // Antwort 1 (Index) = "4"
	fmt.Printf("Antwort richtig: %t\n", correct)
	
	// Score anzeigen
	score, total := q.GetScore()
	fmt.Printf("Aktueller Score: %d/%d\n", score, total)
	
	// Quiz-Status
	fmt.Printf("Quiz beendet: %t\n", q.IsFinished())
	fmt.Printf("Fortschritt: %.1f%%\n", q.GetProgress()*100)
	
	// Beispiel-Quiz verwenden
	fmt.Println("\n🔹 Vordefiniertes Beispiel-Quiz:")
	sampleQuiz := quiz.CreateSampleQuiz()
	fmt.Printf("Titel: %s\n", sampleQuiz.Title)
	fmt.Printf("Anzahl Fragen: %d\n", len(sampleQuiz.Questions))
	
	// Erste Frage des Beispiel-Quiz anzeigen
	if question, ok := sampleQuiz.GetCurrentQuestion(); ok {
		fmt.Printf("Erste Frage: %s\n", question.Text)
		for i, option := range question.Options {
			marker := "  "
			if i == question.Correct {
				marker = "→ "
			}
			fmt.Printf("%s%d) %s\n", marker, i+1, option)
		}
	}
	
	fmt.Println("\n✅ Quiz-Modul erfolgreich demonstriert!")
}