package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	. "github.com/tedraykov/devtools/scripts"
)

type MacOsTools struct {
	tools []string
}

func (t *MacOsTools) Run() error {
	// Ensure Homebrew is installed
	if err := t.ensureHomebrew(); err != nil {
		return err
	}

	t.runCommand("brew", "update")

	for _, tool := range t.tools {
		color.Blue("Setting up %s...", tool)
		var err error

		switch tool {
		case "zsh":
			err = t.InstallZsh()
		case "go":
			err = t.InstallGo()
		case "make":
			err = t.InstallMake()
		case "gcc":
			err = t.InstallGcc()
		case "unzip":
			err = t.InstallUnzip()
		case "docker":
			err = t.InstallDocker()
		case "tmux":
			err = t.InstallTmux()
		case "node":
			err = t.InstallNode()
		case "python":
			err = t.InstallPython()
		case "neovim":
			err = t.InstallNeovim()
		default:
			color.Red("Error: %s is not a valid tool", tool)
		}

		if err != nil {
			color.Red("Error installing %s: %v", tool, err)
			return err
		}

		color.Green("%s installed and configured successfully", tool)
	}

	return nil
}

func (m *MacOsTools) ensureHomebrew() error {
	_, err := exec.LookPath("brew")
	if err == nil {
		return nil // Homebrew is already installed
	}

	color.Blue("Installing Homebrew...")
	return m.runCommand( "curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh")
}

func (m *MacOsTools) InstallNeovim() error {
	color.Blue("Installing Neovim...")
	if err := m.runCommand("brew", "install", "neovim"); err != nil {
		return err
	}

	return m.ConfigureNeovim()
}

func (m *MacOsTools) InstallZsh() error {
	color.Blue("Installing Zsh...")
	if err := m.runCommand("brew", "install", "zsh"); err != nil {
		return err
	}

	// Change default shell to Zsh
	if err := m.runCommand("chsh", "-s", "/bin/zsh"); err != nil {
		return err
	}

	// Install Oh My Zsh
	return m.InstallOhMyZsh()
}

func (m *MacOsTools) InstallGcc() error {
	color.Blue("Installing GCC...")
	return m.runCommand("brew", "install", "gcc")
}

func (m *MacOsTools) InstallMake() error {
	color.Blue("Installing Make...")
	return m.runCommand("brew", "install", "make")
}

func (m *MacOsTools) InstallRipgrep() error {
	color.Blue("Installing Ripgrep...")
	return m.runCommand("brew", "install", "ripgrep")
}

func (m *MacOsTools) InstallUnzip() error {
	color.Blue("Installing Unzip...")
	return m.runCommand("brew", "install", "unzip")
}

func (m *MacOsTools) InstallOhMyZsh() error {
	color.Blue("Installing Oh My Zsh...")
	return m.runCommand("curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh")
}

func (m *MacOsTools) InstallDocker() error {
	color.Blue("Installing Docker...")
	return m.runCommand("brew", "install", "--cask", "docker")
}

func (m *MacOsTools) InstallTmux() error {
	color.Blue("Installing Tmux...")
	if err := m.runCommand("brew", "install", "tmux"); err != nil {
		return err
	}

	color.Blue("Installing fzf...")
	if err := m.runCommand("brew", "install", "fzf"); err != nil {
		return err
	}

	tmuxSessionizerScriptPath := filepath.Join(LocalBinPath(), "tmux-sessionizer")
	tmuxConfigPath := filepath.Join(HomePath(), ".tmux.conf")

	color.Blue("Installing tmux-sessionizer script...")
	if err := WriteContentToFile(TmuxSessionizer, tmuxSessionizerScriptPath); err != nil {
		return err
	}

	if err := MakeExecutable(tmuxSessionizerScriptPath); err != nil {
		return err
	}

	color.Blue("Configuring tmux...")
	return WriteContentToFile(TmuxConfig, tmuxConfigPath)
}

func (m *MacOsTools) getLastestGoVersion() (string, error) {
	url := "https://golang.org/VERSION?m=text"
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	version := strings.TrimSpace(string(body))
	re := regexp.MustCompile(`go(\d+\.\d+(\.\d+)?)`)
	match := re.FindStringSubmatch(version)
	if len(match) < 2 {
		return "", fmt.Errorf("failed to parse version number from: %s", version)
	}

	return match[1], nil
}

func (m *MacOsTools) InstallGo() error {
	color.Blue("Installing Go...")
	return m.runCommand("brew", "install", "go")
}

func (m *MacOsTools) InstallNode() error {
	color.Blue("Installing NVM...")
	return m.runCommand("brew", "install", "nvm")
}

func (m *MacOsTools) InstallPython() error {
	color.Blue("Installing Python...")
	return m.runCommand("brew", "install", "python")
}

func (m *MacOsTools) ConfigureNeovim() error {
	color.Blue("Configuring Neovim...")
	return m.runCommand("git", "clone", "https://github.com/tedraykov/init.lua.git", "~/.config/nvim")
}

func (m *MacOsTools) InstallPoetry() error {
  color.Blue("Installing Poetry...")
  return m.runCommand("curl -sSL https://install.python-poetry.org | python3 -")
}

func (m *MacOsTools) runCommand(args ...string) error {
	splitArgs := []string{}

	for _, arg := range args {
		splitArgs = append(splitArgs, strings.Split(arg, " ")...)
	}

	fmt.Printf("Running command: %s\n", strings.Join(splitArgs, " "))
	cmd := exec.Command(splitArgs[0], splitArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
