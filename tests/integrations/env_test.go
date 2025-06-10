package integrations

import (
	"os"
	"testing"

	"github.com/galiherlangga/go-attendance/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	t.Run("Successful Load", func(t *testing.T) {
		envContent := "TEST_KEY=TEST_VALUE\n"
		envFile := ".env.test"
		err := os.WriteFile(envFile, []byte(envContent), 0644)
		assert.NoError(t, err)
		defer os.Remove(envFile)

		err = config.LoadEnv(envFile)
		assert.NoError(t, err)

		value, exists := os.LookupEnv("TEST_KEY")
		assert.True(t, exists)
		assert.Equal(t, "TEST_VALUE", value)

		os.Unsetenv("TEST_KEY")
		// Clean up the environment variable after the test
	})

	t.Run("File not found", func(t *testing.T) {
		err := config.LoadEnv("nonexistent.env")
		assert.Error(t, err)
	})

	t.Run("Invalid format", func(t *testing.T) {
		envContent := "TEST_KEY\n"
		envFile := ".env.test"
		err := os.WriteFile(envFile, []byte(envContent), 0644)
		assert.NoError(t, err)
		defer os.Remove(envFile)

		// Load the environment variables from the file
		err = config.LoadEnv(envFile)
		assert.Error(t, err)
	})
}

func TestGetEnv(t *testing.T) {
	// Set an environment variable
	os.Setenv("EXISTING_KEY", "EXISTING_VALUE")
	defer os.Unsetenv("EXISTING_KEY")

	// Test retrieving an existing key
	value := config.GetEnv("EXISTING_KEY", "DEFAULT_VALUE")
	assert.Equal(t, "EXISTING_VALUE", value)

	// Test retrieving a non-existing key with a default value
	value = config.GetEnv("NON_EXISTING_KEY", "DEFAULT_VALUE")
	assert.Equal(t, "DEFAULT_VALUE", value)
}
