<div align="center">
    <img src=".github/assets/logo.png" width="400">

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/guthius/mirage-nova)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/Guthius/mirage-nova/go.yml)
![CodeFactor Grade](https://img.shields.io/codefactor/grade/github/guthius/mirage-nova)
![GitHub License](https://img.shields.io/github/license/Guthius/mirage-nova)

</div>

---

**Mirage Nova** is a 2D game development engine written in Go (Golang), specifically designed for creating small-scale online multiplayer role-playing games.

It leverages Go's powerful concurrency model to enable efficient and scalable networking.

Mirage Nova is based on the Lite 2D Version of the VB6 Mirage Engine available at [https://mirage-engine.uk/](https://mirage-engine.uk/).

## Features

- Player account management and login system
- Account management with multiple characters per account
- Real-time multiplayer communication with efficient TCP networking

The following features are currently in development:

- In-game chat system (global, private, and local)
- Item inventory system with support for trading and equipping items
- NPC spawning, movement, and interactions
- Basic combat mechanics
- Experience and leveling system for player progression

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
go build -o ../bin/
```

### Run the server

Execute the built server binary:
```bash
cd ../bin
./server
```

## License

This project is licensed under the MIT License. For the complete license text, please refer to the [LICENSE](LICENSE) file.