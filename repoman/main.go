package main

import (
    "fmt"
    "os"
)

type commandError interface {
    Error() string
    GetExitCode() int
    GetCause() int
}

type commandErrorImpl struct {
    // The error message to display.
    msg string

    // The exit code value that this error should cause the process to exit with.
    exitCode int
    
    // The error that caused this error. If not nil, this this commandError will an "Caused by <cause>" 
    cause error
}

// Implement the error interface.
func (err *commandErrorImpl) Error() string {
    if err.cause != nil {
        return fmt.Sprintf("%s\n    Caused by: %s", err.msg, err.cause)
    } else {
        return err.msg
    }
}

func (err *commandErrorImpl) GetExitCode() int {
    return err.exitCode
}

func (err *commandErrorImpl) GetCause() error {
    return err.cause
}


type commandFunc func(args ...string) commandError

type CommandInfo struct {
    // The function to call to execute this command.
    CmdFunc commandFunc
    
    // String to display in the list of commands. Should not contain any newlines.
    UsageSummary string
}


var commands map[string]CommandInfo

func main() {
    // Initialize the command map.
    commands = map[string]CommandInfo {
        "help": CommandInfo{CmdFunc: helpCommand, UsageSummary: "Shows a list of available commands and some basic information about them."},
        "test": CommandInfo{CmdFunc: testCommand, UsageSummary: "This is a test command."},
    }

    // Get the command line arguments.
    args := os.Args

    // There must be at least one argument (the sub-command). If not, print the help text and exit.
    if len(args) <= 1 {
        executeCommand(commands["help"])
        os.Exit(1)
    }

    // If there is a command argument, get it.
    cmd := args[1]

    // Look up the command in the command map.
    if cmdInfo, ok := commands[cmd]; ok {
        // Run the command.
        os.Exit(executeCommand(cmdInfo, args[2:]...))
    } else {
        // If the command doesn't exist, print the help message and exit.
        os.Exit(executeCommand(commands["help"], args[2:]...))
    }

    return
}

// Executes the given command and returns the exit code that the process should exit with.
func executeCommand(cmd CommandInfo, args ...string) int {
    err := cmd.CmdFunc(args...)
    if err == nil {
        return 0
    } else {
        fmt.Fprintf(os.Stderr, err.Error())
        return err.GetExitCode()
    }
}


///////////////////////////////////
//////// COMMAND FUNCTIONS ////////
///////////////////////////////////

// Returns a string containing the list of commands that should be printed with the help message when no command is specified.
func helpCommand(args ...string) commandError {
    help := fmt.Sprintf("Usage: %s COMMAND [arg...]\n", os.Args[0])
    
    for cmdStr, cmdInfo := range commands {
        help += fmt.Sprintf("    %-10.10s%s\n", cmdStr, cmdInfo.UsageSummary)
    }

    fmt.Fprintf(os.Stderr, help)

    return nil
}

func testCommand(args ...string) commandError {
    fmt.Printf("It worked!")
    return nil
}

