# Chatgpt clone with Ollama

This repository contains a simple Go web server application demonstrating the use of various Go packages for web development, including chi for routing, godotenv for environment variable management, and custom packages for business logic.


## Features

- Routing: Utilizes the chi router for handling web requests with middleware support.
- Environment Variables: Uses godotenv to load environment variables from a .env file, facilitating easy configuration management.
- Graceful Shutdown: Supports graceful shutdown on receiving interrupt signals, ensuring that ongoing requests are not abruptly terminated.
Custom Business 
- Logic: Incorporates ollama as the LLM.

## Getting Started

### Prerequisites

- Go (version 1.15 or later recommended)
- Optionally, a .env file in the root directory to specify the PORT environment variable. If not set, the server defaults to port 8080.

### Installation

1. Clone the repository: 

```bash 
git clone https://github.com/ana-sofia-silva/chatgpt-clone-ollama
```

2. Navigate to the cloned directory.

3. Load environment variables (optional):
    - Create a .env file in the root directory.
    - Add the PORT variable

4. Run the server:

```bash 
go run main.go
```

### Usage

Once the server is running, it will listen on the specified port for incoming HTTP requests. You can interact with the server using any HTTP client by sending requests to localhost:<PORT> where <PORT> is the port number specified in your .env file or the default 8080.