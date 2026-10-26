package main

import (
	"fmt"
	"os"
)

func handle_create(args []string) error {
	command := args[0]
	path := args[1]
	if command == "--file" {
	} else {
		return fmt.Errorf("not supported sub command")
	}
	return nil
}

func handle_pin(args []string) error {
	command := args[0]
	path := args[1]
	if command == "--file" {
	} else {
		return fmt.Errorf("not supported sub command")
	}
	return nil
}

func handle_args(argsv []string) error {
	argsc := len(argsv)
	if argsc != 3 {
		return fmt.Errorf("command line arguments are incorrect!")
	}

	switch command := argsv[0]; command {
	case "create":
		handle_create(argsv[1:])
		break
	case "pin":
		handle_pin(argsv[1:])
		break
	default:
		return fmt.Errorf("command line argument is not suported")
	}
	return nil
}

func main() {
	args := os.Args[1:]

	if err := handle_args(args); err != nil {
		os.Exit(-1)
	}
}
