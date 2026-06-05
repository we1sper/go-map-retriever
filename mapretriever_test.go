package mapretriever

import (
	"strings"
	"testing"
)

// --- test data builders ---

func nestedMap() map[string]interface{} {
	return map[string]interface{}{
		"author": "welsper",
		"emails": []string{
			"welsper@qq.com",
			"welsper@nit.com",
			"welsper@nekomi.com",
		},
		"details": map[string]interface{}{
			"country": "China",
			"age":     18,
			"job": map[string]string{
				"company": "China Merchants Bank",
				"salary":  "100,000/m",
			},
		},
		"others": map[interface{}]string{
			1:      "test",
			"key1": "value1",
		},
		"active": true,
		"score":  95.5,
	}
}

func typedMap() map[string]interface{} {
	return map[string]interface{}{
		"intVal":      int(42),
		"int8Val":     int8(8),
		"int16Val":    int16(16),
		"int32Val":    int32(32),
		"int64Val":    int64(64),
		"uintVal":     uint(100),
		"uint8Val":    uint8(8),
		"uint16Val":   uint16(16),
		"uint32Val":   uint32(32),
		"uint64Val":   uint64(64),
		"float32Val":  float32(3.14),
		"float64Val":  float64(2.718),
		"boolVal":     true,
		"strVal":      "hello",
		"strSliceVal": []string{"a", "b", "c"},
		"sliceVal":    []interface{}{"x", "y", "z"},
	}
}

// =============================================================================
// NewMapRetriever
// =============================================================================

func TestNewMapRetriever(t *testing.T) {
	t.Run("with valid map", func(t *testing.T) {
		mr := NewMapRetriever(nestedMap())
		if mr == nil {
			t.Fatal("expected non-nil MapRetriever")
		}
		if !mr.Success() {
			t.Error("expected Success() to be true for valid map")
		}
		if mr.Value() == nil {
			t.Error("expected non-nil Value()")
		}
	})

	t.Run("with nil", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		if mr == nil {
			t.Fatal("expected non-nil MapRetriever")
		}
		if mr.Value() != nil {
			t.Errorf("expected nil Value(), got %v", mr.Value())
		}
	})
}

// =============================================================================
// Get
// =============================================================================

func TestGet(t *testing.T) {
	data := nestedMap()
	mr := NewMapRetriever(data)

	t.Run("single key success", func(t *testing.T) {
		result := mr.Get("author")
		if !result.Success() {
			t.Errorf("expected success, got error: %v", result.Error())
		}
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper" {
			t.Errorf("expected 'welsper', got %q", val)
		}
	})

	t.Run("nested keys", func(t *testing.T) {
		result := mr.Get("details", "country")
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "China" {
			t.Errorf("expected 'China', got %q", val)
		}
	})

	t.Run("deeply nested keys", func(t *testing.T) {
		result := mr.Get("details", "job", "company")
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "China Merchants Bank" {
			t.Errorf("expected 'China Merchants Bank', got %q", val)
		}
	})

	t.Run("int key in interface{} map", func(t *testing.T) {
		result := mr.Get("others", 1)
		if !result.Success() {
			t.Errorf("expected success, got error: %v", result.Error())
		}
		val, _ := result.String()
		if val != "test" {
			t.Errorf("expected 'test', got %q", val)
		}
	})

	t.Run("string key in interface{} map", func(t *testing.T) {
		result := mr.Get("others", "key1")
		if !result.Success() {
			t.Errorf("expected success, got error: %v", result.Error())
		}
		val, _ := result.String()
		if val != "value1" {
			t.Errorf("expected 'value1', got %q", val)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		result := mr.Get("nonexistent")
		if result.Success() {
			t.Error("expected failure for missing key")
		}
		if result.Error() == nil {
			t.Error("expected non-nil error for missing key")
		}
	})

	t.Run("wrong key type", func(t *testing.T) {
		result := mr.Get("details", "job", 123)
		if result.Success() {
			t.Error("expected failure for wrong key type")
		}
	})

	t.Run("get on non-map value", func(t *testing.T) {
		result := mr.Get("author", "subkey")
		if result.Success() {
			t.Error("expected failure when getting key on non-map value")
		}
	})

	t.Run("get on nil parent", func(t *testing.T) {
		result := NewMapRetriever(nil).Get("key")
		if result.Success() {
			t.Error("expected failure when getting key on nil parent")
		}
	})
}

