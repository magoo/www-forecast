package main

import (
	"encoding/csv"
	"fmt"
	"github.com/go-errors/errors"
	_ "github.com/joho/godotenv/autoload"
	"github.com/magoo/www-forecast/app/models"
	"io"
	"os"
)

func printError(err error) {
	fmt.Println("Error reading file:", errors.Wrap(err, 2).ErrorStack())
}

func main() {
	fmt.Println(os.Getenv("E6E_AWS_ENV"))

	preexistingQuestions := models.GetAllCalibrationQuestions()

	f, _ := os.Open("./scripts/mock_data.csv")
	defer f.Close()

	r := csv.NewReader(f)
	//var columnNames []string
	columnNames, err := r.Read()
	if err != nil {
		printError(err)
		return
	}
	fmt.Println(columnNames)
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			printError(err)
			return
		}

		//timestamp := line[0]
		description := line[1]
		correctAnswerString := line[2]
		if correctAnswerString != "true" && correctAnswerString != "false" {
			fmt.Println("Error: Row has invalid true/false value:", line)
			return
		}
		correctAnswer := correctAnswerString == "true"
		answerDetail := line[3]
		answerSource := line[4]
		difficulty := line[5]
		uniqueId := line[6]
		updateQuestion := false
		for i := 0; i < len(preexistingQuestions); i++ {
			q := preexistingQuestions[i]
			//fmt.Println("Checking pre-existing question: ")
			//fmt.Println("  description:", description, q.Description)
			//fmt.Println("  correctAnswer:", correctAnswer, q.CorrectAnswer)
			if uniqueId == q.Id {
				fmt.Println("Overwriting preexisting question: ")
				fmt.Println("  description:", description)
				fmt.Println("  correctAnswer:", correctAnswer)
				updateQuestion = true
				break
			}
		}
		if updateQuestion {
			models.UpdateCalibrationQuestion(uniqueId, description, correctAnswer, answerDetail, answerSource, difficulty)
		} else {
			models.CreateCalibrationQuestion(uniqueId, description, correctAnswer, answerDetail, answerSource, difficulty)
		}


		fmt.Println(line)
	}
}
