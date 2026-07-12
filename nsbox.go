// nsbox - A non-sandboxed command executor with compatible interface to trae-sandbox

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

const version = "1.0.0"

func printHelp() {
	fmt.Print(`A lightweight command executor (non-sandboxed)

Usage: nsbox <COMMAND>

Commands:
  exec    Run command directly without sandbox isolation
  help    Print this message or the help of the given subcommand(s)

Options:
  -h, --help     Print help
  -v, --version  Print version
`)
}

func printExecHelp() {
	fmt.Print(`Run command directly without sandbox isolation

Usage: nsbox exec [OPTIONS] --storage-path <STORAGE_PATH> --config-name <CONFIG_NAME> --shell-path <SHELL_PATH> --command-line <COMMAND_LINE>

Options:
      --storage-path <STORAGE_PATH>  Configuration file storage directory (accepted but ignored)
      --config-name <CONFIG_NAME>    Configuration file name (accepted but ignored)
      --shell-path <SHELL_PATH>      Shell path to use (e.g., "/bin/bash", "/bin/zsh")
      --command-line <COMMAND_LINE>  Command line arguments for the process to start
  -h, --help                         Print help
`)
}

type execConfig struct {
	storagePath string
	configName  string
	shellPath   string
	commandLine string
	dryRun      bool
}

func parseExecArgs(args []string) (*execConfig, error) {
	cfg := &execConfig{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--storage-path":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--storage-path requires a value")
			}
			i++
			cfg.storagePath = args[i]
		case "--config-name":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--config-name requires a value")
			}
			i++
			cfg.configName = args[i]
		case "--shell-path":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--shell-path requires a value")
			}
			i++
			cfg.shellPath = args[i]
		case "--command-line":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--command-line requires a value")
			}
			i++
			cfg.commandLine = args[i]
		case "--dry-run":
			cfg.dryRun = true
		case "-h", "--help":
			printExecHelp()
			os.Exit(0)
		default:
			return nil, fmt.Errorf("unknown option: %s", args[i])
		}
	}
	if cfg.shellPath == "" {
		return nil, fmt.Errorf("--shell-path is required")
	}
	if cfg.commandLine == "" {
		return nil, fmt.Errorf("--command-line is required")
	}
	return cfg, nil
}

func runExec(args []string) error {
	cfg, err := parseExecArgs(args)
	if err != nil {
		return err
	}

	if cfg.dryRun {
		fmt.Printf("[dry-run] Would execute: %s -c \"%s\"\n", cfg.shellPath, cfg.commandLine)
		return nil
	}

	cmd := exec.Command(cfg.shellPath, "-c", cfg.commandLine)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sigChan)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	go func() {
		for sig := range sigChan {
			if cmd.Process != nil {
				cmd.Process.Signal(sig)
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "exec":
		if err := runExec(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		printHelp()
	case "--version", "-v":
		fmt.Printf("nsbox %s\n", version)
	default:
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		commandLine := strings.Join(os.Args[1:], " ")
		cmd := exec.Command(shell, "-c", commandLine)
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
