package fetcher

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type KeyValue struct {
	prefix string
	value  map[string]interface{}
}

func ParseYAMLToConfigMap(yamlContent string) (map[string]interface{}, error) {
	var parsedYAML map[string]interface{}
	configMap := make(map[string]interface{})

	// Unmarshal YAML into a map
	err := yaml.Unmarshal([]byte(yamlContent), &parsedYAML)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Queue for BFS processing of the YAML map
	queue := []KeyValue{
		{
			prefix: "",
			value:  parsedYAML,
		},
	}

	// Process the queue
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for key, val := range current.value {
			// Construct the new key
			fullKey := current.prefix + "." + key
			if current.prefix == "" {
				fullKey = key
			}
			fullKey = strings.ToLower(fullKey)

			switch valTyped := val.(type) {
			case map[string]interface{}:
				// If the value is a nested map, add it to the queue for further processing
				queue = append(queue, KeyValue{
					prefix: fullKey,
					value:  valTyped,
				})
			case map[interface{}]interface{}:
				// Ensure all keys are strings
				for interfaceKey := range valTyped {
					if _, isString := interfaceKey.(string); !isString {
						return nil, fmt.Errorf("non-string key found: %T for key %v", interfaceKey, interfaceKey)
					}
				}
			default:
				// Handle non-map types, ensuring uniqueness and non-nil values
				if val == nil {
					return nil, fmt.Errorf("nil value encountered for key: %s", fullKey)
				}
				if _, exists := configMap[fullKey]; exists {
					return nil, fmt.Errorf("duplicate key found: %s (keys are case insensitive)", fullKey)
				}
				configMap[fullKey] = val
			}
		}
	}

	return configMap, nil
}