// =============================================================================
// At
// =============================================================================

func TestAt(t *testing.T) {
	data := nestedMap()
	mr := NewMapRetriever(data)

	t.Run("valid index", func(t *testing.T) {
		result := mr.Get("emails").At(0)
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@qq.com" {
			t.Errorf("expected 'welsper@qq.com', got %q", val)
		}
	})

	t.Run("second valid index", func(t *testing.T) {
		result := mr.Get("emails").At(1)
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@nit.com" {
			t.Errorf("expected 'welsper@nit.com', got %q", val)
		}
	})

	t.Run("negative index wraps around", func(t *testing.T) {
		result := mr.Get("emails").At(-1)
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@nekomi.com" {
			t.Errorf("expected 'welsper@nekomi.com', got %q", val)
		}
	})

	t.Run("out of bounds positive", func(t *testing.T) {
		result := mr.Get("emails").At(10)
		if result.Success() {
			t.Error("expected failure for out of bounds index")
		}
	})

	t.Run("out of bounds negative", func(t *testing.T) {
		result := mr.Get("emails").At(-100)
		if result.Success() {
			t.Error("expected failure for out of bounds negative index")
		}
	})

	t.Run("multiple positions via At", func(t *testing.T) {
		result := mr.Get("emails").At(0, 0) // first char? No — emails[0] is a string, At on string fails
		// emails[0] is a string, not a slice — so At(0) on it should fail
		_ = result
	})

	t.Run("at on non-slice", func(t *testing.T) {
		result := mr.Get("author").At(0)
		if result.Success() {
			t.Error("expected failure for At on non-slice")
		}
	})

	t.Run("at on nil parent", func(t *testing.T) {
		result := NewMapRetriever(nil).At(0)
		if result.Success() {
			t.Error("expected failure for At on nil parent")
		}
	})
}

// =============================================================================
// Head / Tail
// =============================================================================

func TestHead(t *testing.T) {
	data := nestedMap()
	mr := NewMapRetriever(data)

	t.Run("head returns first element", func(t *testing.T) {
		result := mr.Get("emails").Head()
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@qq.com" {
			t.Errorf("expected 'welsper@qq.com', got %q", val)
		}
	})

	t.Run("head on empty slice", func(t *testing.T) {
		emptySlice := NewMapRetriever([]int{})
		result := emptySlice.Head()
		if result.Success() {
			t.Error("expected failure for Head on empty slice")
		}
	})
}

func TestTail(t *testing.T) {
	data := nestedMap()
	mr := NewMapRetriever(data)

	t.Run("tail returns last element", func(t *testing.T) {
		result := mr.Get("emails").Tail()
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@nekomi.com" {
			t.Errorf("expected 'welsper@nekomi.com', got %q", val)
		}
	})

	t.Run("tail on empty slice", func(t *testing.T) {
		emptySlice := NewMapRetriever([]int{})
		result := emptySlice.Tail()
		if result.Success() {
			t.Error("expected failure for Tail on empty slice")
		}
	})
}

// =============================================================================
// Fetch
// =============================================================================

