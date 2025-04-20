# Config Package

The **config** package provides a robust mechanism to load, manage, and access your application's configuration settings. By utilizing this package, you can pull configuration data from either a local YAML file or from AWS AppConfig (when using the cloud fetcher) and keep it updated via periodic polling.

> **Important:** Always initialize the configuration using the `Init` function (with a valid polling duration of at least 10 seconds) **before** attempting to use any of the context, middleware, or helper functions. Failure to do so will break configuration access.

## What Does This Package Do?

- **Configuration Initialization:**  
  The package is initialized with the `Init(pollDuration time.Duration)` function. It reads the `CONFIG_SOURCE` environment variable to determine whether to load the configuration from a local YAML file or to fetch it via AWS AppConfig.  
  - For local configurations, set `CONFIG_SOURCE` to `local`.
  - For cloud configurations, set `CONFIG_SOURCE` to a string starting with `appconfig:` (e.g., `appconfig:myApplication`).

- **Configuration Polling and Updating:**  
  An observer is initialized to fetch the initial configuration and then continuously poll for changes on a configurable interval. The updated configuration is stored and made available safely for concurrent access.

- **Context Integration:**  
  The package provides functions such as `TODOContext` and `SetConfigInContext` to attach the configuration to a Go context. This allows you to share the configuration across your application.

- **Pre-built Getters:**  
  Helper functions—`Get`, `GetBool`, `GetString`, `GetInt`, `GetInt64`, `GetUint`, `GetUint64`, `GetFloat32`, `GetFloat64`, `GetTime`, `GetDuration`, and more—allow you to retrieve configuration values in a type-safe manner from the context.

- **Gin Middleware Support:**  
  The `Middleware` function adds the current configuration into every Gin HTTP request's context so that handlers can access the configuration directly.

## How to Initialize and Use the Config Package

### Setting Up the Environment

Before starting your application, set the `CONFIG_SOURCE` environment variable. For example:

```sh
# For local configuration
export CONFIG_SOURCE=local

# For AWS AppConfig-based cloud configuration:
export CONFIG_SOURCE=appconfig:myApplication
```

### Example 1: Standalone Application

The following example demonstrates how to initialize the configuration, attach it to a context, and then use the helper functions to retrieve values:

```go
package main

import (
	"fmt"
	"time"

	"github.com/mercor/payment-service/pkg/config"
)

func main() {
	// Initialize configuration.
	// Ensure that the CONFIG_SOURCE environment variable is set appropriately.
	err := config.Init(15 * time.Second)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize config: %v", err))
	}

	// Example: Use TODOContext to get a context with configuration attached.
	ctx, err := config.TODOContext()
	if err != nil {
		panic(fmt.Sprintf("Failed to attach config to context: %v", err))
	}

	// Retrieve a configuration value using helper functions.
	val := config.GetString(ctx, "exampleKey")
	fmt.Println("exampleKey:", val)
}
```

### Example 2: Gin-based Web Application

```go
package main

import (
	"github.com/gin-gonic/gin"
	"time"

	"github.com/mercor/payment-service/pkg/config"
)

func main() {
	// Initialize configuration.
	err := config.Init(15 * time.Second)
	if err != nil {
		panic(err)
	}

	// Create a new Gin router.
	router := gin.Default()

	// Use the configuration middleware to inject config into each request's context.
	router.Use(config.Middleware())

	// Define a route that retrieves and returns the current configuration.
	router.GET("/config", func(c *gin.Context) {
		// Retrieve the configuration from the Gin context.
		conf, exists := c.Get("config")
		if !exists {
			c.JSON(500, gin.H{"error": "config not found"})
			return
		}
		c.JSON(200, gin.H{"config": conf})
	})

	// Start the Gin server.
	router.Run() // Uses the default port (8080)
}
```

## Helper Functions

Once the configuration is attached to your context, you can easily retrieve values by key using the following functions:

- **Generic Getter:**
  - `Get(ctx context.Context, key string) interface{}`
  
- **Type-Specific Getters:**
  - `GetBool(ctx context.Context, key string) bool`
  - `GetString(ctx context.Context, key string) string`
  - `GetInt(ctx context.Context, key string) int`
  - `GetInt64(ctx context.Context, key string) int64`
  - `GetUint(ctx context.Context, key string) uint`
  - `GetUint64(ctx context.Context, key string) uint64`
  - `GetFloat32(ctx context.Context, key string) float32`
  - `GetFloat64(ctx context.Context, key string) float64`
  - `GetTime(ctx context.Context, key string) time.Time`
  - `GetDuration(ctx context.Context, key string) time.Duration`
  - `GetSlice(ctx context.Context, key string) []interface{}`
  - `GetBoolSlice(ctx context.Context, key string) []bool`
  - `GetIntSlice(ctx context.Context, key string) []int`
  - `GetStringSlice(ctx context.Context, key string) []string`

Additionally, the package includes environment-checking functions such as `IsDevelopment`, `IsStaging`, `IsProduction`, and `IsLocal` to help tailor behavior based on the current environment.

## Summary

The **config** package enables your application to:

- Dynamically load configuration from local YAML files or AWS AppConfig.
- Continuously poll for configuration updates and safely share the latest configuration.
- Attach the configuration to Go contexts for easy access across your application.
- Retrieve configuration values with type-safe helper functions.
- Seamlessly integrate configuration data into Gin-based web applications through middleware.

Remember: **Always call `Init` before accessing any configuration.** This ensures that your configuration is properly loaded and attached to the context, avoiding runtime errors.