package main

import (
	"fmt"
	"os"
	"github.com/nsf/termbox-go"
)

func main() {
	if len(os.Args) < 2 {
		die("Please supply path to ROM file as the first argument.")
	}

	err := termbox.Init()
	if err != nil {
		die("Termbox did not initialise correctly:", err)
	}

	emu := NewChip8()
	emu.LoadROM(os.Args[1])
	for {
		emu.Cycle()
	}

	termbox.Close()
}


// Println to stderr and exit
func die(msg ...interface{}) {
	if termbox.IsInit {
		termbox.Close()
	}

	fmt.Fprintln(os.Stderr, msg...)
	os.Exit(1)
}

// Printf to stderr and exit
func dief(format string, v ...interface{}) {
	if termbox.IsInit {
		termbox.Close()
	}

	fmt.Fprintf(os.Stderr, format, v...)
	os.Exit(1)
}
