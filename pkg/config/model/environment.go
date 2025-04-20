package model

type Environment string

const (
	Development Environment = "development"
	Dev         Environment = "dev" // Deprecated: Use Development instead.

	Staging    Environment = "staging"
	Production Environment = "production"
	Prod       Environment = "prod" // Deprecated: Use Production instead.

	Local   Environment = "local"
	Invalid Environment = "invalid"
)

// NewEnvironment creates a new Environment instance from the given string.
// If the provided string does not represent a valid environment, it returns Invalid.
// Parameters:
//   - env: A string representing the environment.
//
// Returns:
//   - Environment: A valid Environment instance or Invalid if the input is not valid.
func NewEnvironment(env string) Environment {
	environment := Environment(env)
	if !environment.IsValid() {
		return Invalid
	}
	return environment
}

// String converts the Environment type to its string representation.
// It implements the Stringer interface.
func (e Environment) String() string {
	return string(e)
}

// checks if the environment is valid
func (e Environment) IsValid() bool {
	switch e {
	case Development, Dev, Staging, Production, Prod, Local:
		return true
	default:
		return false
	}
}

// isProduction returns true if the environment is production
func (e Environment) IsProduction() bool {
	return e == Production || e == Prod
}

// isDevelopment returns true if the environment is development
func (e Environment) IsDevelopment() bool {
	return e == Development || e == Dev
}

// isStaging returns true if the environment is staging
func (e Environment) IsStaging() bool {
	return e == Staging
}

// isLocal returns true if the environment is local
func (e Environment) IsLocal() bool {
	return e == Local
}