func TestFetch(t *testing.T) {
	data := nestedMap()
	mr := NewMapRetriever(data)

	t.Run("mixed string and int path success", func(t *testing.T) {
		result := mr.Fetch("emails", 0)
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@qq.com" {
			t.Errorf("expected 'welsper@qq.com', got %q", val)
		}
	})

	t.Run("multiple mixed segments", func(t *testing.T) {
		result := mr.Fetch("details", "job", "company")
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "China Merchants Bank" {
			t.Errorf("expected 'China Merchants Bank', got %q", val)
		}
	})

	t.Run("out of bounds via Fetch", func(t *testing.T) {
		result := mr.Fetch("emails", 10)
		if result.Success() {
			t.Error("expected failure for out of bounds via Fetch")
		}
	})

	t.Run("int16 as index via Fetch", func(t *testing.T) {
		result := mr.Fetch("emails", int16(1))
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@nit.com" {
			t.Errorf("expected 'welsper@nit.com', got %q", val)
		}
	})

	t.Run("int8 as index via Fetch", func(t *testing.T) {
		result := mr.Fetch("emails", int8(2))
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@nekomi.com" {
			t.Errorf("expected 'welsper@nekomi.com', got %q", val)
		}
	})

	t.Run("int32 as index via Fetch", func(t *testing.T) {
		result := mr.Fetch("emails", int32(2))
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@nekomi.com" {
			t.Errorf("expected 'welsper@nekomi.com', got %q", val)
		}
	})

	t.Run("int64 as index via Fetch", func(t *testing.T) {
		result := mr.Fetch("emails", int64(0))
		val, err := result.String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "welsper@qq.com" {
			t.Errorf("expected 'welsper@qq.com', got %q", val)
		}
	})
}

// =============================================================================
// Type extractors — success cases
// =============================================================================

func TestTypeExtractors_Success(t *testing.T) {
	data := typedMap()
	mr := NewMapRetriever(data)

	t.Run("Bool", func(t *testing.T) {
		val, err := mr.Get("boolVal").Bool()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != true {
			t.Errorf("expected true, got %v", val)
		}
	})

	t.Run("Int", func(t *testing.T) {
		val, err := mr.Get("intVal").Int()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 42 {
			t.Errorf("expected 42, got %v", val)
		}
	})

	t.Run("Int8", func(t *testing.T) {
		val, err := mr.Get("int8Val").Int8()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 8 {
			t.Errorf("expected 8, got %v", val)
		}
	})

	t.Run("Int16", func(t *testing.T) {
		val, err := mr.Get("int16Val").Int16()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 16 {
			t.Errorf("expected 16, got %v", val)
		}
	})

	t.Run("Int32", func(t *testing.T) {
		val, err := mr.Get("int32Val").Int32()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 32 {
			t.Errorf("expected 32, got %v", val)
		}
	})

	t.Run("Int64", func(t *testing.T) {
		val, err := mr.Get("int64Val").Int64()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 64 {
			t.Errorf("expected 64, got %v", val)
		}
	})

	t.Run("Uint", func(t *testing.T) {
		val, err := mr.Get("uintVal").Uint()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 100 {
			t.Errorf("expected 100, got %v", val)
		}
	})

	t.Run("Uint8", func(t *testing.T) {
		val, err := mr.Get("uint8Val").Uint8()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 8 {
			t.Errorf("expected 8, got %v", val)
		}
	})

	t.Run("Uint16", func(t *testing.T) {
		val, err := mr.Get("uint16Val").Uint16()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 16 {
			t.Errorf("expected 16, got %v", val)
		}
	})

	t.Run("Uint32", func(t *testing.T) {
		val, err := mr.Get("uint32Val").Uint32()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 32 {
			t.Errorf("expected 32, got %v", val)
		}
	})

	t.Run("Uint64", func(t *testing.T) {
		val, err := mr.Get("uint64Val").Uint64()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 64 {
			t.Errorf("expected 64, got %v", val)
		}
	})

	t.Run("Float32", func(t *testing.T) {
		val, err := mr.Get("float32Val").Float32()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != float32(3.14) {
			t.Errorf("expected 3.14, got %v", val)
		}
	})

	t.Run("Float64", func(t *testing.T) {
		val, err := mr.Get("float64Val").Float64()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != 2.718 {
			t.Errorf("expected 2.718, got %v", val)
		}
	})

	t.Run("String", func(t *testing.T) {
		val, err := mr.Get("strVal").String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "hello" {
			t.Errorf("expected 'hello', got %q", val)
		}
	})

	t.Run("StringSlice", func(t *testing.T) {
		val, err := mr.Get("strSliceVal").StringSlice()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(val) != 3 {
			t.Fatalf("expected length 3, got %d", len(val))
		}
		if val[0] != "a" || val[1] != "b" || val[2] != "c" {
			t.Errorf("expected [a b c], got %v", val)
		}
	})
}

