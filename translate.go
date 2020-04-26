package main

import (
	"fmt"

	gt "github.com/bas24/googletranslatefree"
)

// createTranslateMessage returns the input translated from inLang to outLang
func createTranslateMessage(text, inLang, outLang string) (translated string) {
	translated, err := gt.Translate(text, inLang, outLang)
	if err != nil {
		fmt.Println("Error translating,", err)
	}

	return
}
