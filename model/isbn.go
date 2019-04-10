package model

import (
	"regexp"
)

const ISBNRegex = `(ISBN|isbn)?(-13|-10)?[ .:]*(97.)?([-]?[\dX]){10}`

func FindISBN(input string) []string {

	compile := regexp.MustCompile(ISBNRegex)

	return compile.FindAllString(input, -1)
}

func ReformatISBN(input string) string {

	compile := regexp.MustCompile(`(ISBN|isbn|ISBN-13|ISBN-10|isbn-10|isbn-13)?([-: ])?`)

	i := compile.ReplaceAllString(input, "")

	return i;
}
