
# WebSocket Chat Server

## Overview

This project implements a WebSocket-based chat server in Go using the Gorilla WebSocket package. The server supports basic message broadcasting between clients and is designed to handle multiple clients simultaneously. Each client is identified by a unique `userId`.

## Features

- **WebSocket Communication**: Real-time messaging using WebSocket connections.
- **Client Management**: Register and unregister clients dynamically.
- **Message Broadcasting**: Send messages to specific clients.
- **Graceful Handling of Disconnections**: Handles unexpected client disconnections and cleans up resources.

## Planned Features

- **Redis Support**: Planned integration with Redis for handling client presence and message queuing, enabling a scalable and distributed system.

## Project Structure

- `internal/Client.go`: Manages individual client connections, handling both incoming and outgoing messages.
- `internal/Hub.go`: Manages all connected clients, handling registration, unregistration, and message broadcasting.

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or higher)
- [Redis](https://redis.io/download) (if planning to use Redis support)

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/tarun-kavipurapu/chat-app.git
   cd chat-app
   ```

2. **Install Dependencies**

   This project uses Gorilla WebSocket. To install the package, run:

   ```bash
   go get github.com/gorilla/websocket
   ```

3. **Run the Server**

   ```bash
   go run main.go
   ```

## Usage

### Establishing a WebSocket Connection

Clients can connect to the server using the WebSocket protocol. The server listens for connections and handles communication based on the client-defined `userId`.

### Sending and Receiving Messages

Messages are sent and received in JSON format. Each message should include the following fields:

- `from`: The sender's userId
- `to`: The recipient's userId
- `content`: The message content

### Example Message Structure

```json
{
  "from": "user1",
  "to": "user2",
  "content": "Hello, World!"
}
```

## Redis Integration (Planned)

To enable Redis support:

1. **Install Redis Client Library**

   ```bash
   go get github.com/go-redis/redis/v8
   ```

2. **Configure Redis in the Project**

   Modify the project to use Redis for storing active clients and message queuing. This will allow for scalability and distribution of the WebSocket server across multiple instances.

## TODO

- Implement Redis-based message queuing and presence tracking.
- Develop a more robust client reconnection strategy.
- Add support for room-based messaging.
- Implement logging and monitoring features.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

=