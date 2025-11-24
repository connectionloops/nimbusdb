# Local Setup

This guide will help you set up the NimbusDb project on your local machine.

## Prerequisites

### Go Installation

This project requires **Go 1.25** or later.

#### Recommended: Using asdf (Version Manager)

If you work with multiple Go projects that require different Go versions, we recommend using [asdf](https://asdf-vm.com/) to manage Go versions.

##### Installing asdf

**macOS:**

```bash
brew install asdf
```

**Linux:**

```bash
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0
echo '. "$HOME/.asdf/asdf.sh"' >> ~/.bashrc
echo '. "$HOME/.asdf/completions/asdf.bash"' >> ~/.bashrc
source ~/.bashrc
```

**Windows (using WSL/Git Bash):**

```bash
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0
echo '. "$HOME/.asdf/asdf.sh"' >> ~/.bashrc
echo '. "$HOME/.asdf/completions/asdf.bash"' >> ~/.bashrc
source ~/.bashrc
```

##### Installing the Go Plugin

```bash
asdf plugin add golang
```

##### Basic asdf Usage

```bash
# Install a specific Go version
asdf install golang 1.25.0

# Set the global Go version (applies to all projects)
asdf global golang 1.25.0

# Set the local Go version (for this project only)
asdf local golang 1.25.0

# List installed Go versions
asdf list golang

# List all available Go versions
asdf list all golang

# Show current Go version
asdf current golang

# Update Go plugin to see latest versions
asdf plugin update golang
```

This project already has a `.tool-versions` file (which asdf will automatically pickup).

##### Verify Go Installation

From project root,

```bash
go version
```

Should output something like: `go version go1.25.0 darwin/amd64`

#### Alternative: Direct Go Installation

If you prefer a direct installation or don't need version management:

**macOS:**

```bash
brew install go
```

**Linux:**

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y golang-go

# Fedora
sudo dnf install -y golang
```

**Windows:**

1. Download the installer from [golang.org/dl](https://go.dev/dl/)
2. Run the installer and follow the prompts
3. Verify installation: `go version`

**From Source:**
Follow the instructions at [golang.org/doc/install/source](https://go.dev/doc/install/source)

## Project Setup

### Clone the Repository

```bash
git clone <repository-url>
cd NimbusDb
```

### Install Dependencies

The project uses Go modules. Dependencies will be automatically downloaded when you build or test the project.

To explicitly download dependencies:

```bash
go mod download
```

To verify dependencies are correct:

```bash
go mod verify
```

### Configuration & Secrets

- A sample configuration file is given `.config.yml`
- However secrets are recommended to be injected by a secret manager
- If you are a Connection Loops employee then doppler file is already added for you. just run `doppler setup` to set it up.
- To run the project with secrets being injected from doppler, run `doppler run -- go run .`

### Build the Project

```bash
go build
```

This will create an executable in the current directory.

### Run the Application

```bash
go run main.go
```

Or if you've built it:

```bash
./NimbusDb
```

> Please make sure to supply the needed secrets see Configuration & Secrets section above.

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests in a Specific Package

```bash
go test ./configurations
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
# Run tests and generate coverage report
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

This will open a coverage report in your browser showing which lines are covered by tests.

### Run Specific Test Functions

```bash
# Run a specific test function
go test -run TestLoad_FromYAML ./configurations

# Run tests matching a pattern
go test -run TestLoad ./configurations
```

### Run Tests in Watch Mode

For continuous testing during development, you can use tools like:

**Using `entr` (Linux/macOS):**

```bash
# Install entr
brew install entr  # macOS
# or: sudo apt-get install entr  # Linux

# Watch and run tests
find . -name "*.go" | entr -c go test ./...
```

**Using `air` (with test configuration):**
Install [air](https://github.com/cosmtrek/air) for a more full-featured testing experience.

## Format with `gofmt`

Run

```bash
gofmt -l -w .
```

to auto-format code.

## Troubleshooting

WIP

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [asdf Documentation](https://asdf-vm.com/guide/getting-started.html)
- [Go Testing Documentation](https://go.dev/doc/tutorial/add-a-test)
