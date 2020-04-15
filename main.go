package main

import (
	"fmt"
	"log"
	"sync"
	"syscall"
	"unsafe"
)

const command = "./example/WriteConsoleA.exe"

func main() {
	hConsole, err := CreateConsoleScreenBuffer()
	if err != nil {
		log.Fatalln("could not create console screen buffer", err)
	}
	defer syscall.CloseHandle(hConsole)

	// we have to start the process manually since
	// we dont have access to syscall.ProcessInformation
	// using exec.Command

	argvp, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		log.Fatalln("could not convert the command input to a utf16ptr", err)
	}

	si := new(syscall.StartupInfo)
	si.Cb = uint32(unsafe.Sizeof(*si))
	si.Flags = syscall.STARTF_USESTDHANDLES
	si.Flags |= StartFForceOffFeedBack
	si.ShowWindow = syscall.SW_SHOW

	si.StdOutput = hConsole // <- This is the part which we cannot change using the current golang api
	si.StdErr = hConsole    // <-

	pi := new(syscall.ProcessInformation)

	err = syscall.CreateProcess(nil, argvp, nil, nil, true, 0, nil, nil, si, pi)
	if err != nil {
		log.Fatalln("could not create the process", err)
	}

	// we dont have access to pi.Process as well, but we could use win32api OpenProcess to get an handle to it
	outputChan, errorsChan := ReadConsoleOutput(pi.Process, hConsole)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for buf := range outputChan {
			fmt.Print(buf)
		}

		fmt.Println()
		wg.Done()
	}()

	go func() {
		for err := range errorsChan {
			fmt.Println(err)
		}

		fmt.Println()
		wg.Done()
	}()

	wg.Wait()
}
