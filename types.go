package main

import (
	"unsafe"
)

// Coord represents an X-Y coordinate
type Coord struct {
	X int16
	Y int16
}

// CoordToPointer casts the Coord to a uintptr pointer
func CoordToPointer(c Coord) uintptr {
	return uintptr(*((*uint32)(unsafe.Pointer(&c))))
}

// SmallRect represents an rectangle area
type SmallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

// ConsoleScreenBufferInfo holds information about an ConsoleScreenBuffer
type ConsoleScreenBufferInfo struct {
	Size              Coord
	CursorPosition    Coord
	Attributes        uint16
	Window            SmallRect
	MaximumWindowSize Coord
}
