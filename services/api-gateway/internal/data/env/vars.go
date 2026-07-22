package env

import (
	"fmt"
	"os"
)

const (
	FieldApiGatewayAddr  = "API_GATEWAY_ADDR"
	FieldAuthServiceAddr = "AUTH_SERVICE_ADDR"
)

type Vars map[string]string

type requiredField string

func ParseVars[F string | requiredField](fields ...F) (Vars, error) {
	vars := make(Vars, len(fields))

	for _, field := range fields {
		switch any(field).(type) {
		case string:
			s := any(field).(string)
			vars[s] = os.Getenv(s)

		case requiredField:
			s := string(any(field).(requiredField))
			v := os.Getenv(s)
			if v == "" {
				return nil, fmt.Errorf("env variable %q required", s)
			}
			vars[s] = v
		}
	}

	return vars, nil
}

func Required(field string) requiredField {
	return requiredField(field)
}

func (v Vars) Get(field string) string {
	return v[field]
}