// =============================================================================
// Type extractors — error cases (type mismatch)
// =============================================================================

func TestTypeExtractors_TypeMismatch(t *testing.T) {
	data := typedMap()
	mr := NewMapRetriever(data)

	tests := []struct {
		name    string
		path    []any
		extract func(*MapRetriever) error
	}{
		{"Bool from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Bool(); return err }},
		{"Int from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Int(); return err }},
		{"Int8 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Int8(); return err }},
		{"Int16 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Int16(); return err }},
		{"Int32 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Int32(); return err }},
		{"Int64 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Int64(); return err }},
		{"Uint from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Uint(); return err }},
		{"Uint8 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Uint8(); return err }},
		{"Uint16 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Uint16(); return err }},
		{"Uint32 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Uint32(); return err }},
		{"Uint64 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Uint64(); return err }},
		{"Float32 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Float32(); return err }},
		{"Float64 from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.Float64(); return err }},
		{"String from int", []any{"intVal"}, func(m *MapRetriever) error { _, err := m.String(); return err }},
		{"StringSlice from string", []any{"strVal"}, func(m *MapRetriever) error { _, err := m.StringSlice(); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.extract(mr.Get(tt.path...))
			if err == nil {
				t.Error("expected error for type mismatch, got nil")
			}
		})
	}
}

// =============================================================================
// ValueSlice
// =============================================================================

func TestValueSlice(t *testing.T) {
	t.Run("slice of interface{}", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		result := mr.Get("sliceVal")
		children, err := result.ValueSlice()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(children) != 3 {
			t.Fatalf("expected 3 children, got %d", len(children))
		}
		for i, child := range children {
			if !child.Success() {
				t.Errorf("child %d: expected success", i)
			}
		}
	})

	t.Run("string slice via ValueSlice", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		result := mr.Get("emails")
		children, err := result.ValueSlice()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(children) != 3 {
			t.Fatalf("expected 3 children, got %d", len(children))
		}
		val, _ := children[0].String()
		if val != "welsper@qq.com" {
			t.Errorf("expected 'welsper@qq.com', got %q", val)
		}
	})

	t.Run("on nil value", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		_, err := mr.ValueSlice()
		if err == nil {
			t.Error("expected error for ValueSlice on nil")
		}
	})

	t.Run("on non-slice value", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		_, err := mr.Get("strVal").ValueSlice()
		if err == nil {
			t.Error("expected error for ValueSlice on non-slice value")
		}
	})
}

// =============================================================================
// Unsafe — panics on nil, returns zero value on error
// =============================================================================

func TestUnsafe(t *testing.T) {
	t.Run("String success", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("author").Unsafe().String()
		if val != "welsper" {
			t.Errorf("expected 'welsper', got %q", val)
		}
	})

	t.Run("String on error returns zero value", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		val := mr.Get("missing").Unsafe().String()
		if val != "" {
			t.Errorf("expected empty string, got %q", val)
		}
	})

	t.Run("Bool on error returns false", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		val := mr.Get("missing").Unsafe().Bool()
		if val != false {
			t.Errorf("expected false, got %v", val)
		}
	})

	t.Run("Int on error returns 0", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		val := mr.Get("missing").Unsafe().Int()
		if val != 0 {
			t.Errorf("expected 0, got %v", val)
		}
	})

	t.Run("Int8 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("int8Val").Unsafe().Int8()
		if val != 8 {
			t.Errorf("expected 8, got %v", val)
		}
	})

	t.Run("Int16 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("int16Val").Unsafe().Int16()
		if val != 16 {
			t.Errorf("expected 16, got %v", val)
		}
	})

	t.Run("Int32 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("int32Val").Unsafe().Int32()
		if val != 32 {
			t.Errorf("expected 32, got %v", val)
		}
	})

	t.Run("Int64 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("int64Val").Unsafe().Int64()
		if val != 64 {
			t.Errorf("expected 64, got %v", val)
		}
	})

	t.Run("Uint success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("uintVal").Unsafe().Uint()
		if val != 100 {
			t.Errorf("expected 100, got %v", val)
		}
	})

	t.Run("Uint8 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("uint8Val").Unsafe().Uint8()
		if val != 8 {
			t.Errorf("expected 8, got %v", val)
		}
	})

	t.Run("Uint16 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("uint16Val").Unsafe().Uint16()
		if val != 16 {
			t.Errorf("expected 16, got %v", val)
		}
	})

	t.Run("Uint32 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("uint32Val").Unsafe().Uint32()
		if val != 32 {
			t.Errorf("expected 32, got %v", val)
		}
	})

	t.Run("Uint64 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("uint64Val").Unsafe().Uint64()
		if val != 64 {
			t.Errorf("expected 64, got %v", val)
		}
	})

	t.Run("Float32 success", func(t *testing.T) {
		data := typedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("float32Val").Unsafe().Float32()
		if val != float32(3.14) {
			t.Errorf("expected 3.14, got %v", val)
		}
	})

	t.Run("Float64 on error returns 0", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		val := mr.Get("missing").Unsafe().Float64()
		if val != 0 {
			t.Errorf("expected 0, got %v", val)
		}
	})

	t.Run("StringSlice on error returns nil", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		val := mr.Get("missing").Unsafe().StringSlice()
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("ValueSlice on error returns nil", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		val := mr.Get("missing").Unsafe().ValueSlice()
		if val != nil {
			t.Errorf("expected nil, got %v", val)
		}
	})

	t.Run("chained unsafe access", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("details", "job", "company").Unsafe().String()
		if val != "China Merchants Bank" {
			t.Errorf("expected 'China Merchants Bank', got %q", val)
		}
	})
}

