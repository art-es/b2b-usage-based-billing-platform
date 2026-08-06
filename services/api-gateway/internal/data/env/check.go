package env

import (
	"fmt"
	"os"
)

func CheckEmpty(vars ...string) error {
	for _, name := range vars {
		if os.Getenv(name) == "" {
			return fmt.Errorf("env var %q is empty", name)
		}
	}
	return nil
}
