# Git Experience Helper (gitexp)

A Git configuration and SSH key management utility inspired by Gitea's environment-to-ini tool.

## Overview

This tool helps developers manage their Git configurations and SSH keys more efficiently. It provides a unified interface for common Git setup tasks and supports environment-based configuration.

## Features

- **Git Configuration Management**: Set user name, email, and other Git settings
- **SSH Key Management**: Generate, list, and manage SSH keys for Git repositories
- **Environment-based Configuration**: Configure Git settings from environment variables
- **Both Local and Global Configuration**: Support for repository-specific and global Git settings

## Installation

```bash
# Clone the repository
git clone https://github.com/kashrobinson71-tech/Sshappkey-puppygui4got.git
cd Sshappkey-puppygui4got

# Build the tool
go build -o gitexp main.go

# Or install it
go install
```

## Usage

### Configure Git User

```bash
# Set user name and email globally
./gitexp config user --name "Your Name" --email "your.email@example.com"

# Set user name and email for current repository only
./gitexp config user --name "Your Name" --email "your.email@example.com" --global=false
```

### List Git Configuration

```bash
# Show current Git configuration
./gitexp config list
```

### SSH Key Management

```bash
# List existing SSH keys
./gitexp ssh list

# Generate a new ED25519 SSH key (recommended)
./gitexp ssh generate --email "your.email@example.com"

# Generate an RSA SSH key
./gitexp ssh generate --email "your.email@example.com" --type rsa
```

### Environment-based Configuration

Set environment variables and run:

```bash
# Set environment variables
export GIT_USER_NAME="Your Name"
export GIT_USER_EMAIL="your.email@example.com"
export GIT_DEFAULT_BRANCH="main"

# Apply configuration from environment
./gitexp env
```

## Environment Variables

The tool supports the following environment variables:

- `GIT_USER_NAME`: Sets the Git user name
- `GIT_USER_EMAIL`: Sets the Git user email  
- `GIT_DEFAULT_BRANCH`: Sets the default branch name for new repositories

## Similar to Gitea's environment-to-ini

This tool follows a similar pattern to Gitea's environment-to-ini utility but focuses on Git configuration management rather than general INI file manipulation. It provides:

- Command-line interface using urfave/cli/v3
- Environment variable support
- Structured subcommands for different operations
- Error handling and user feedback

## Examples

### Quick Setup for New Developer

```bash
# Set up Git user configuration from environment
export GIT_USER_NAME="John Doe"
export GIT_USER_EMAIL="john.doe@company.com"
./gitexp env

# Generate SSH key
./gitexp ssh generate --email "john.doe@company.com"

# Verify configuration
./gitexp config list
```

### Project-specific Configuration

```bash
# Navigate to your project
cd /path/to/your/project

# Set project-specific user (useful for work vs personal projects)
./gitexp config user --name "Work Name" --email "work@company.com" --global=false
```

## Requirements

- Go 1.21 or later
- Git installed and available in PATH
- SSH client tools (for SSH key generation)

## License

This project follows the same licensing approach as referenced in the Gitea project.