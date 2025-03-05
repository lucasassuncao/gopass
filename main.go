package main

import (
	"fmt"
)

type password struct {
	Size            int
	StartsWith      string
	EndsWith        string
	WithNumbers     bool
	WithSpecialChar bool
}

func main() {
	var withNumbers, withSpecialChar bool

	size, err := inputPasswordSize()
	if err != nil {
		fmt.Println("Error while getting password size:", err)
		return
	}

	if size < 4 {
		fmt.Println("Mininum password size is 4")
		return
	}

	startsWith, _ := promptCharacterType("Your password should starts with?")
	endsWith, _ := promptCharacterType("Your password should ends with?")

	withNumbers, _ = promptYesOrNo("Should the password contain numbers?")
	withSpecialChar, _ = promptYesOrNo("Should the password contain special characters?")

	password := &password{
		Size:            size,
		StartsWith:      startsWith,
		EndsWith:        endsWith,
		WithNumbers:     withNumbers,
		WithSpecialChar: withSpecialChar,
	}

	generatedPassword, err := generatePassword(password)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = hasAllRequiredCharacters(generatedPassword, password)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = hasExpectedLength(generatedPassword, password.Size)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nGenerated password:", generatedPassword)
}
