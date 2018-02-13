package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		die("Please supply path to ROM file as the first argument.")
	}

	emu := NewChip8()
	emu.LoadROM(os.Args[1])
	for {
		emu.Cycle()
	}
}


// Println to stderr and exit
func die(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

// Printf to stderr and exit
func dief(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	os.Exit(1)
}
