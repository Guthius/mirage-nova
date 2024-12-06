![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Guthius/mirage-nova/go.yml)
![GitHub License](https://img.shields.io/github/license/Guthius/mirage-nova)

# Mirage Nova

**Mirage Nova** is a modernized 2D game development engine written in Go (Golang) designed for creating online multiplayer role-playing games (MMORPGs).

It leverages Go's powerful concurrency model to enable efficient and scalable networking.

## Features

*TBC*

## Build Instructions

### Clone the repository
```bash
git clone https://github.com/Guthius/mirage-nova.git
cd mirage-nova
```

### Install Go

Ensure that [Go](https://golang.org/dl/) is installed on your system. The minimum required version of Go is 1.23.

You can verify the installation by running:
```bash
go version
```

### Build the server

Navigate to the server directory and build the server:
```bash
cd server
go build -o mirage-nova-server
```

### Run the server

Execute the built server binary:
```bash
./mirage-nova-server
```

## Legacy

The Mirage Nova Engine represents a significant step forward from its predecessor, the VB6-based Mirage Engine, by embracing modern technologies and best practices.

Mirage Nova is based on the Lite 2D Version of the VB6 Mirage Engine available at [https://mirage-engine.uk/](https://mirage-engine.uk/).

## License

This project is licensed under the MIT License. For the complete license text, please refer to the [LICENSE](LICENSE) file.