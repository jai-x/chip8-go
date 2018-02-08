package main

import (
	"fmt"
)

// dump all stack
func (c *Chip8) DumpStack() {
	fmt.Println("Stack:")
	for i, _ := range(c.STACK) {
		c.PrintStackAt(i)
	}
}

// print one stack value
func (c *Chip8) PrintStackAt(i int) {
	format := "0x%X: [ 0x%04X ]"
	if uint8(i) == c.SP {
		format += " <--SP"
	}
	format += "\n"

	fmt.Printf(format, i, c.STACK[i])
}

// dump all registers
func (c *Chip8) DumpRegisters() {
	fmt.Println("Registers:")
	for i, _ := range(c.V) {
		c.PrintRegisterAt(i)
	}
}

// print one register value
func (c *Chip8) PrintRegisterAt(i int) {
	fmt.Printf("0x%X: [ 0x%02X ]\n", i, c.V[i])
}

// print currently stored opcode
func (c *Chip8) DumpOp() {
	fmt.Printf("Opcode: 0x%04X\n", c.OPCODE)
}

