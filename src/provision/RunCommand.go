package provision

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func Exec(arg ...string) (*exec.Cmd, error) {
	// Command to start the package manager (e.g., pacman or apt-get)
	cmd := exec.Command(arg[0], arg[1:]...) // Change this command as needed

	// Create pipes for stdin, stdout, and stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal("Error creating stdin pipe:", err)
		return cmd, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating stdout pipe:", err)
		return cmd, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error creating stderr pipe:", err)
		return cmd, err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting command:", err)
		return cmd, err
	}

	// Create a scanner to read outputs from the command
	outStream := bufio.NewScanner(stdout)
	errStream := bufio.NewScanner(stderr)

	// Start a goroutine to read and print output from the command
	go func() {
		for outStream.Scan() {
			log.Println(outStream.Text())
			if err != nil {
				return
			}
		}
	}()

	// Start a goroutine to read and print output from the command
	go func() {
		for errStream.Scan() {
			log.Println(errStream.Text())
		}
	}()

	// Start a goroutine to capture keyboard input and write it to stdin
	go keyboardListener(stdin)

	// Wait for the command to finish
	//if err := cmd.Wait(); err != nil {
	//	log.Fatal("Error waiting for command:", err)
	//}

	return cmd, cmd.Wait()
}

func keyboardListener(stdin io.WriteCloser) {
	keyboard := bufio.NewScanner(os.Stdin)
	for keyboard.Scan() {
		input := keyboard.Text()
		_, err := fmt.Fprintln(stdin, input)
		if err != nil {
			log.Fatal("Error writing to stdin:", err)
			return
		}
	}
}
