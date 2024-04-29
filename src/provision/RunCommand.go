package provision

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Exec(name string, arg ...string) (*exec.Cmd, error) {
	// Command to start the package manager (e.g., pacman or apt-get)
	cmd := exec.Command(name, arg...) // Change this command as needed

	// Create pipes for stdin, stdout, and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return cmd, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return cmd, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return cmd, err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return cmd, err
	}

	// Create a scanner to read outputs from the command
	outStream := bufio.NewScanner(stdout)
	errStream := bufio.NewScanner(stderr)

	// Start a goroutine to read and print output from the command
	go func() {
		for outStream.Scan() {
			fmt.Println(outStream.Text())
			if err != nil {
				return
			}
		}
	}()

	// Start a goroutine to read and print output from the command
	go func() {
		for errStream.Scan() {
			fmt.Println(errStream.Text())
		}
	}()

	// Start a goroutine to capture keyboard input and write it to stdin
	go keyboardListener(stdin)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
	}

	return cmd, nil
}

func keyboardListener(stdin io.WriteCloser) {
	keyboard := bufio.NewScanner(os.Stdin)
	for keyboard.Scan() {
		input := keyboard.Text()
		_, err := fmt.Fprintln(stdin, input)
		if err != nil {
			fmt.Println("Error writing to stdin:", err)
			return
		}
	}
}
