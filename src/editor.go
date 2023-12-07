package main

import (
	"os"
	"os/exec"
)

func EditJournalEntry() (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "journal-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file afterwards

	// Function to try opening an editor
	openEditor := func(editor string) error {
		cmd := exec.Command(editor, tmpFile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Try NeoVim first, then fall back to Vim
	if err := openEditor("nvim"); err != nil {
		if err := openEditor("vim"); err != nil {
			return "", err // Both NeoVim and Vim failed
		}
	}

	// Read the content of the temporary file
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return string(content), nil
}
