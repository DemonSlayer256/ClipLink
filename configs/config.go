package configs

import (
	"bufio"
	"os"
	"strings"
)

func LoadEnv(keys ...string) string {
	env := make(map[string]string)

	file, err := os.Open(".env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		// Split line into key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) || (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
			value = value[1 : len(value)-1]
		}
		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	// Filter only the requested keys
	var required string
	for _, key := range keys {
		if value, exists := env[key]; exists {
			required = value
			break
		}
	}

	return required
}
