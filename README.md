# Peep: AutoReload

Peep is a command line application that monitors specified directories for file changes and executes
user-defined commands, facilitating development automation processes. It allows you to automate
builds, restart programs, and gain insights into command execution through a verbose mode.

## Installation

To install Peep, follow these steps:

1. Clone the repository:

   ```bash
   git clone git@github.com:d3code/peep.git
   ```

2. Navigate to the project directory:

   ```bash
   cd peep
   ```

3. Build and install the application:

   ```bash
   go install
   ```

## Usage

Peep offers various flags to tailor the monitoring behavior and execution of commands. Here are the
available flags:

- `-v, --verbose`: Show additional information about command execution.

- `-d, --directory`: Specify directories to watch (default: current directory).

- `-e, --extension`: Specify extensions to watch (default: all extensions).

- `-i, --ignore`: Specify files or sub-directories to ignore (default: .git).

- `-r, --run`: Specify the command to run and restart on file change.

### Examples

Monitor the current directory with verbose output and run a build command on file changes:

```bash
peep -v -r "go run main.go"
```

Monitor `.go` files in the current and another directory, ignore `.git`, `.idea` and `bin` folders,
and run a custom command:

```bash
peep \
  -d . \
  -d ../other-module \
  -e ".go" \
  -r "go build -o bin/server ./cmd/server;bin/server" \
  -r "cd ../other-module;go build -o bin/server ./cmd/server;bin/server" \
  -i ".git",".idea","bin"
```

Note that you can specify multiple directories and commands to run. Peep will run the commands in
the order they are specified.

### Configuration

> Note: The configuration file feature is currently not implemented.

You can also use a configuration file to specify the directories and commands to run. The
configuration file must be named `peep.yaml` and placed in the directory you want to monitor. Here's
an example configuration file:

```yaml
directories:
  - .
  - ../other-module
extensions:
  - .go
ignore:
  - .git
  - .idea
  - bin
commands:
  - go build -o bin/server ./cmd/server
  - bin/server
```

## Contributing

Contributions to Peep are welcome! If you'd like to contribute, please follow the guidelines
mentioned in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

Peep is licensed under the [MIT License](LICENSE.md).