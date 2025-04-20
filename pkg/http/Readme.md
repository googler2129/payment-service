# HTTP Package

This package provides a robust and flexible HTTP client and server implementation built on top of popular Go frameworks. It offers a clean API for making HTTP requests and setting up HTTP servers with middleware support.

## Features

- HTTP Client with support for all standard HTTP methods (GET, POST, PUT, PATCH, DELETE)
- HTTP Server implementation using Gin framework
- Customizable timeouts and middleware support
- Built-in status code handling and type safety
- NewRelic integration
- Graceful shutdown support

## HTTP Client Usage

### Connection Pooling Configuration

The HTTP client uses connection pooling to improve performance. Here are the key configuration parameters and their implications:

```go
transport := &nethttp.Transport{
	MaxIdleConns:        100,  // Maximum number of idle connections across all hosts
	MaxIdleConnsPerHost: 100,  // Maximum number of idle connections per host
	IdleConnTimeout:     90 * time.Second,  // How long to keep idle connections alive
}
```

#### Connection Pool Parameters Explained

1. **MaxIdleConns**
   - Controls the maximum number of idle connections across all hosts
   - Default: 100
   - Advantages:
     - Reduces connection establishment overhead
     - Improves response times for subsequent requests
     - Reduces CPU and memory usage for connection creation
   - Drawbacks:
     - Each idle connection consumes memory
     - Too many idle connections can waste system resources
   - Best practices:
     - Set based on expected concurrent connections
     - Consider available system memory
     - Monitor connection usage patterns

2. **MaxIdleConnsPerHost**
   - Controls maximum idle connections per host
   - Default: 2 (in Go's default transport)
   - Advantages:
     - Prevents single host from consuming all connection slots
     - Helps distribute resources across different backends
     - Improves connection reuse for frequently accessed hosts
   - Drawbacks:
     - Too low can cause connection churn
     - Too high can waste resources on rarely accessed hosts
   - Best practices:
     - Set equal to or lower than MaxIdleConns
     - Consider traffic distribution across hosts
     - For single-host scenarios, can match MaxIdleConns

3. **Connection Pooling Scenarios**

```go
// High-traffic, multi-host scenario
transport := &nethttp.Transport{
	MaxIdleConns:        1000,
	MaxIdleConnsPerHost: 100,
	IdleConnTimeout:     90 * time.Second,
}

// Single-host, high-throughput scenario
transport := &nethttp.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
	IdleConnTimeout:     90 * time.Second,
}

// Low-traffic, multi-host scenario
transport := &nethttp.Transport{
	MaxIdleConns:        50,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     60 * time.Second,
}
```

### Creating a New Client

```go
import (
	"github.com/mercor/payment-service/pkg/http"
	nethttp "net/http"
)

// Create transport
transport := &nethttp.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 100,
}

// Initialize client with base URL
client, err := http.NewHTTPClient(
	"my-service",           // client service name
	"https://api.example.com", // base URL
	transport,
	http.WithTimeout(30 * time.Second), // optional timeout
)
if err != nil {
	// Handle error
}
```

### Making HTTP Requests

```go
// Create request
request := &http.Request{
	Url: "/users",
	Body: map[string]interface{}{
		"name": "John Doe",
		"email": "john@example.com",
	},
	Headers: map[string][]string{
		"Content-Type": {"application/json"},
	},
	QueryParams: url.Values{
		"page": []string{"1"},
		"limit": []string{"10"},
	},
	Timeout: 5 * time.Second, // Optional request-specific timeout
}

// For GET request
var response YourResponseStruct
resp, err := client.Get(request, &response)

// For POST request
resp, err := client.Post(request, &response)

// For PUT request
resp, err := client.Put(request, &response)

// For PATCH request
resp, err := client.Patch(request, &response)

// For DELETE request
resp, err := client.Delete(request, &response)
```

## HTTP Server Usage

### Creating and Starting a Server

```go
import (
	"github.com/mercor/payment-service/pkg/http"
	"time"
)

// Initialize server with custom timeouts
server := http.InitializeServer(
	":8080",                // listen address
	10 * time.Second,       // read timeout
	10 * time.Second,       // write timeout
	70 * time.Second,       // idle timeout
	// Add custom middleware here...
)

// Add routes
server.GET("/health", func(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
})

// Start the server
if err := server.StartServer("my-service"); err != nil {
	// Handle error
}
```

### Status Codes

The package provides type-safe status codes and helper methods:

```go
// Using status codes
if statusCode.Is2xx() {
	// Handle success
}

if statusCode.Is4xx() {
	// Handle client error
}

if statusCode.Is5xx() {
	// Handle server error
}

// Get status code text
text := http.StatusBadRequest.String() // Returns "Invalid Request"
```

## Features and Best Practices

1. **Automatic Request ID**: The server automatically adds request IDs if not present in headers
2. **Recovery Middleware**: Built-in panic recovery
3. **Graceful Shutdown**: Proper shutdown handling with timeout
4. **NewRelic Integration**: Built-in support for NewRelic monitoring
5. **Flexible Middleware**: Support for custom middleware in both client and server
6. **Timeout Management**: Configurable timeouts at both client and request level

## Error Handling

The package provides proper error handling and status code management. All operations return appropriate errors that should be handled by the caller.

```go
resp, err := client.Get(request, &response)
if err != nil {
	// Handle error
}

if !http.StatusCode(resp.StatusCode()).Is2xx() {
	// Handle non-2xx response
}
```

## Dependencies

- github.com/gin-gonic/gin
- github.com/go-resty/resty/v2
- github.com/newrelic/go-agent/v3/newrelic
- Standard Go packages (context, net/http, etc.)
