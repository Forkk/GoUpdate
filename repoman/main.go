// Copyright 2013 MultiMC Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"strconv"
)

type CommandError interface {
	Error() string
	GetExitCode() int
	GetCause() error
	ShouldPrintUsage() bool
}

type NormalCommandError struct {
	// The error message to display.
	msg string

	// The exit code value that this error should cause the process to exit with.
	exitCode int

	// The error that caused this error. If not nil, this this CommandError will an "Caused by <cause>"
	cause error

	// If true, this error should print the command's usage message.
	printUsage bool
}

// Implement the error interface.
func (err NormalCommandError) Error() string {
	if err.msg == "" {
		return ""
	} else {
		if err.cause != nil {
			return fmt.Sprintf("%s\n    Caused by: %s", err.msg, err.cause)
		} else {
			return err.msg
		}
	}
}

func (err NormalCommandError) GetExitCode() int {
	return err.exitCode
}

func (err NormalCommandError) GetCause() error {
	return err.cause
}

func (err NormalCommandError) ShouldPrintUsage() bool {
	return err.printUsage
}

// Returns a NormalCommandError with the given message, exit code, and cause.
func CausedError(message string, exitCode int, cause error) CommandError {
	return NormalCommandError{msg: message, exitCode: exitCode, cause: cause, printUsage: false}
}

// Returns a NormalCommandError with the given message and exit code, but no cause.
func ErrorMessage(message string, exitCode int) CommandError {
	return NormalCommandError{msg: message, exitCode: exitCode, cause: nil, printUsage: false}
}

// Returns a NormalCommandError with no message that will print usage information.
func CmdUsageError() CommandError {
	return NormalCommandError{msg: "", exitCode: 1, cause: nil, printUsage: true}
}

// Returns a NormalCommandError with a message that will also print usage information.
func CmdUsageErrorMsg(message string) CommandError {
	return NormalCommandError{msg: message, exitCode: 1, cause: nil, printUsage: true}
}

type commandFunc func(args ...string) CommandError

type CommandInfo struct {
	// The function to call to execute this command.
	CmdFunc commandFunc

	// String to display in the list of commands. Should not contain any newlines.
	HelpSummary string

	// Longer usage message to display with the command's help that tells the user how to use the command.
	UsageMessage string
}

var commands map[string]CommandInfo

func main() {
	// Initialize the command map.
	commands = map[string]CommandInfo{
		"help":   CommandInfo{CmdFunc: helpCommand, HelpSummary: "Shows a list of available commands and some basic information about them.", UsageMessage: "help"},
		"create": CommandInfo{CmdFunc: createCommand, HelpSummary: "Creates a new, blank repository in the specified directory.", UsageMessage: "create REPO_DIR"},
		"update": CommandInfo{CmdFunc: updateCommand, HelpSummary: "Updates an existing repository with a set of files.", UsageMessage: "update REPO_DIR FILE_STORAGE URL_BASE NEW_VERSION_DIR VERSION_NAME VERSION_ID\n    Updates the repository in REPO_DIR with the files in NEW_VERSION_DIR with the version name VERSION_NAME and the ID VERSION_ID, storing files in FILE_STORAGE with the base URL URL_BASE."},
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
		if err.Error() != "" {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		if err.ShouldPrintUsage() {
			fmt.Fprintf(os.Stderr, "Usage: %s\n", cmd.UsageMessage)
		}
		return err.GetExitCode()
	}
}

///////////////////////////////////
//////// COMMAND FUNCTIONS ////////
///////////////////////////////////

// Returns a string containing the list of commands that should be printed with the help message when no command is specified.
func helpCommand(args ...string) CommandError {
	help := fmt.Sprintf("Usage: %s COMMAND [arg...]\n", os.Args[0])

	for cmdStr, cmdInfo := range commands {
		help += fmt.Sprintf("    %-10.10s%s\n", cmdStr, cmdInfo.HelpSummary)
	}

	fmt.Fprintf(os.Stderr, help)

	return nil
}

func createCommand(args ...string) CommandError {
	// Determine what directory to create the repository in.
	if len(args) <= 0 {
		return CmdUsageErrorMsg("'create' command requires at least one argument.")
	} else {
		repoDir := args[0]
		return CreateRepo(repoDir)
	}
}

func updateCommand(args ...string) CommandError {
	if len(args) < 6 {
		return CmdUsageErrorMsg("'update' command requires at least six arguments.")
	} else {
		repoDir := args[0]
		filesDir := args[1]
		urlBase := args[2]
		newVersionDir := args[3]
		versionName := args[4]
		versionIdStr := args[5]

		versionId, err := strconv.ParseInt(versionIdStr, 10, 0)

		if err != nil {
			return CmdUsageErrorMsg("Version ID must be a positive integer.")
		} else {
			return UpdateRepo(repoDir, filesDir, urlBase, newVersionDir, versionName, int(versionId))
		}
	}
}
