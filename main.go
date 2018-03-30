// Package main contains common and useful utils for the Go project development.
//
// Inclusion criteria:
//  - Only rely on the Go standard package
//  - Functions or lightweight packages
//  - Non-business related general tools
package main

import (
	"fmt"
	"log"

	"github.com/aljiwala/gutils/filepathx"
	"github.com/aljiwala/gutils/timex"
	"github.com/aljiwala/gutils/timex/monthx"
)

func main() {
	fmt.Println(timex.EndOfMonth(2018, monthx.January.TimeMonth()))
	lines, err := filepathx.GrepFile("-i TimeMonth()", "/home/mohsin/Work/Projects/go1.8/src/github.com/aljiwala/gutils/main.go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(lines)

	// // We'll start with a simple command that takes no
	// // arguments or input and just prints something to
	// // stdout. The `exec.Command` helper creates an object
	// // to represent this external process.
	// dateCmd := exec.Command("date")
	//
	// // `.Output` is another helper than handles the common
	// // case of running a command, waiting for it to finish,
	// // and collecting its output. If there were no errors,
	// // `dateOut` will hold bytes with the date info.
	// dateOut, err := dateCmd.Output()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println("> date")
	// fmt.Println(string(dateOut))
}
