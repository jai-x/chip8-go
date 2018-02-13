package main

import (
	"math/rand"
)

// Each function implements the 2 byte opcode beginning with the named
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
	inst := uint8(c.OPCODE)

	switch inst {
	// Display clear
	case 0xE0:
		// Zero all values in screen array
		// Quick reinitialise of 2d array
		// 64 x 32 monochrome display
		c.SCREEN = [SCREENW][SCREENH]bool{}

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

// 0x2nnn - Call subroutine at address nnn
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
	// x consists of bits 4 - 7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// y consists of bits 8 - 11
	y := uint8(c.OPCODE >> 4) & 0x0F
	// Last 4 bits of the opcode determines the operation
	inst := uint8(c.OPCODE) >> 4

	switch inst {
	// Vx = Vy
	case 0x0:
		c.V[x] = c.V[y]

	// Vx = Vx | Vy
	case 0x1:
		c.V[x] = c.V[x] | c.V[y]

	// Vx = Vx & Vy
	case 0x2:
		c.V[x] = c.V[x] & c.V[y]

	// Vx = Vx ^ Vy
	case 0x3:
		c.V[x] = c.V[x] ^ c.V[y]

	// Vx += Vy
	case 0x4:
		c.V[x] += c.V[y]

	// Vx -= Vy
	case 0x5:
		c.V[x] -= c.V[y]

	// Vx = Vy = Vy >> 1 (Vf = Vy LSB)
	case 0x6:
		// Put LSB into Vf
		c.V[0xF] = c.V[y] & 0x1
		// Bitshift and assign to both registers
		res := c.V[y] >> 1
		c.V[x], c.V[y] = res, res

	// Vx = Vy - Vx (Vf = carry ? 1 : 0)
	case 0x7:
		carry := 0
		if c.V[y] > c.V[x] {
			carry = 1
		}
		c.V[0xF] = uint8(carry)
		c.V[x] = c.V[y] - c.V[x]

	// Vx = Vy = Vy << 1 (Vf = Vy MSB)
	case 0xE:
		// Put MSB into Vf
		c.V[0xF] = c.V[y] & 0x8
		// Bitshift and assign to both registers
		res := c.V[y] << 1
		c.V[x], c.V[y] = res, res
	}

	// Move to next 2 byte opcode
	c.PC += 0x2
}

// 0x9xy0 - Conditional check of registers Vx and Vy
// If Vx != Vy, skip the next instruction
func (c *Chip8) Opcode9() {
	// x consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// y consists of the bits 8-11
	y := uint8(c.OPCODE >> 4) & 0x0F

	// A full instruction  is 2 bytes, so a skip is 4 bytes
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
	// nn consists of the bits 8-15
	nn := uint8(c.OPCODE)

	c.V[x] = uint8(rand.Intn(256)) & nn

	// Move to next 2 byte opcode
	c.PC += 0x2
}

// dank bits boi
func bytetobitarray(theByte uint8) (out [8]bool) {
	for i := uint8(0); i < 8; i++ {
		out[i] = ((theByte & (1 << i)) >> i) != 0
	}
	return
}

// 0xDxyn - Display sprite at screen location x,y of size n, data from offset I
// Sprite is of fixed width 8 pixels with data read starting from address
// pointed to by I. Sprites will wraparound screen top and bottom.
// Sprite values are XOR'd(^) with current screen values. If any screen values
// are set to 0, Vf is set to 1.
func (c *Chip8) OpcodeD() {
	// x value consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F
	// y value consists of the bits 8-11
	y := uint8(c.OPCODE >> 4) & 0x0F
	// n value consists of the bits 12 - 16
	n := uint8(c.OPCODE) & 0x0F

	// Read n bytes of sprite data starting from address at I
	sprite := make([][8]bool, n)
	for offset := 0; uint8(offset) < n; offset++ {
		line := c.MEM[c.I + uint16(offset)]
		sprite[offset] = bytetobitarray(line)
	}

	for w, line := range(sprite) {
		for h, bit := range(line) {
			// Determine postion of bit with screen wrap-around
			xpos := (uint8(w) + x) % SCREENW
			ypos := (uint8(h) + y) % SCREENH
			// XOR bit onto screen
			c.SCREEN[xpos][ypos] = c.SCREEN[xpos][ypos] != bit
		}
	}

	// Move to next 2 byte opcode
	c.PC += 0x2
}

// 0xEx9E - Skip next instruction if the key value stored in Vx is currently pressed
// 0xExA1 - Skip next instruction if the key value stored in Vx is not currenly pressed
func (c *Chip8) OpcodeE() {
	// Last 8 bits determine instruction of opcode
	inst := uint8(c.OPCODE)
	// x value consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F

	var skip bool

	switch inst {
	// Skip next instruction if KEY[Vx] == 1
	case 0x9E:
		if c.KEY[c.V[x]] == 1 {
			skip = true
		}

	// Skip next instruction if KEY[Vx] == 0
	case 0xA1:
		if c.KEY[c.V[x]] == 0 {
			skip = true
		}
	}

	// A full instruction is 2 bytes, so a skip is 4 bytes
	if skip {
		c.PC += 0x4
	} else {
		c.PC += 0x2
	}
}

// 0xFx07 - Set Vx to value of the delay timer
// 0xFx0A - Block await for a keypress. Store key value into Vx
// 0xFx15 - Set the delay timer to value of Vx
// 0xFx18 - Set the sound timer to value of Vx
// 0xFx1E - I += Vx
// 0xFx29 - Set I to the location of the sprite values for the value in Vx
// 0xFx33 - Store binary coded decimal value of Vx and store in memory from I
// 0xFx55 - Dump registers V0 to Vx to memory at address starting from I
// 0xFx65 - Load registers V0 to Vx from memory at address starting from I
func (c *Chip8) OpcodeF() {
	// Last 8 bits determine instruction of opcode
	inst := uint8(c.OPCODE)
	// x value consists of the bits 4-7
	x := uint8(c.OPCODE >> 8) & 0x0F

	switch inst {
	// Vx = DT
	case 0x07:
		c.V[x] = c.DT

	// Block await for a keypress. Store key value into Vx
	case 0x0A:
		c.V[x] = <-c.KEYCHAN

	// DT = Vx
	case 0x15:
		c.DT = c.V[x]

	// ST = Vx
	case 0x18:
		c.ST = c.V[x]

	// I += Vx
	case 0x1E:
		c.I += uint16(c.V[x])

	// Set I to the location of the sprite values for the value in Vx
	case 0x29:
		// Sprites are stored in interpreter memory from 0x000 to 0x1FF
		// Sprite of hex value 0xn will be at offset in memory 5 * 0xn
		c.I = 5 * uint16(c.V[x])

	// Store binary coded decimal value of Vx and store in memory from I
	case 0x33:
		// Hundreds
		c.MEM[c.I]     = uint8(int(c.V[x]) % 1000 / 100)
		// Tens
		c.MEM[c.I + 1] = uint8(int(c.V[x]) % 100 / 10)
		// Ones
		c.MEM[c.I + 2] = uint8(int(c.V[x]) % 10)

	// Dump registers V0 to Vx to memory at address starting from I
	case 0x55:
		copy(c.MEM[c.I:], c.V[0:x])

	// Load registers V0 to Vx from memory at address starting from I
	case 0x65:
		copy(c.V[0:x], c.MEM[c.I:])
	}

	// Move to next 2 byte opcode
	c.PC += 0x2
}
