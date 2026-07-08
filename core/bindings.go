package core

import (
	"fmt"
	"reflect"
	"strings"
)

// boundMethod holds the unexported runtime reflection data for execution targeting.
// Keeping it lowercase hides it completely from the gobind parser.
type boundMethod struct {
	ParamTypes  []reflect.Type
	MethodValue reflect.Value
}

// parseBindings loops through all bound interfaces and structures to extract call metadata.
func (a *Application) parseBindings() error {
	for _, service := range a.options.Bind {
		val := reflect.ValueOf(service)
		typ := reflect.TypeOf(service)

		// Verification: Ensure the developer is passing pointers to structs for proper state allocation
		if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("bound target must be a pointer to a struct, got %s", val.Kind())
		}

		structName := typ.Elem().Name()
		fmt.Printf("[sweet-juice] Mapping service layer: %s\n", structName)

		// Iterate through all methods exposed on the struct pointer type
		for i := 0; i < val.NumMethod(); i++ {
			method := val.Method(i)
			methodName := typ.Method(i).Name

			// Enforce standard Go visibility rules: skip unexported lowercase methods
			if strings.ToUpper(methodName[0:1]) != methodName[0:1] {
				continue
			}

			methodType := method.Type()
			paramTypes := make([]reflect.Type, methodType.NumIn())
			for j := 0; j < methodType.NumIn(); j++ {
				paramTypes[j] = methodType.In(j)
			}

			// Generate the clean lookup key (e.g., "AcademicService.GetDashboardStats")
			bindingKey := fmt.Sprintf("%s.%s", structName, methodName)

			// Store data securely inside generic methods map using the unexported type
			a.methods[bindingKey] = boundMethod{
				ParamTypes:  paramTypes,
				MethodValue: method,
			}
			fmt.Printf("  -> Bound method identifier: %s\n", bindingKey)
		}
	}
	return nil
}
