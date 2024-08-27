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



type UbuntuTools struct {
	tools []string
}

func (t *UbuntuTools) Run() error {
  t.runCommand("sudo", "apt", "update")

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

func (u *UbuntuTools) InstallNeovim() error {
	color.Blue("Removing Vim if installed...")
	if err := u.runCommand("sudo apt remove -y vim vim-runtime gvim"); err != nil {
		return err
	}

	color.Blue("Downloading latest Neovim...")
	if err := u.runCommand("curl -LO https://github.com/neovim/neovim/releases/download/0.9.5/nvim-linux64.tar.gz"); err != nil {
		return err
	}

	color.Blue("Extracting Neovim...")
	if err := u.runCommand("tar xzf nvim-linux64.tar.gz"); err != nil {
		return err
	}

	color.Blue("Moving and renaming Neovim executable to /usr/local/bin/vim...")
	if err := u.runCommand("sudo mv nvim-linux64/bin/nvim /usr/local/bin/vim"); err != nil {
		return err
	}

	color.Blue("Setting correct permissions...")
	if err := u.runCommand("sudo chmod +x /usr/local/bin/vim"); err != nil {
		return err
	}

	color.Blue("Cleaning up...")
	if err := u.runCommand("rm -rf nvim-linux64 nvim-linux64.tar.gz"); err != nil {
		return err
	}

	color.Green("Neovim installation complete. You can now use 'vim' to run Neovim.")

  return u.ConfigureNeovim()
}

func (u *UbuntuTools) InstallZsh() error {
    fmt.Println("Installing Zsh...")
    if err := u.runCommand("sudo apt install -y zsh"); err != nil {
        return err
    }

    // Change default shell to Zsh
    if err := u.runCommand("sudo chsh -s /bin/zsh"); err != nil {
        return err
    }

    // Install Oh My Zsh
    if err := u.InstallOhMyZsh(); err != nil {
        return err
    }

    return nil
}

func (u *UbuntuTools) InstallGcc() error {
    fmt.Println("Installing GCC...")
    return u.runCommand("sudo apt install -y gcc")
}

func (u *UbuntuTools) InstallMake() error {
    fmt.Println("Installing Make...")
    return u.runCommand("sudo apt install -y make")
}

func (u *UbuntuTools) InstallRipgrep() error {
    fmt.Println("Installing Ripgrep...")
    return u.runCommand("sudo apt install -y ripgrep")
}

func (u *UbuntuTools) InstallUnzip() error {
    fmt.Println("Installing Unzip...")
    return u.runCommand("sudo apt install -y unzip")
}

func (u *UbuntuTools) InstallOhMyZsh() error {
    fmt.Println("Installing Oh My Zsh...")
    return u.runCommand("curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh")
}

func (u *UbuntuTools) InstallDocker() error {
    fmt.Println("Installing Docker...")
    cmds := [][]string{
        {"sudo apt-get install ca-certificates curl"},
        {"sudo install -m 0755 -d /etc/apt/keyrings"},
        {"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc"},
        {"sudo chmod a+r /etc/apt/keyrings/docker.asc"},
        {"echo \"deb [signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null"},
        {"sudo apt-get update"},
        {"sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin"},
    }

    for _, cmd := range cmds {
        if err := u.runCommand(cmd...); err != nil {
            return err
        }
    }

    fmt.Println("Running Docker post-installation configuration...")

    cmds = [][]string{
        {"sudo groupadd docker"},
        {"sudo usermod -aG docker $USER"},
        {"newgrp docker"},
    }

    for _, cmd := range cmds {
        if err := u.runCommand(cmd...); err != nil {
          color.Red("Error running Docker post-installation configuration: %v", err)
        }
    }

    return nil
}


func (u *UbuntuTools) InstallTmux() error {
    fmt.Println("Installing Tmux...")
    u.runCommand("sudo apt install -y tmux")


    fmt.Println("Installing fzf...")
    u.runCommand("sudo apt install -y fzf")


    tmuxSessionizerScriptPath := filepath.Join(LocalBinPath(), "tmux-sessionizer")
    tmuxConfigPath := filepath.Join(HomePath(), ".tmux.conf")

    fmt.Println("Installing tmux-sessionizer script...")
    if err := WriteContentToFile(TmuxSessionizer, tmuxSessionizerScriptPath); err != nil {
        return err
    }

    if err := MakeExecutable(tmuxSessionizerScriptPath); err != nil {
        return err
    }

    fmt.Println("Configuring tmux...")
    if err := WriteContentToFile(TmuxConfig, tmuxConfigPath); err != nil {
        return err
  }

    return nil
}

func (u *UbuntuTools) getLastestGoVersion() (string, error) {
	// URL to fetch the latest Go version
	url := "https://golang.org/VERSION?m=text"

	// Make an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert body to string and trim any whitespace
	version := strings.TrimSpace(string(body))

	// Use regex to extract just the version number
	re := regexp.MustCompile(`go(\d+\.\d+(\.\d+)?)`)
	match := re.FindStringSubmatch(version)
	if len(match) < 2 {
		return "", fmt.Errorf("failed to parse version number from: %s", version)
	}

	// Return the version number (without the 'go' prefix)
	return match[1], nil
}

func (u *UbuntuTools) InstallGo() error {
    // Get latest version of Go
    version, err := u.getLastestGoVersion(); if err != nil {
        return err
    }
    filename := fmt.Sprintf("go%s.linux-amd64.tar.gz", version)

    fmt.Printf("Installing Go %s...\n", version)

    // Construct the download URL
    downloadURL := fmt.Sprintf("https://go.dev/dl/%s",filename)

    // Download the Go tarball
    if err = u.runCommand("curl", "-LO", downloadURL); err != nil {
        return  err
    }

    // Remove existing Go installation and extract the new one
    err = u.runCommand("sudo rm -rf /usr/local/go")
    err = u.runCommand("sudo tar -C /usr/local -xzf", filename)

    // Add Go binary directory to PATH
    if err := AddToRCFiles("export PATH=$PATH:/usr/local/go/bin"); err != nil {
        return err
    }

    // Clean up
    if err := DeleteFile(filename); err != nil {
      return err
    }

    return nil
}

func (u *UbuntuTools) InstallNode() error {
    fmt.Println("Installing NVM...")
    return u.runCommand("curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh")
}

func (u *UbuntuTools) InstallPython() error {
    fmt.Println("Installing Python...")
    return u.runCommand("sudo apt install -y python3")
}

func (u *UbuntuTools) ConfigureNeovim() error {
    color.Blue("Configuring Neovim...")

    if err := u.runCommand("git clone https://github.com/tedraykov/init.lua.git ~/.config/nvim"); err != nil {
        return err
    }

    return nil
}

func (u *UbuntuTools) InstallPoetry() error {
    fmt.Println("Installing Poetry...")
    return u.runCommand("curl -sSL https://install.python-poetry.org | python3 -")
}

func (u *UbuntuTools) runCommand(args ...string) error {
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
