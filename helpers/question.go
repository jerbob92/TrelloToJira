package helpers

import (
	"fmt"
)

func AskForSelect(question string, options []string) string {
	fmt.Println(question)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return AskForSelect(question, options)
	}

	for _, option := range options {
		if response == option {
			return option;
		}
	}

	return AskForSelect(question, options)
}

func AskForConfirmation(question string, default_option *bool) bool {
	fmt.Println(question)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		if default_option != nil {
			return *default_option
		}
		return AskForConfirmation(question, default_option)
	}

	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		if default_option != nil && response == "" {
			return *default_option
		}
		return AskForConfirmation(question, default_option)
	}
}

// You might want to put the following two functions in a separate utility package.

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}