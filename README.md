# Run: AutoReload

Run is a command line application that monitors specified directories for file changes and executes
user-defined commands, facilitating development automation processes. It allows you to automate
builds, restart programs, and gain insights into command execution through a verbose mode.

## Installation

To install Run, follow these steps:

1. Clone the repository:

   ```bash
   git clone git@github.com:d3code/run.git
   ```

2. Navigate to the project directory:

   ```bash
   cd run
   ```

3. Build and install the application:

   ```bash
   go install ./cmd/run
   ```

## Usage

Run offers various flags to tailor the monitoring behavior and execution of commands. Here are the
available flags:

- `-v, --verbose`: Show additional information about command execution.

- `-d, --directory`: Specify directories to watch (default: current directory).

- `-e, --extension`: Specify extensions to watch (default: all extensions).

- `-i, --ignore`: Specify files or sub-directories to ignore (default: .git).

- `-c, --command`: Specify the command to run and restart on file change.

### Examples

Monitor the current directory with verbose output and run a build command on file changes:

```bash
run -v -r "go run main.go"
```

Monitor `.go` files in the current and another directory, ignore `.git`, `.idea` and `bin` folders,
and run a custom command:

```bash
run \
  -d . \
  -d ../other-module \
  -e ".go" \
  -r "go build -o bin/server ./cmd/server;bin/server" \
  -r "cd ../other-module;go build -o bin/server ./cmd/server;bin/server" \
  -i ".git",".idea","bin"
```

Note that you can specify multiple directories and commands to run. Run will execute the commands in
the order they are specified.

## Contributing

Contributions to Run are welcome! If you'd like to contribute, please follow the guidelines
mentioned in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## License

Run is licensed under the [MIT License](LICENSE.md).
