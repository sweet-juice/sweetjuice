package core

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// InvokeCall parses incoming payloads from the WebView container and executes the matching Go routine.
func (a *Application) InvokeCall(methodKey string, rawArgs []json.RawMessage) (interface{}, error) {
	// Pull the generic wrapper and cast it back to the internal unexported boundMethod structure
	rawBound, exists := a.methods[methodKey]
	if !exists {
		return nil, fmt.Errorf("method identity '%s' not registered with application", methodKey)
	}

	bound := rawBound.(boundMethod)

	// Safety check to ensure argument counts align with Go function signature parameters
	if len(rawArgs) != len(bound.ParamTypes) {
		return nil, fmt.Errorf("argument length mismatch: expected %d, got %d", len(bound.ParamTypes), len(rawArgs))
	}

	invokingArgs := make([]reflect.Value, len(rawArgs))
	for i, argRaw := range rawArgs {
		targetType := bound.ParamTypes[i]

		// Create a pointer instance to unmarshal raw bytes into
		allocatedPtr := reflect.New(targetType)
		if err := json.Unmarshal(argRaw, allocatedPtr.Interface()); err != nil {
			return nil, fmt.Errorf("failed to parse parameter %d to type %s: %w", i, targetType, err)
		}
		invokingArgs[i] = allocatedPtr.Elem()
	}

	// Trigger execution via standard reflection mechanisms
	results := bound.MethodValue.Call(invokingArgs)
	if len(results) == 0 {
		return nil, nil
	}

	// Support typical core multi-value results returns: (data struct, error)
	if len(results) == 2 {
		errVal := results[1].Interface()
		if errVal != nil {
			if err, ok := errVal.(error); ok {
				return nil, err
			}
		}
		return results[0].Interface(), nil
	}

	return results[0].Interface(), nil
}
