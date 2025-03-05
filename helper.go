package main

import (
	"fmt"
	"math/rand"
	"unicode"

	"github.com/pterm/pterm"
	"github.com/spf13/cast"
)

func getRandomLetterLowercase() rune {
	defaultRune := rune('a')
	defaultRune = rune(int(defaultRune) + rand.Intn(26))

	return defaultRune
}

func getRandomLetterUppercase() rune {
	defaultRune := rune('A')
	defaultRune = rune(int(defaultRune) + rand.Intn(26))

	return defaultRune
}

func getRandomNumber() rune {
	defaultRune := rune('0')
	defaultRune = rune(int(defaultRune) + rand.Intn(10))

	return defaultRune
}

func getRandomSpecialChar() rune {
	specialChars := []rune{'!', '@', '#', '$', '%', '&', '*'}
	return specialChars[rand.Intn(len(specialChars))]
}

func getRandomCharacter(withNumbers, withSpecial bool) rune {
	ranges := []rune{'L', 'U'}
	if withNumbers {
		ranges = append(ranges, 'N')
	}
	if withSpecial {
		ranges = append(ranges, 'S')
	}

	switch ranges[rand.Intn(len(ranges))] {
	default:
		fallthrough
	case 'L':
		return getRandomLetterLowercase()
	case 'U':
		return getRandomLetterUppercase()
	case 'N':
		return getRandomNumber()
	case 'S':
		return getRandomSpecialChar()
	}
}

func selectStartsWith() (string, error) {
	options := []string{"Lowercase", "Uppercase", "Number", "Special Character"}
	return pterm.DefaultInteractiveSelect.WithDefaultText("Your password should starts with?").WithOptions(options).Show()
}

func selectEndsWith() (string, error) {
	options := []string{"Lowercase", "Uppercase", "Number", "Special Character"}
	return pterm.DefaultInteractiveSelect.WithDefaultText("Your password should ends with?").WithOptions(options).Show()
}

func inputYesOrNo(text string) (bool, error) {
	return pterm.DefaultInteractiveConfirm.WithDefaultText(text).Show()
}

func inputPasswordSize() (int, error) {
	var input string
	var err error

	for {
		input, err = pterm.DefaultInteractiveTextInput.WithDefaultText("New password size (type 'exit' to cancel)").Show()
		if err != nil {
			return 0, err
		}

		if input == "exit" {
			return 0, fmt.Errorf("user exited")
		}

		var valid bool = true
		for _, v := range input {
			if !unicode.IsDigit(v) {
				pterm.Error.Println("Please enter a number")
				valid = false
				break
			}
		}

		if valid {
			break
		}
	}

	return cast.ToInt(input), nil
}

func generatePassword(password *password) string {
	var characters = make([]rune, password.Size)

	switch password.StartsWith {
	default:
		fallthrough
	case "Lowercase":
		characters = append(characters, rune(getRandomLetterLowercase()))
	case "Uppercase":
		characters = append(characters, rune(getRandomLetterUppercase()))
	case "Number":
		characters = append(characters, rune(getRandomNumber()))
	case "Special Character":
		characters = append(characters, rune(getRandomSpecialChar()))
	}

	for i := 1; i < password.Size-1; i++ {
		characters = append(characters, getRandomCharacter(password.WithNumbers, password.WithSpecialChar))
	}

	switch password.EndsWith {
	default:
		fallthrough
	case "Special Character":
		characters = append(characters, rune(getRandomSpecialChar()))
	case "Lowercase":
		characters = append(characters, rune(getRandomLetterLowercase()))
	case "Uppercase":
		characters = append(characters, rune(getRandomLetterUppercase()))
	case "Number":
		characters = append(characters, rune(getRandomNumber()))
	}

	return string(characters)
}
