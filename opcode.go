package main

import (
	"math/rand"
)

// Each function implemments the 2 byte opcode beginning with the named
// hexadecimal value in the function name

// 0x0nnn - Calls RCA 1802 Program at address nnn
// Not needed can ignore
// ----
// 0x00E0 - Display Clear
// Zero memory area tied to to display
// ----
// 0x00EE - Return from subroutine
// Set the program counter at the previous address found in the stack
// Increment program counter to next instruction
func (c *Chip8) Opcode0() {

	// Only last byte needed to determine what to do
	last := uint8(c.OPCODE)
	switch last {
	// Display clear
	case 0xE0:
		// Zero all values in screen array
		for x, _ := range(c.SCREEN) {
			for y, _ := range(c.SCREEN[x]) {
				c.SCREEN[x][y] = 0
			}
		}

	// Subroutine return
	case 0xEE:
		// Decrement stack pointer to stored address
		c.SP--
		// Set program to stored address
		c.PC = c.STACK[c.SP]
	}

	// Move to next 2 byte opcode
	c.PC += 0x2
}

// 0x1nnn - Direct jump to address at nnn
func (c *Chip8) Opcode1() {
	// Discard first 4 bits to obtain address
	nnn := c.OPCODE & 0x0FFF
	// Set program counter to address in a direct jump
	c.PC = nnn
}

// Ox2nnn - Call subroutine at address nnn
// Store current program counter in stack location indicated by stack pointer
// Increment stack pointer and set program counter to subroutine address
func (c *Chip8) Opcode2() {
	// Discard first 4 bits to obtain address
	nnn := c.OPCODE & 0x0FFF
	// Store current program counter in stack
	c.STACK[c.SP] = c.PC
	// Increment stack pointer
	c.SP++
	// Set program counter to new address
	c.PC = nnn
}

// 0x3xnn - Conditional check of register Vx and given value nn
// If Vx == nn, skip the next instruction
func (c *Chip8) Opcode3() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// nn consists of the bits 8-16
	nn := uint8(c.OPCODE)

	// A full instruction is 2 bytes, so a skip is 4 bytes
	if c.V[x] == nn {
		c.PC += 0x4
	} else {
		c.PC += 0x2
	}
}

// 0x4xnn - Conditional check of register Vx and given value nn
// If Vx != nn, skip the next instruction
func (c *Chip8) Opcode4() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// nn consists of the bits 8-16
	nn := uint8(c.OPCODE)

	// A full instruction is 2 bytes, so a skip is 4 bytes
	if c.V[x] != nn {
		c.PC += 0x4
	} else {
		c.PC += 0x2
	}
}

// 0x5xy0 - Conditonal check of registers Vx and Vy
// If Vx == Vy, skip the next instruction
func (c *Chip8) Opcode5() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// y consists of the bits 8-11
	y := uint8(c.OPCODE >> 4) & 0x0F

	// A full instruction  is 2 bytes, so a skip by 4 bytes
	if c.V[x] == c.V[y] {
		c.PC += 0x4
	} else {
		c.PC += 0x2
	}
}

// 0x6xnn - Assign Vx = nn
func (c *Chip8) Opcode6() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// nn consists of the bits 8-16
	nn := uint8(c.OPCODE)

	c.V[x] = nn

	// Move to next 2 byte opcode
	c.PC += 0x2
}

// 0x7xnn - Add Vx += nn
func (c *Chip8) Opcode7() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// nn consists of the bits 8-16
	nn := uint8(c.OPCODE)

	c.V[x] += nn

	// Move to next 2 byte opcode
	c.PC += 0x2
}

// 0x8xyn - Maths operations
// 0x8xy0 - Assign   Vx = Vy
// 0x8xy1 - Bitwise  Vx = Vx | Vy
// 0x8xy2 - Bitwise  Vx = Vx & Vy
// 0x8xy3 - Bitwise  Vx = Vx ^ Vy
// 0x8xy4 - Add      Vx += Vy
// 0x8xy5 - Subtract Vx -= Vy
// 0x8xy6 - Bitwise  Vx = Vy = Vy >> 1 (Vf = Vy LSB)
// 0x8xy7 - Subtract Vx = Vy - Vx (Vf = carry ? 1 : 0)
// 0x8xyE - Bitwise  Vx = Vy = Vy << 1 (Vf = Vy MSB)
func (c *Chip8) Opcode8() {
	dief("Opcode not implemmented yet: 0x%04X\n", c.OPCODE)
}

// 0x9xy0 - Conditional check of registers Vx and Vy
// If Vx != Vy, skip the next instruction
func (c *Chip8) Opcode9() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// y consists of the bits 8-11
	y := uint8(c.OPCODE >> 4) & 0x0F

	// A full instruction  is 2 bytes, so a skip by 4 bytes
	if c.V[x] != c.V[y] {
		c.PC += 0x4
	} else {
		c.PC += 0x2
	}
}

// 0xAnnn - Set I to nnn
func (c *Chip8) OpcodeA() {
	// Discard first 4 bits to obtain address
	nnn := c.OPCODE & 0x0FFF

	c.I = nnn

	// Move to next 2 byte opcode
	c.PC += 0x2
}

// 0xBnnn - Jump with offset
// Jump to address nnn + V0
func (c *Chip8) OpcodeB() {
	// Discard first 4 bits to obtain address
	nnn := c.OPCODE & 0x0FFF

	c.PC = nnn + uint16(c.V[0])
}

// 0xCxnn - Random byte with bitwise AND(&)
// Generate random byte (0, 255) and biwise AND(&) with nn
// Store the result in Vx
func (c *Chip8) OpcodeC() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// nn consists of the bits 8-16
	nn := uint8(c.OPCODE)

	c.V[x] = uint8(rand.Intn(256)) & nn

	// Move to next 2 byte opcode
	c.PC += 0x2
}

//:0xDxyn - Display sprite at screen location x,y of size n, date from offset I
// Sprite is of fixed width 8 pixels with data read starting from address
// pointed to by I. Sprites will wraparound screen top and bottom.
// Sprite values are XOR'd(^) with current screen values. If any screen values
// are set to 0, Vf is set to 1.
func (c *Chip8) OpcodeD() {
	dief("Opcode not implemmented yet: 0x%04X\n", c.OPCODE)
}

// 0xEx9E - Skip next instruction if the key value stored in Vx is currently pressed
// 0xExA1 - Skip next instruction if the key value stored in Vx is not currenly pressed
func (c *Chip8) OpcodeE() {
	dief("Opcode not implemmented yet: 0x%04X\n", c.OPCODE)
}

// 0xFx07 - Set Vx to value of the delay timer
// 0xFx0A - Block await for a keypress. Store key value into Vx
// 0xFx15 - Set the delay timer to value of Vx
// 0xFx18 - Set the sound timer to value of Vx
// 0xFx1E - I += Vx
// 0xFx29 - Set I to the location of the sprite values for the value in Vx
// 0xFx33 - BCD something TODO
// 0xFx55 - Dump registers V0 to Vx to memory at address starting from I
// 0xFx65 - Load registers V0 to Vx from memory at address starting from I
func (c *Chip8) OpcodeF() {
	dief("Opcode not implemmented yet: 0x%04X\n", c.OPCODE)
}
