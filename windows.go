// +build windows

package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const (
	StartFForceOffFeedBack = 0x00000080

	GenericRead  = 0x80000000
	GenericWrite = 0x40000000

	FileShareRead  = 0x00000001
	FileShareWrite = 0x00000002

	ConsoleTextModeBuffer = 0x1
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	createConsoleScreenBufferProc   = kernel32.NewProc("CreateConsoleScreenBuffer")
	getConsoleScreenBufferInfoProc  = kernel32.NewProc("GetConsoleScreenBufferInfo")
	readConsoleOutputCharacterAProc = kernel32.NewProc("ReadConsoleOutputCharacterA")
	setConsoleCursorPositionProc    = kernel32.NewProc("SetConsoleCursorPosition")
)

func CreateConsoleScreenBuffer() (syscall.Handle, error) {
	sa := new(syscall.SecurityAttributes)
	sa.InheritHandle = 1
	sa.SecurityDescriptor = 0
	sa.Length = uint32(unsafe.Sizeof(*sa))

	r1, r2, err := createConsoleScreenBufferProc.Call(uintptr(GenericRead|GenericWrite),
		uintptr(FileShareRead|FileShareWrite), uintptr(unsafe.Pointer(sa)), uintptr(ConsoleTextModeBuffer), 0)

	return syscall.Handle(r1), checkError(r1, r2, err)
}

func GetConsoleScreenBufferInfo(handle syscall.Handle) (*ConsoleScreenBufferInfo, error) {
	info := ConsoleScreenBufferInfo{}
	err := checkError(getConsoleScreenBufferInfoProc.Call(uintptr(handle), uintptr(unsafe.Pointer(&info)), 0))
	return &info, err
}

func SetConsoleCursorPosition(handle syscall.Handle, coord Coord) error {
	r1, r2, err := setConsoleCursorPositionProc.Call(uintptr(handle), CoordToPointer(coord))
	return checkError(r1, r2, err)
}

func ReadConsoleOutput(hProcess syscall.Handle, hConsole syscall.Handle) (<-chan string, <-chan error) {
	var lastPosition Coord
	var origin Coord
	var csbi *ConsoleScreenBufferInfo

	output := make(chan string)
	errors := make(chan error)

	go func() {
		for {
			r, err := syscall.WaitForSingleObject(hProcess, 0)
			if err != nil {
				errors <- err
				close(output)
				close(errors)
				break
			}

			if r == syscall.WAIT_OBJECT_0 {
				log.Println("process finalized")
				close(output)
				close(errors)
				break
			}

			if r != syscall.WAIT_TIMEOUT {
				errors <- fmt.Errorf("syscall.WaitForSingleObject != syscall.WAIT_TIMEOUT: response %v", r)
				close(output)
				close(errors)
				break
			}

			csbi, err = GetConsoleScreenBufferInfo(hConsole)
			if err != nil {
				errors <- fmt.Errorf("GetConsoleScreenBufferInfo error: response %v", err)
				close(output)
				close(errors)
				break
			}

			lineWidth := csbi.Size.X

			if csbi.CursorPosition.X == lastPosition.X && csbi.CursorPosition.Y == lastPosition.Y {
				continue
			} else {
				count := (csbi.CursorPosition.Y-lastPosition.Y)*lineWidth + csbi.CursorPosition.X - lastPosition.X

				buf := make([]byte, count)
				r1, r2, err := readConsoleOutputCharacterAProc.Call(uintptr(hConsole), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), CoordToPointer(lastPosition), uintptr(unsafe.Pointer(&count)))
				err = checkError(r1, r2, err)
				if err != nil {
					errors <- fmt.Errorf("readConsoleOutputCharacterAProc error: response %v", err)
					close(output)
					close(errors)
					break
				}

				lastPosition = csbi.CursorPosition
				csbi, err = GetConsoleScreenBufferInfo(hConsole)
				if err != nil {
					errors <- fmt.Errorf("GetConsoleScreenBufferInfo error: response %v", err)
					close(output)
					close(errors)
					break
				}

				if csbi.CursorPosition.X == lastPosition.X && csbi.CursorPosition.Y == lastPosition.Y {
					_ = SetConsoleCursorPosition(hConsole, origin)
					lastPosition = origin
				}

				output <- string(buf[0:count])
			}
		}
	}()

	return output, errors
}

func checkError(r1, r2 uintptr, err error) error {
	if r1 != 0 {
		return nil
	}

	if err != nil {
		return err
	}

	return syscall.EINVAL
}
