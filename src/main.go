package main

import (
	"fmt"
	"os"

	"mousehole/src/crypt"
)

const (
	REQUIRED_ARGS          = 4
	REQUIRED_ARGS_IDENTITY = 3
)

func validateFileExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file does not exist")
	}
	return nil
}

func handle_create(args []string) error {
	command := args[0]
	path := args[1]
	if command == "--file" {
		signer := crypt.NewSigner()

		if err := signer.New(); err != nil {
			return err
		}

		if err := signer.Save(path); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not supported sub command")
	}
	return nil
}

func handle_pin(args []string) error {
	command := args[0]
	path := args[1]
	if command == "--file" {
		if err := validateFileExists(path); err != nil {
			return err
		}
		signer := crypt.NewSigner()

		if err := signer.Load(path); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("not supported sub command")
	}
	return nil
}

func handle_identity(argsv []string) error {
	argsc := len(argsv)
	if argsc != REQUIRED_ARGS_IDENTITY {
		return fmt.Errorf("command line arguments are incorrect!")
	}

	switch command := argsv[0]; command {
	case "create":
		return handle_create(argsv[1:])
	case "pin":
		return handle_pin(argsv[1:])
	default:
		return fmt.Errorf("command line argument is not suported")
	}
	return nil
}

func handle_args(argsv []string) error {
	argsc := len(argsv)
	if argsc != REQUIRED_ARGS {
		return fmt.Errorf("command line arguments are incorrect!")
	}

	switch command := argsv[0]; command {
	case "identity":
		return handle_identity(argsv[1:])
		break
	default:
		return fmt.Errorf("command line argument is not suported")
	}
	return nil
}

func main() {
	args := os.Args[1:]

	if err := handle_args(args); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(-1)
	}
}
