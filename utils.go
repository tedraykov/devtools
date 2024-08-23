package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func  WriteContentToFile(content, path string) error {
  // Create the file if it doesn't exist
  if _, err := os.Stat(path); os.IsNotExist(err) {
    if _, err := os.Create(path); err != nil {
      return fmt.Errorf("failed to create file: %w", err)
    }
  }

  // Write the content to the file overwriting any existing content
  if err := os.WriteFile(path, []byte(content), 0644); err != nil {
    return fmt.Errorf("failed to write to file: %w", err)
  }

  return nil
}

func HomePath() string {
  return os.Getenv("HOME")
}

func LocalBinPath() string {
  return filepath.Join(HomePath(), ".local", "bin")
}

func MakeExecutable(path string) error {
  if err := os.Chmod(path, 0755); err != nil {
    return fmt.Errorf("failed to make file executable: %w", err)
  }

  return nil
}

func AddToRCFiles(content string) error {
  // Check both the bashrc and zshrc files in the home path
  rcFiles := []string{".bashrc", ".zshrc"}

	for _, rcFile := range rcFiles {
		filePath := filepath.Join(HomePath(), rcFile)

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("%s does not exist, skipping.\n", rcFile)
			continue
		}

		// Check if content is already in the file
		if exists, err := ContentExists(filePath, content); err != nil {
			return fmt.Errorf("error checking %s: %w", rcFile, err)
		} else if exists {
			fmt.Printf("Content already exists in %s, skipping.\n", rcFile)
			continue
		}

		// Append content to the file
		if err := AppendToFile(filePath, content); err != nil {
			return fmt.Errorf("failed to append to %s: %w", rcFile, err)
		}

		fmt.Printf("Successfully added content to %s\n", rcFile)
  }

  return nil
}

func ContentExists(filePath, content string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(content) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func AppendToFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content + "\n"); err != nil {
		return err
	}

	return nil
}

func DeleteFile(path string) error {
  if err := os.Remove(path); err != nil {
    return fmt.Errorf("failed to delete file: %w", err)
  }

  return nil
}
