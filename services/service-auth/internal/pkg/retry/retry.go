package retry

import "fmt"

func Retry(n int, fn func() error) error {
	var err error
	for range n {
		if err = fn(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("reached maximum number of errors: %w", err)
}
