# Package env

## Overview
The env package manages application configuration through environment variables. It provides functions to load, parse, and validate configuration settings from the environment.

## Key Components
- Environment Parsers: Read and convert env variables.
- Configuration Loaders: Initialize application settings.
- Validation Utilities: Ensure required configurations are set.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/mercor/payment-service/pkg/env"
)

func main() {
	config := env.LoadConfig()
	fmt.Println("Loaded configuration:", config)
}
~~~

## Notes
- Essential for 12-factor app design.
