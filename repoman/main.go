package main

import (
    "fmt"
    "os"
)

type commandError interface {
    Error() string
    GetExitCode() int
    GetCause() error
    ShouldPrintUsage() bool
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
func (err commandErrorImpl) Error() string {
    if err.cause != nil {
        return fmt.Sprintf("%s\n    Caused by: %s", err.msg, err.cause)
    } else {
        return err.msg
    }
}

func (err commandErrorImpl) GetExitCode() int {
    return err.exitCode
}

func (err commandErrorImpl) GetCause() error {
    return err.cause
}

func (err commandErrorImpl) ShouldPrintUsage() bool {
    return false
}



type commandFunc func(args ...string) commandError

type CommandInfo struct {
    // The function to call to execute this command.
    CmdFunc commandFunc
    
    // String to display in the list of commands. Should not contain any newlines.
    HelpSummary string

    // Longer usage message to display with the command's help that tells the user how to use the command.
    UsageMessage string
}


// Special command error that prints the command's usage info.
type commandUsageError struct {
    cmd CommandInfo
}


func (err commandUsageError) Error() string { return "Invalid arguments to command." }

func (err commandUsageError) GetExitCode() int { return 1 }

func (err commandUsageError) GetCause() error { return nil }

func (err commandUsageError) ShouldPrintUsage() bool { return true }

var commands map[string]CommandInfo

func main() {
    // Initialize the command map.
    commands = map[string]CommandInfo {
        "help": CommandInfo{CmdFunc: helpCommand, HelpSummary: "Shows a list of available commands and some basic information about them.", UsageMessage: "help"},
        "create": CommandInfo{CmdFunc: createCommand, HelpSummary: "Creates a new, blank repository in the specified directory.", UsageMessage: "create REPO_DIR"},
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
        help += fmt.Sprintf("    %-10.10s%s\n", cmdStr, cmdInfo.HelpSummary)
    }

    fmt.Fprintf(os.Stderr, help)

    return nil
}

func createCommand(args ...string) commandError {
    usage := commands["create"].UsageMessage

    // Determine what directory to create the repository in.
    if len(args) <= 0 {
        return commandErrorImpl{fmt.Sprintf("Usage: %s\n", usage), 1, nil}
    } else {
        repoDir := args[0]

        err := CreateRepo(repoDir)
        
        if err == nil {
            return nil
        } else {
            return commandErrorImpl{"An error occurred.", 1, err}
        }
    }
}