// =============================================================================
// Success / Error / Value / Unsafe
// =============================================================================

func TestSuccess(t *testing.T) {
	t.Run("true for valid path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		if !mr.Get("author").Success() {
			t.Error("expected Success() to be true")
		}
	})

	t.Run("true for root with no error", func(t *testing.T) {
		mr := NewMapRetriever(nestedMap())
		if !mr.Success() {
			t.Error("expected root Success() to be true")
		}
	})

	t.Run("false after error", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		result := mr.Get("missing_key")
		if result.Success() {
			t.Error("expected Success() to be false after missing key")
		}
	})
}

func TestError(t *testing.T) {
	t.Run("nil for successful path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		if err := mr.Get("author").Error(); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("returns root cause error", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		result := mr.Get("details", "job", "missing", "deep")
		err := result.Error()
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		// The root cause should be about "missing" key, not about "deep"
		if !strings.Contains(err.Error(), "missing") {
			t.Errorf("expected error to mention 'missing', got: %v", err)
		}
	})
}

func TestValue(t *testing.T) {
	t.Run("returns raw value", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("author").Value()
		if val != "welsper" {
			t.Errorf("expected 'welsper', got %v", val)
		}
	})

	t.Run("returns nil for error path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		val := mr.Get("missing").Value()
		if val != nil {
			t.Errorf("expected nil for missing key, got %v", val)
		}
	})
}

// =============================================================================
// Trace
// =============================================================================

