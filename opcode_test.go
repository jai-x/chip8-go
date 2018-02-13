package main

import (
	"testing"
	"math/rand"
)

// Split uint16 opcodes into uint8 bytes storable in Chip8 memory
func byteChunk(code []uint16) []uint8 {
	var out []uint8
	for _, op := range(code) {
		out = append(out, uint8(op >> 8))
		out = append(out, uint8(op))
	}
	return out
}

// Test abilty to perform a direct jump to an address
func TestDirectJump(t *testing.T) {
	// Opcode to jump to address 0x456
	code := []uint16{0x1456}
	data := byteChunk(code)

	emu := NewChip8()
	emu.LoadProgram(data)

	// Complete one cycle such that program counter should now be set to the jump address
	emu.Cycle()

	// Test for jump at program counter
	if emu.PC != 0x456 {
		t.Errorf("Jump instruction failed (PC), target: 0x%04X, actual: 0x%04X", 0x456, emu.PC)
	}
}

// Test ability to perform a subroutine call and return
func TestSubroutineCallAndReturn(t *testing.T) {
	// Opcode to call routine at 0x204 which will immediatly return
	code := []uint16{0x2204, 0x0000, 0x00EE}
	data := byteChunk(code)

	emu := NewChip8()
	emu.LoadProgram(data)

	// Complete one cycle such that the subroutine has been called
	emu.Cycle()

	// Test for previous address stored in stack
	if emu.STACK[0] != 0x200 {
		t.Errorf("Stack address storage failed (STACK[0]), target: 0x%04X, actual: 0x%04X", 0x200, emu.STACK[0])
	}

	// Test if stack pointer has incremented
	if emu.SP != 1 {
		t.Errorf("Stack pointer increment failed (SP), target: 0x%02X, actual: 0x%02X", 1, emu.SP)
	}

	// Test if program counter has been set to new instruction
	if emu.PC != 0x204 {
		t.Errorf("Subroutine call failed (PC), target: 0x%04X, actual: 0x%04X", 0x204, emu.PC)
	}

	// Complete one cycle such that the will execute
	// The subroutine has only on instruction which is to return
	emu.Cycle()

	// Test if the stack pointer has decremented
	if emu.SP != 0 {
		t.Errorf("Stack pointer decrement failed (SP), target: 0x%02X, actual: 0x%02X", 0, emu.SP)
	}

	// Test if the program counter has been set to the previous memory address + 0x2
	if emu.PC != 0x202 {
		t.Errorf("Subroutine return failed (PC), target: 0x%04X, actual: 0x%04X", 0x202, emu.PC)
	}
}

func TestVRegisterAssign(t *testing.T) {
	// Opcode to assign register V0 with 0xAB and V1 with 0xCD
	code := []uint16{0x60AB, 0x61CD}
	data := byteChunk(code)

	emu := NewChip8()
	emu.LoadProgram(data)

	// Complete two cycles such that both values are stored
	emu.Cycle()
	emu.Cycle()

	// Test correct values have been stored
	if emu.V[0] != 0xAB {
		t.Errorf("Register assign failed (V0), target: 0x%02X, actual: 0x%02X", 0xAB, emu.V[0])
	}

	if emu.V[1] != 0xCD {
		t.Errorf("Register assign failed (V1), target: 0x%02X, actual: 0x%02X", 0xCD, emu.V[1])
	}
}

func TestDisplayClear(t *testing.T) {
	// Opcode to clear display
	code := []uint16{0x00E0}
	data := byteChunk(code)

	emu := NewChip8()
	emu.LoadProgram(data)

	// Pollute display area memory with random values
	for x, _ := range(emu.SCREEN) {
		for y, _ := range(emu.SCREEN[x]) {
			// Random true or false for monochrome display
			emu.SCREEN[x][y] = rand.Intn(1) != 0
		}
	}

	// Complete once cycle to execute the opcode to clear the display
	emu.Cycle()

	// Check that all display values have been cleared
	for x, _ := range(emu.SCREEN) {
		for y, _ := range(emu.SCREEN[x]) {
			val := emu.SCREEN[x][y]
			if val != false {
				t.Errorf("Display clear failed (SCREEN[%d][%d]), target: false, actual: %b", x, y, val)
			}
		}
	}
}
