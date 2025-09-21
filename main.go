// Git Configuration Helper
// A utility to help manage Git configurations and SSH keys
// Inspired by Gitea's environment-to-ini tool

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "gitexp",
		Usage: "Git Experience Helper - Manage Git configurations and SSH keys",
		Description: `A helper utility to manage Git configurations, SSH keys, and provide
a better Git experience. This tool can help with:

- Setting up Git user configuration
- Managing SSH keys for Git repositories
- Configuring Git for different environments
- Environment-based Git configuration

Similar to Gitea's environment-to-ini, this tool provides a convenient
way to manage Git settings through environment variables and direct configuration.`,
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "Manage Git configuration",
				Commands: []*cli.Command{
					{
						Name:   "user",
						Usage:  "Configure Git user settings",
						Action: configureUser,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Usage:   "Set Git user name",
							},
							&cli.StringFlag{
								Name:    "email",
								Aliases: []string{"e"},
								Usage:   "Set Git user email",
							},
							&cli.BoolFlag{
								Name:  "global",
								Usage: "Set configuration globally",
								Value: true,
							},
						},
					},
					{
						Name:   "list",
						Usage:  "List current Git configuration",
						Action: listConfig,
					},
				},
			},
			{
				Name:  "ssh",
				Usage: "Manage SSH keys for Git",
				Commands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "List SSH keys",
						Action: listSSHKeys,
					},
					{
						Name:   "generate",
						Usage:  "Generate new SSH key",
						Action: generateSSHKey,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "email",
								Aliases: []string{"e"},
								Usage:   "Email for SSH key",
							},
							&cli.StringFlag{
								Name:  "type",
								Usage: "SSH key type (rsa, ed25519)",
								Value: "ed25519",
							},
						},
					},
				},
			},
			{
				Name:   "env",
				Usage:  "Configure Git from environment variables",
				Action: configureFromEnv,
				Description: `Configure Git settings from environment variables.
Supported environment variables:
- GIT_USER_NAME: Set Git user name
- GIT_USER_EMAIL: Set Git user email
- GIT_DEFAULT_BRANCH: Set default branch name`,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func configureUser(ctx context.Context, c *cli.Command) error {
	name := c.String("name")
	email := c.String("email")
	global := c.Bool("global")

	scope := "--local"
	if global {
		scope = "--global"
	}

	if name != "" {
		cmd := exec.Command("git", "config", scope, "user.name", name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set user name: %w", err)
		}
		fmt.Printf("Set user name to: %s\n", name)
	}

	if email != "" {
		cmd := exec.Command("git", "config", scope, "user.email", email)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set user email: %w", err)
		}
		fmt.Printf("Set user email to: %s\n", email)
	}

	if name == "" && email == "" {
		return fmt.Errorf("please provide at least one of --name or --email")
	}

	return nil
}

func listConfig(ctx context.Context, c *cli.Command) error {
	cmd := exec.Command("git", "config", "--list")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list Git configuration: %w", err)
	}

	fmt.Println("Current Git Configuration:")
	fmt.Println(string(output))
	return nil
}

func listSSHKeys(ctx context.Context, c *cli.Command) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return fmt.Errorf("failed to read SSH directory: %w", err)
	}

	fmt.Println("SSH Keys:")
	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".pub") || 
			(!strings.Contains(entry.Name(), ".") && !strings.HasSuffix(entry.Name(), ".pub"))) {
			fmt.Printf("  %s\n", entry.Name())
		}
	}

	return nil
}

func generateSSHKey(ctx context.Context, c *cli.Command) error {
	email := c.String("email")
	keyType := c.String("type")

	if email == "" {
		return fmt.Errorf("email is required for SSH key generation")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("failed to create SSH directory: %w", err)
	}

	keyPath := filepath.Join(sshDir, fmt.Sprintf("id_%s", keyType))
	
	var cmd *exec.Cmd
	switch keyType {
	case "ed25519":
		cmd = exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", keyPath, "-N", "")
	case "rsa":
		cmd = exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-C", email, "-f", keyPath, "-N", "")
	default:
		return fmt.Errorf("unsupported key type: %s", keyType)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate SSH key: %w", err)
	}

	fmt.Printf("SSH key generated successfully!\n")
	fmt.Printf("Private key: %s\n", keyPath)
	fmt.Printf("Public key: %s.pub\n", keyPath)
	
	// Display the public key
	pubKeyContent, err := os.ReadFile(keyPath + ".pub")
	if err == nil {
		fmt.Printf("\nPublic key content:\n%s\n", string(pubKeyContent))
	}

	return nil
}

func configureFromEnv(ctx context.Context, c *cli.Command) error {
	configured := false

	if name := os.Getenv("GIT_USER_NAME"); name != "" {
		cmd := exec.Command("git", "config", "--global", "user.name", name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set user name from environment: %w", err)
		}
		fmt.Printf("Set user name from GIT_USER_NAME: %s\n", name)
		configured = true
	}

	if email := os.Getenv("GIT_USER_EMAIL"); email != "" {
		cmd := exec.Command("git", "config", "--global", "user.email", email)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set user email from environment: %w", err)
		}
		fmt.Printf("Set user email from GIT_USER_EMAIL: %s\n", email)
		configured = true
	}

	if branch := os.Getenv("GIT_DEFAULT_BRANCH"); branch != "" {
		cmd := exec.Command("git", "config", "--global", "init.defaultBranch", branch)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to set default branch from environment: %w", err)
		}
		fmt.Printf("Set default branch from GIT_DEFAULT_BRANCH: %s\n", branch)
		configured = true
	}

	if !configured {
		fmt.Println("No Git environment variables found.")
		fmt.Println("Supported variables: GIT_USER_NAME, GIT_USER_EMAIL, GIT_DEFAULT_BRANCH")
	}

	return nil
}