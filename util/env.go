// util/env.go
//
// Load environment variables from a local .env file (when not deployed).

package util

import (
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

// ============================================================================
// Exported constants
// ============================================================================

// Path relative to project root holding the ".env" variable assignments.
const LOCAL_ENV = "tmp/env"

// ============================================================================
// Internal variables
// ============================================================================

var envLoaded bool

// ============================================================================
// Exported functions
// ============================================================================

// Get the value of an environment variable, loading LOCAL_ENV if necessary.
func Getenv(name string) string {
    value, found := os.LookupEnv(name)
    if found || envLoaded {
        return value
    } else {
        LoadEnv()
        return os.Getenv(name)
    }
}

// Load environment variables defined in LOCAL_ENV.
//  NOTE: Desktop only.
func LoadEnv() {
    if !envLoaded {
        env := envFile()
        err := godotenv.Load(env)
        if (err != nil) && testing.Testing() {
            panic(err)
        }
        envLoaded = true
    }
}

// ============================================================================
// Internal functions
// ============================================================================

// Return the absolute pathname of the environment file.
func envFile() string {
    if env := LOCAL_ENV; strings.HasPrefix(env, "/") {
        return env
    } else {
        return RootPath() + "/" + env
    }
}
