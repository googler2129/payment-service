# Package shutdown

## Overview
The shutdown package provides a robust mechanism for implementing graceful shutdown in Go applications. It handles system signals (SIGTERM, SIGINT) and ensures that all registered cleanup operations are executed properly before the application exits. The package implements a timeout-based shutdown mechanism to prevent indefinite hanging during cleanup.

## Features
- **Signal Handling**: Automatically catches SIGTERM and SIGINT signals
- **Two-Phase Shutdown**: Supports both drain and shutdown callbacks
- **Timeout Protection**: Enforces a maximum termination time (25 seconds)
- **Named Callbacks**: Each callback can be registered with a descriptive name for better logging
- **Global Server**: Singleton pattern for application-wide shutdown management
- **Wait Channel**: Provides a channel to block until shutdown is complete

## Installation
```bash
go get github.com/mercor/payment-service/pkg
```

## Core Components

### Callback Interface
```go
type Callback interface {
    Close() error
}
```
Any struct implementing this interface can be registered as a shutdown or drain callback.

### Main Functions
- `RegisterShutdownCallback(name string, callback Callback)`: Register a callback for the shutdown phase
- `RegisterDrainCallback(name string, callback Callback)`: Register a callback for the drain phase
- `GetWaitChannel() <-chan bool`: Get a channel that closes when shutdown is complete

## Usage Examples

### Basic Usage
```go
package main

import (
    "fmt"
    "github.com/mercor/payment-service/pkg/shutdown"
)

type CleanupHandler struct{}

func (c *CleanupHandler) Close() error {
    fmt.Println("Performing cleanup...")
    return nil
}

func main() {
    // Register shutdown callback
    shutdown.RegisterShutdownCallback("cleanup", &CleanupHandler{})

    fmt.Println("Application is running...")
    // Wait for shutdown signal
    <-shutdown.GetWaitChannel()
}
```

### Two-Phase Shutdown Example
```go
package main

import (
    "fmt"
    "github.com/mercor/payment-service/pkg/shutdown"
)

type DatabaseConnection struct{}
type HTTPServer struct{}

func (db *DatabaseConnection) Close() error {
    fmt.Println("Closing database connections...")
    return nil
}

func (srv *HTTPServer) Close() error {
    fmt.Println("Stopping HTTP server...")
    return nil
}

func main() {
    // Register drain callback (executed first)
    shutdown.RegisterDrainCallback("http-server", &HTTPServer{})
    
    // Register shutdown callback (executed after drain)
    shutdown.RegisterShutdownCallback("database", &DatabaseConnection{})

    fmt.Println("Application is running...")
    <-shutdown.GetWaitChannel()
}
```

## Shutdown Process
1. When a termination signal (SIGTERM/SIGINT) is received:
   - The package logs the received signal
   - Initiates the graceful shutdown process

2. Drain Phase:
   - All drain callbacks are executed in the order they were registered
   - Useful for stopping incoming traffic/requests

3. Shutdown Phase:
   - All shutdown callbacks are executed in the order they were registered
   - Used for cleaning up resources and connections

4. Timeout Protection:
   - The entire shutdown process must complete within 25 seconds
   - If exceeded, the application exits with a non-zero status code

## Best Practices
1. Register drain callbacks for services that need to stop accepting new work
2. Register shutdown callbacks for cleanup operations
3. Implement proper error handling in Close() methods
4. Use meaningful names when registering callbacks for better logging
5. Keep cleanup operations efficient to complete within the timeout period

## Error Handling
- Each callback's Close() method can return an error
- Errors are logged but don't stop the shutdown process
- If the shutdown process exceeds the timeout, the application exits with status code 1

## Thread Safety
The package is designed to be thread-safe and can handle concurrent registrations of callbacks.

## Limitations
- Maximum termination time is fixed at 25 seconds
- Callbacks are executed sequentially, not in parallel
- Once shutdown begins, new callbacks cannot be effectively registered

## Contributing
Contributions are welcome! Please ensure that any pull requests include appropriate tests and documentation updates.