func TestTrace(t *testing.T) {
	t.Run("simple path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		trace := mr.Get("author").Trace()
		expected := "source.get(author)"
		if trace != expected {
			t.Errorf("expected %q, got %q", expected, trace)
		}
	})

	t.Run("nested path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		trace := mr.Get("details", "country").Trace()
		expected := "source.get(details).get(country)"
		if trace != expected {
			t.Errorf("expected %q, got %q", expected, trace)
		}
	})

	t.Run("at path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		trace := mr.Get("emails").At(0).Trace()
		expected := "source.get(emails).at(0)"
		if trace != expected {
			t.Errorf("expected %q, got %q", expected, trace)
		}
	})

	t.Run("fetch path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		trace := mr.Fetch("emails", 0).Trace()
		expected := "source.get(emails).at(0)"
		if trace != expected {
			t.Errorf("expected %q, got %q", expected, trace)
		}
	})

	t.Run("head path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		trace := mr.Get("emails").Head().Trace()
		expected := "source.get(emails).at(0)"
		if trace != expected {
			t.Errorf("expected %q, got %q", expected, trace)
		}
	})

	t.Run("tail path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		trace := mr.Get("emails").Tail().Trace()
		expected := "source.get(emails).at(-1)"
		if trace != expected {
			t.Errorf("expected %q, got %q", expected, trace)
		}
	})
}

// =============================================================================
// Debug
// =============================================================================

func TestDebug(t *testing.T) {
	t.Run("success path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		debug := mr.Get("author").Debug()
		if !strings.Contains(debug, ": success") {
			t.Errorf("expected debug to contain ': success', got: %s", debug)
		}
		if !strings.Contains(debug, "get(author)") {
			t.Errorf("expected debug to contain path, got: %s", debug)
		}
	})

	t.Run("error path", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		debug := mr.Get("missing").Debug()
		if strings.Contains(debug, ": success") {
			t.Error("expected debug NOT to contain ': success' for error path")
		}
		// Should contain the error message
		if !strings.Contains(debug, "missing") {
			t.Errorf("expected debug to mention 'missing', got: %s", debug)
		}
	})

	t.Run("nested error shows root cause", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		debug := mr.Get("details", "job", "missing", "deep").Debug()
		if strings.Contains(debug, ": success") {
			t.Error("expected debug NOT to contain ': success' for nested error")
		}
		// Should point to where error occurred (missing key)
		if !strings.Contains(debug, "missing") {
			t.Errorf("expected debug to mention 'missing', got: %s", debug)
		}
	})
}

// =============================================================================
// Edge cases & chained operations
// =============================================================================

func TestEdgeCases(t *testing.T) {
	t.Run("nil root map", func(t *testing.T) {
		mr := NewMapRetriever(nil)
		if mr.Value() != nil {
			t.Error("expected nil value for nil root")
		}
		if !mr.Success() {
			t.Error("expected root Success() to be true even with nil raw")
		}
	})

	t.Run("empty map", func(t *testing.T) {
		mr := NewMapRetriever(map[string]interface{}{})
		result := mr.Get("anything")
		if result.Success() {
			t.Error("expected failure for key in empty map")
		}
	})

	t.Run("chained error propagation", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		// Get a key that doesn't exist, then try to go deeper
		result := mr.Get("details", "nonexistent", "deeper")
		if result.Success() {
			t.Error("expected failure for chained error")
		}
		// Error should be about "nonexistent", not "deeper"
		err := result.Error()
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if !strings.Contains(err.Error(), "nonexistent") {
			t.Errorf("expected error to mention 'nonexistent', got: %v", err)
		}
	})

	t.Run("Error walks to root cause", func(t *testing.T) {
		data := nestedMap()
		mr := NewMapRetriever(data)
		// First error is at "missing", then we try to go deeper (which creates child with nil parent error)
		result := mr.Get("details", "job", "missing", "deep")
		err := result.Error()
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		// The root cause should be the first error encountered
		if !strings.Contains(err.Error(), "missing") {
			t.Errorf("expected root error about 'missing', got: %v", err)
		}
	})

	t.Run("array as root", func(t *testing.T) {
		arr := []interface{}{"first", "second", "third"}
		mr := NewMapRetriever(arr)
		if !mr.Success() {
			t.Error("expected success for array root")
		}
		val, err := mr.At(0).String()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "first" {
			t.Errorf("expected 'first', got %q", val)
		}
	})
}
