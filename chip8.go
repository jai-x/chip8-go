package main

import (
	"io/ioutil"
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
	SCREEN [64][32]uint8
}

func NewChip8() *Chip8 {
	var out Chip8
	// Program counter starts at 0x200 where the ROM is stored in memory
	out.PC = 0x200
	return &out
}

// Load ROM file from string filepath
func (c *Chip8) LoadROM(rompath string) {
	data, err := ioutil.ReadFile(rompath)
	if err != nil {
		die("Could not read ROM file at path:", rompath)
	}

	c.LoadData(data)
}

// Load array of program data directly into ROM space in the Chip8 memory bank
func (c *Chip8) LoadData(data []uint8) {
	// MEM from 0x000 to 0x1FF is reserved for the interpreter
	// Loaded roms start from 0x200
	copy(c.MEM[0x200:], data)
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
