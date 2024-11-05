# TCP Chat Application in Go (NETCAT)

This project is a NetCat-like TCP chat application built in Go, designed to support server-client communication over TCP with a simple group chat interface. The application can operate in server mode to accept incoming connections on a specified port or in client mode, where it connects to a server to join the chat. It includes concurrency control, message broadcasting, and client connection management.

## Table of Contents

- [About the Project](#about-the-project)
- [Features](#features)
- [Usage](#usage)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage Instructions](#usage-instructions)
- [Error Handling](#error-handling)

---

## About the Project

This application is a NetCat (nc) inspired tool that supports multiple clients in a TCP-based chat server. Each client connecting to the server is prompted for a username, after which they can send messages to all other clients connected to the server. Incoming and outgoing messages are timestamped, and clients are notified when others join or leave the chat. The server maintains a history of messages, which new clients receive upon connection, and supports a maximum of 10 concurrent clients.

## Features

- **TCP Server-Client Architecture**: Allows multiple clients to connect to a single server.
- **Client Identification**: Each client must provide a non-empty name.
- **Connection Control**: Maximum of 10 clients at a time.
- **Message Broadcasting**: All clients can send and receive messages.
- **Timestamped Messages**: Displays the sender’s name and timestamp with each message.
- **Connection Notifications**: Notifies all clients when someone joins or leaves the chat.
- **Message History**: New clients receive previous messages upon connection.
- **Default and Custom Ports**: Listens on port 8989 by default but can be customized.
- **Enhanced with Concurrency**: Uses goroutines and mutexes for concurrent handling.
- **Error Handling**: Robust error handling on both server and client sides.

## Requirements

- Go version 1.15+
- Network access for TCP/IP communication

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/A-fethi/net-cat.git
   cd net-cat
   ```

2. Compile the project:
   ```bash
   go build -o TCPChat
   ```

## Usage

Start the server:

```bash
$ ./TCPChat
Listening on port :8989
```

Or specify a custom port:

```bash
$ ./TCPChat 2525
Listening on port :2525
```

Connect to the server as a client:

```bash
$ nc localhost 8989
```

Example interaction:
```plaintext
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: afethi
[2023-01-20 16:03:43][afethi]: Hello everyone!
```

## Usage Instructions

1. Run the server to start listening for incoming client connections.
2. Each client must enter a unique, non-empty name to join the chat.
3. Messages are broadcast to all clients, including timestamps and usernames.
4. Type a message and press Enter to send.
5. Type `Ctrl+C` to exit the chat.

**Usage Errors**:
If the port is not specified correctly, the program will respond with:
```plaintext
[USAGE]: ./TCPChat $port
```

## Error Handling

- **Server-Side**: Handles client connection errors, limiting connections to 10.
- **Client-Side**: Checks for empty names on connection, prevents sending of empty messages.
- **Network Disruptions**: Connections automatically handle leaving clients, notifying remaining users.

---

**Note**: This application stores a chat history in memory and does not persist it across sessions.