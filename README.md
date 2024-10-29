# Mailman

Mailman is a simple client application for generating html emails from markdown files and sending them using an email server.

# Dependencies

This project uses the following dependencies:

- [Go 1.21.5](https://go.dev/doc/install)

# Installation

Install all go dependencies by running the following command.

```bash
go install
```

# Usage in Go

You need to define a config file in yaml format and a markdown file as email. View the example config files in the [examples](examples) directory.

Shows the help message for the application. 
```bash
go run mailman -h
```

Generates a custom html file from a markdown File and a html template.

```bash
go run mailman --config ./examples/ua-service.yaml generate
```

Send the custom email as html file.

```bash
go run mailman --config ./examples/ua-service.yaml send
```

# Build

Build the project by running the following command for a desired platform.

```bash
export GOOS=linux
export GOARCH=amd64
go build -o mailman-linux-amd64
```

# Usage as Executable

Copy the executable to the desired location and make it executable.

```bash
chmod +x mailman-linux-amd64
mv mailman-linux-amd64 ~/.local/bin/mailman
```

Shows the help message for the application. 
```bash
mailman -h
```



