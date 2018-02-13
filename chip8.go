package main

import (
	"io/ioutil"
)

const (
	SCREENW uint8 = 64
	SCREENH uint8 = 32
)

type Chip8 struct {
	// 16 general purpose 8 bit registers
	V [0x10]uint8

	// Special 16 bit register called I
	I uint16

	// 16 bit program counter
	PC uint16

	// 8 bit stack pointer
	SP uint8

	// 16 stack locations to store current 16bit memory address
	STACK [0x10]uint16

	// Main memory
	// 4096 bytes of contiguous memory
	MEM [0x1000]uint8

	// Current fetched opcode
	OPCODE uint16

	// 64x32 size monochrome display
	SCREEN [SCREENW][SCREENH]bool

	// Delay timer
	// When above zero, will decrement at a rate of 60Hz
	DT uint8

	// Sound timer
	// When above zero, will play a tone and also decrement at a rate of 60Hz
	ST uint8

	// Keyboard values
	// Represents the state of the each hexadecimal key value press
	// 1 for pressed, 0 for not pressed
	KEY [0xF]uint8

	// Channel to await a keypress blocking
	// Receives a hexadecimal value (0x0 - 0xF)
	KEYCHAN chan uint8
}

func NewChip8() *Chip8 {
	var out Chip8

	// Load the sprite table into memory
	// Sprites are monochrome bitmaps representing the hexadecimal characters
	// Sprite of hex value 0xn will be at offset in memory 5 * 0xn

	// http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#2.4
	sprites := []uint8 {
		0x0F, 0x90, 0x90, 0x90, 0xF0, // Sprite of 0
		0x20, 0x60, 0x20, 0x20, 0x70, // Sprite of 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // Sprite of 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // Sprite of 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // Sprite of 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // Sprite of 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // Sprite of 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // Sprite of 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // Sprite of 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // Sprite of 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // Sprite of A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // Sprite of B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // Sprite of C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // Sprite of D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // Sprite of E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // Sprite of F
	}

	// Sprites are stored in interpreter memory from 0x000 to 0x1FF
	out.LoadData(sprites, 0)

	// Program counter starts at 0x200 where the ROM is stored in memory
	out.PC = 0x200

	// Initialise keypress channel
	out.KEYCHAN = make(chan uint8)

	return &out
}

// Load ROM file from string filepath
func (c *Chip8) LoadROM(rompath string) {
	data, err := ioutil.ReadFile(rompath)
	if err != nil {
		die("Could not read ROM file at path:", rompath)
	}
	c.LoadProgram(data)
}

// Load program data directly into program space in the Chip8 memory bank
func (c *Chip8) LoadProgram(data []uint8) {
	// MEM from 0x000 to 0x1FF is reserved for the interpreter
	// Loaded roms start from 0x200
	c.LoadData(data, 0x200)
}

// Load array of data directly into Chip8 memory at given offset
func (c *Chip8) LoadData(data []uint8, offset int) {
	copy(c.MEM[offset:], data)
}

// Complete one CPU cycle of the Chip8
func (c *Chip8) Cycle() {
	// Fetch the uint16 opcode from two contiguous uint8 memory locations
	// offset by the program counter variable
	// Bitshift first byte to left by 8 and or with second byte
	c.OPCODE = uint16(c.MEM[c.PC])<<8 | uint16(c.MEM[c.PC+1])


	// The first 4 bits of the opcode contain the instruction type
	// Bitwise AND(&) with 0xF000 to remove the trailing 12 bits
	instruction := c.OPCODE & 0xF000

	// Dispatch to opcode implemmentation
	switch instruction {
	case 0x0000:
		c.Opcode0()
	case 0x1000:
		c.Opcode1()
	case 0x2000:
		c.Opcode2()
	case 0x3000:
		c.Opcode3()
	case 0x4000:
		c.Opcode4()
	case 0x5000:
		c.Opcode5()
	case 0x6000:
		c.Opcode6()
	case 0x7000:
		c.Opcode7()
	case 0x8000:
		c.Opcode8()
	case 0x9000:
		c.Opcode9()
	case 0xA000:
		c.OpcodeA()
	case 0xB000:
		c.OpcodeB()
	case 0xC000:
		c.OpcodeC()
	case 0xD000:
		c.OpcodeD()
	case 0xE000:
		c.OpcodeE()
	case 0xF000:
		c.OpcodeF()
	default:
		dief("Unknown instruction 0x%04X in opcode 0x%04X\n", instruction, c.OPCODE)
	}
}
