## Example capturing output from an executable that uses ConsoleScreenBuffer

**NOTE: This should not be used on production.**

On Windows if you use `WriteConsole` family of functions, you cannot redirect the Stdout, Stderr or capture it
using `exec.Command` (if you use a custom writer or custom file it will be empty, if you use the `os.Stdout` it will be piped correctly, but you cannot capture it, the `StdoutPipe` will not work as well).
This behaviour occurs with others programming languages as well. It will not work while trying to pipe on terminal as well.


You can try using the example C++ application on the example folder.


To capture the output you have to create an `ConsoleScreenBuffer`, attach it as the stdout and/or stderr of the process and start the process.
You can then read the ConsoleScreenBuffer and convert it to string. 

**In this repository there is an example of running a subprocess that writes to a ConsoleScreenBuffer and reading the output on GOLANG**.

We create a process using low-level system calls and use a lot of `win32api` calls to make this works, its not perfect or even
a good implementation (I think is missing a lot of error checking and process state handling). 

If the GO standard library allowed us to inform the `syscall.ProcessInformation` we could use a different approach, a simpler, short and direct way to resolve this issue.
 
