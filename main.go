package main

import (
	"fmt"
	"os"
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

	startsWith, _ := selectStartsWith()
	endsWith, _ := selectEndsWith()

	withNumbers, _ = inputYesOrNo("Should the password contain numbers?")
	withSpecialChar, _ = inputYesOrNo("Should the password contain special characters?")

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
		os.Exit(1)
	}

	fmt.Println("\nGenerated password:", generatedPassword)
}
