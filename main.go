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
	size, err := inputPasswordSize()
	if err != nil {
		fmt.Println("Error while getting password size:", err)
		return
	}

	if size == 0 {
		fmt.Println("Password size is 0, exiting...")
		return
	}

	startsWith, _ := selectStartsWith()
	endsWith, _ := selectEndsWith()

	withNumbers, _ := inputYesOrNo("Should the password contain numbers?")
	withSpecialChar, _ := inputYesOrNo("Should the password contain special characters?")

	password := &password{
		Size:            size,
		StartsWith:      startsWith,
		EndsWith:        endsWith,
		WithNumbers:     withNumbers,
		WithSpecialChar: withSpecialChar,
	}

	generatedPassword := generatePassword(password)

	fmt.Println("Generated password:", generatedPassword)
}
