package cli

import (
	"fmt"
	"strings"

	"github.com/shrimpsizemoose/trekker/logger"
)

func ConfirmAction(message string) bool {
	logger.Question.Println(message)
	logger.Question.Println("Погнали? y/n/yes/no/д/н/да/нет")

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		logger.Error.Fatal(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes", "д", "да":
		return true
	case "n", "no", "н", "нет":
		return false
	default:
		fmt.Println("Я могу понять, только если ввести y/n/yes/no/д/н/да/нет, вот такой я дурачок.")
		return ConfirmAction(message)
	}
}

func SetupUsage(usageMessage string) {
	fmt.Println(usageMessage)
}
