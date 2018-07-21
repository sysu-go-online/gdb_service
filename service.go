package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// GenerateMakefileContent generates generic Makefile is the `Debug/` directory
func GenerateMakefileContent(pn string) {
	// TODO: generate makefile content
	var makefile string

	// create Debug/temp folder if not exists
	err := os.MkdirAll("Debug/temp", os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// save makefile file
	ioutil.WriteFile("Debug/Makefile", []byte(makefile), os.ModeAppend)
}
