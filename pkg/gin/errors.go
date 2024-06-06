package gin

import "fmt"

type HTTPConfigurationError struct {
	Input string
}

func (e *HTTPConfigurationError) Error() string {
	return fmt.Sprintf("error validating configuration: %s", e.Input)
}
