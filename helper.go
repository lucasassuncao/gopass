package main

import (
	"errors"
	"fmt"
	"math/rand"
	"unicode"
	"unicode/utf8"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/pterm/pterm"
	"github.com/spf13/cast"
)

var specialChars = []rune{'!', '@', '#', '$', '%', '&', '*'}
var typs = []string{"Lowercase", "Uppercase", "Number", "Special Character"}

// getRune returns a rune based on the specified charType.
// It defaults to returning a random letter.
func getRune(charType string) rune {
	switch charType {
	default:
		fallthrough
	case "Lowercase":
		return rune('a') + rand.Int31n(26)
	case "Uppercase":
		return rune('A') + rand.Int31n(26)
	case "Number":
		return rune('0') + rand.Int31n(10)
	case "Special Character":
		return specialChars[rand.Intn(len(specialChars))]
	}
}

// getRandomRune returns a random rune based on the specified options.
// The function generates a random number to determine the character type:
// - 70% chance to return a random lowercase or uppercase letter.
// - 20% chance to return a random number if withNumbers is true.
// - 10% chance to return a random special character if withSpecial is true.
// If the conditions for numbers or special characters are not met, it defaults to returning a random letter.
func getRandomRune(password *password) rune {
	roll := rand.Intn(100)

	switch {
	case roll < 70:
		if rand.Intn(2) == 0 {
			return rune('a') + rand.Int31n(26)
		}
		return rune('A') + rand.Int31n(26)
	case roll < 90 && password.WithNumbers:
		return rune('0') + rand.Int31n(10)
	case roll < 100 && password.WithSpecialChar:
		return specialChars[rand.Intn(len(specialChars))]
	default:
		if rand.Intn(2) == 0 {
			return rune('a') + rand.Int31n(26)
		}
		return rune('A') + rand.Int31n(26)
	}
}

// hasAllRequiredCharacters checks if a given string meets the password requirements.
// It verifies the presence of lowercase letters, uppercase letters, numbers, and special characters.
// The function returns an error message if any required character type is missing.
func hasAllRequiredCharacters(s string, password *password) error {
	hasLower, hasUpper, hasNumber, hasSpecial := false, false, false, false

	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasLower {
		return errors.New("your password doesn't have lowercase characters")
	}

	if !hasUpper {
		return errors.New("your password doesn't have upper characters")
	}

	if password.WithNumbers && !hasNumber {
		return errors.New("your password doesn't have numbers")
	}

	if password.WithSpecialChar && !hasSpecial {
		return errors.New("your password doesn't have special characters")
	}

	return nil
}

// hasExpectedLength checks if the given string has the expected length.
// It counts the number of runes (characters) in the string and compares it to the expected size.
func hasExpectedLength(s string, expectedSize int) error {
	currentSize := utf8.RuneCountInString(s)

	if currentSize != expectedSize {
		return fmt.Errorf("your password %s is %d characters long... expected %d characters long", s, currentSize, expectedSize)
	}

	return nil
}

// promptStartCharacterType prompts the user to select the type of character the password should start with.
func promptStartCharacterType() (string, error) {
	return pterm.DefaultInteractiveSelect.WithDefaultText("Your password should starts with?").WithOptions(typs).Show()
}

// promptEndCharacterType prompts the user to select the type of character the password should end with.
func promptEndCharacterType() (string, error) {
	return pterm.DefaultInteractiveSelect.WithDefaultText("Your password should ends with?").WithOptions(typs).Show()
}

// promptYesOrNo prompts the user with a Yes/No question using an interactive confirmation.
func promptYesOrNo(text string) (bool, error) {
	return pterm.DefaultInteractiveConfirm.WithDefaultText(text).Show()
}

// inputPasswordSize prompts the user to enter the desired password size using an interactive text input.
// It ensures that the user provides a valid numeric input.
func inputPasswordSize() (int, error) {
	var input string
	var err error

	for {
		input, err = pterm.DefaultInteractiveTextInput.WithDefaultText("What's the size of the new password?").Show()
		if err != nil {
			return 0, err
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

// generatePassword generates a password based on the provided `password` struct.
// It ensures that the password meets the required character types (lowercase, uppercase, number, special character)
// and follows specific constraints such as the starting and ending characters, as well as the desired length.
// The function ensures that the password contains at least one of each required character type, fills the remaining
// characters randomly, and shuffles the middle characters to avoid predictable patterns.
// It also checks that the generated password meets the length and character requirements before returning it.
func generatePassword(password *password) (string, error) {
	var characters = make([]rune, 0, password.Size)
	var remainingChars int = password.Size - 2

	set := hashset.New("Lowercase", "Uppercase", "Number", "Special Character")
	set.Remove(password.StartsWith)
	set.Remove(password.EndsWith)

	// sets the first character
	characters = append(characters, getRune(password.StartsWith))

	if !password.WithNumbers {
		set.Remove("Number")
	}

	if !password.WithSpecialChar {
		set.Remove("Special Character")
	}

	// Ensures at least one of each required type
	for _, value := range set.Values() {
		characters = append(characters, getRune(value.(string))) // explicitly casting to string
		remainingChars--
	}

	// Fills the remaining characters randomly
	for i := 1; i <= remainingChars; i++ {
		characters = append(characters, getRandomRune(password))
	}

	// sets the last character
	characters = append(characters, getRune(password.EndsWith))

	// Shuffles the characters in the middle of the password
	// First and last characters are preserved
	middle := characters[1 : len(characters)-1]
	rand.Shuffle(len(middle), func(i, j int) {
		middle[i], middle[j] = middle[j], middle[i]
	})

	// Reassemble the characters slice with shuffled middle
	characters = append([]rune{characters[0]}, append(middle, characters[len(characters)-1])...)

	finalPassword := string(characters)

	err := hasAllRequiredCharacters(finalPassword, password)
	if err != nil {
		return "", err
	}

	err = hasExpectedLength(finalPassword, password.Size)
	if err != nil {
		return "", err
	}

	return finalPassword, nil
}
