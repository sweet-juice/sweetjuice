package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestService struct{}

func (s *TestService) Greet(name string) (string, error) {
	return fmt.Sprintf("Hello %s", name), nil
}

func (s *TestService) Multiply(a, b int) int {
	return a * b
}

func TestApplicationBindingsAndNativeMethods(t *testing.T) {
	app := NewApplication(Options{
		Name: "test-app",
		Bind: []interface{}{&TestService{}},
	})

	defer SetGlobalApp(nil)

	if err := app.Run(); err != nil {
		t.Fatalf("unexpected Run error: %v", err)
	}

	result, err := app.InvokeCall("TestService.Greet", []json.RawMessage{[]byte(`"Alice"`)})
	if err != nil {
		t.Fatalf("InvokeCall failed: %v", err)
	}

	if result != "Hello Alice" {
		t.Fatalf("expected Hello Alice, got %v", result)
	}

	result, err = app.InvokeCall("TestService.Multiply", []json.RawMessage{[]byte(`3`), []byte(`4`)})
	if err != nil {
		t.Fatalf("InvokeCall failed: %v", err)
	}

	if result != int(12) {
		t.Fatalf("expected 12, got %v", result)
	}

	_, err = app.InvokeCall("TestService.Unknown", nil)
	if err == nil {
		t.Fatal("expected error for unknown method, got nil")
	}

	app.RegisterNativeMethod("native:echo", func(args []json.RawMessage) (interface{}, error) {
		var value string
		if err := json.Unmarshal(args[0], &value); err != nil {
			return nil, err
		}
		return value, nil
	})

	nativeResult, err := app.InvokeNativeCall("native:echo", []json.RawMessage{[]byte(`"pong"`)})
	if err != nil {
		t.Fatalf("InvokeNativeCall failed: %v", err)
	}

	if nativeResult != "pong" {
		t.Fatalf("expected pong, got %v", nativeResult)
	}
}

func BenchmarkApplicationInvokeCall(b *testing.B) {
	app := NewApplication(Options{
		Name: "bench-app",
		Bind: []interface{}{&TestService{}},
	})
	if err := app.Run(); err != nil {
		b.Fatalf("Run failed: %v", err)
	}

	raw := []json.RawMessage{[]byte(`"Benchmark"`)}
	for i := 0; i < b.N; i++ {
		_, err := app.InvokeCall("TestService.Greet", raw)
		if err != nil {
			b.Fatalf("InvokeCall failed: %v", err)
		}
	}
}
