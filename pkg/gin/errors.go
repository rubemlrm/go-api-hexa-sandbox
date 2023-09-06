package gin

import "fmt"

type HttpConfigurationError struct {
	Input string
}

func (e *HttpConfigurationError) Error() string {
	return fmt.Sprintf("Error validating configuration: %s", e.Input)
}
