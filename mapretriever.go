package mapretriever

import (
	"fmt"
	"reflect"
	"strings"
)

type UnsafeMapRetriever struct {
	*MapRetriever
}

func (u *UnsafeMapRetriever) Bool() bool {
	value, _ := u.MapRetriever.Bool()
	return value
}

func (u *UnsafeMapRetriever) Int() int {
	value, _ := u.MapRetriever.Int()
	return value
}

func (u *UnsafeMapRetriever) Int8() int8 {
	value, _ := u.MapRetriever.Int8()
	return value
}

func (u *UnsafeMapRetriever) Int16() int16 {
	value, _ := u.MapRetriever.Int16()
	return value
}

func (u *UnsafeMapRetriever) Int32() int32 {
	value, _ := u.MapRetriever.Int32()
	return value
}

func (u *UnsafeMapRetriever) Int64() int64 {
	value, _ := u.MapRetriever.Int64()
	return value
}

func (u *UnsafeMapRetriever) Uint() uint {
	value, _ := u.MapRetriever.Uint()
	return value
}

func (u *UnsafeMapRetriever) Uint8() uint8 {
	value, _ := u.MapRetriever.Uint8()
	return value
}

func (u *UnsafeMapRetriever) Uint16() uint16 {
	value, _ := u.MapRetriever.Uint16()
	return value
}

func (u *UnsafeMapRetriever) Uint32() uint32 {
	value, _ := u.MapRetriever.Uint32()
	return value
}

func (u *UnsafeMapRetriever) Uint64() uint64 {
	value, _ := u.MapRetriever.Uint64()
	return value
}

func (u *UnsafeMapRetriever) Float32() float32 {
	value, _ := u.MapRetriever.Float32()
	return value
}

func (u *UnsafeMapRetriever) Float64() float64 {
	value, _ := u.MapRetriever.Float64()
	return value
}

func (u *UnsafeMapRetriever) String() string {
	value, _ := u.MapRetriever.String()
	return value
}

func (u *UnsafeMapRetriever) StringSlice() []string {
	value, _ := u.MapRetriever.StringSlice()
	return value
}

func (u *UnsafeMapRetriever) ValueSlice() []*MapRetriever {
	value, _ := u.MapRetriever.ValueSlice()
	return value
}

type MapRetriever struct {
	raw    any
	from   string
	err    error
	parent *MapRetriever
	unsafe *UnsafeMapRetriever
}

func NewMapRetriever(raw any) *MapRetriever {
	mapRetriever := &MapRetriever{
		raw:  raw,
		from: "",
	}
	mapRetriever.unsafe = &UnsafeMapRetriever{mapRetriever}
	return mapRetriever
}

func (m *MapRetriever) newChild(from string) *MapRetriever {
	result := NewMapRetriever(nil)
	result.from = from
	result.parent = m
	return result
}

func (m *MapRetriever) At(poses ...int) *MapRetriever {
	next := m

	for _, pos := range poses {
		next = next.at(pos)
	}

	return next
}

func (m *MapRetriever) Get(keys ...any) *MapRetriever {
	next := m

	for _, key := range keys {
		next = next.get(key)
	}

	return next
}

func (m *MapRetriever) Head() *MapRetriever {
	return m.at(0)
}

func (m *MapRetriever) Tail() *MapRetriever {
	return m.at(-1)
}

func (m *MapRetriever) Bool() (bool, error) {
	if value, ok := m.raw.(bool); ok {
		return value, nil
	}
	return false, fmt.Errorf("%v is not a bool", m.raw)
}

func (m *MapRetriever) Int() (int, error) {
	if value, ok := m.raw.(int); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a int", m.raw)
}

func (m *MapRetriever) Int8() (int8, error) {
	if value, ok := m.raw.(int8); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a int8", m.raw)
}

func (m *MapRetriever) Int16() (int16, error) {
	if value, ok := m.raw.(int16); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a int16", m.raw)
}

func (m *MapRetriever) Int32() (int32, error) {
	if value, ok := m.raw.(int32); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a int32", m.raw)
}

func (m *MapRetriever) Int64() (int64, error) {
	if value, ok := m.raw.(int64); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a int64", m.raw)
}

func (m *MapRetriever) Uint() (uint, error) {
	if value, ok := m.raw.(uint); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a uint", m.raw)
}

func (m *MapRetriever) Uint8() (uint8, error) {
	if value, ok := m.raw.(uint8); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a uint8", m.raw)
}

func (m *MapRetriever) Uint16() (uint16, error) {
	if value, ok := m.raw.(uint16); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a uint16", m.raw)
}

func (m *MapRetriever) Uint32() (uint32, error) {
	if value, ok := m.raw.(uint32); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a uint32", m.raw)
}

func (m *MapRetriever) Uint64() (uint64, error) {
	if value, ok := m.raw.(uint64); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a uint64", m.raw)
}

func (m *MapRetriever) Float32() (float32, error) {
	if value, ok := m.raw.(float32); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a float32", m.raw)
}

func (m *MapRetriever) Float64() (float64, error) {
	if value, ok := m.raw.(float64); ok {
		return value, nil
	}
	return 0, fmt.Errorf("%v is not a float64", m.raw)
}

func (m *MapRetriever) String() (string, error) {
	if value, ok := m.raw.(string); ok {
		return value, nil
	}
	return "", fmt.Errorf("%v is not a string", m.raw)
}

func (m *MapRetriever) StringSlice() ([]string, error) {
	if value, ok := m.raw.([]string); ok {
		return value, nil
	}
	return nil, fmt.Errorf("%v is not a []string", m.raw)
}

func (m *MapRetriever) ValueSlice() ([]*MapRetriever, error) {
	if m.raw == nil {
		return nil, fmt.Errorf("cast to value slice error: value is nil")
	}

	raw := reflect.ValueOf(m.raw)
	rawKind := raw.Type().Kind()

	if rawKind != reflect.Slice && rawKind != reflect.Array {
		return nil, fmt.Errorf("%v is not a slice or array", m.raw)
	}

	children := make([]*MapRetriever, raw.Len())

	for i := 0; i < raw.Len(); i++ {
		child := m.newChild(fmt.Sprintf("at(%d)", i))
		child.raw = raw.Index(i).Interface()
		children[i] = child
	}

	return children, nil
}

func (m *MapRetriever) Unsafe() *UnsafeMapRetriever {
	return m.unsafe
}

func (m *MapRetriever) Success() bool {
	return m.err == nil
}

func (m *MapRetriever) Error() error {
	return m.err
}

func (m *MapRetriever) Value() any {
	return m.raw
}

func (m *MapRetriever) Trace() string {
	trace := m.from

	for parent := m.parent; parent != nil; parent = parent.parent {
		trace = parent.from + "." + trace
	}

	return "source" + trace
}

func (m *MapRetriever) Debug() string {
	pos := 0
	trace := ""
	rootCause := ""

	for current := m; current != nil; current = current.parent {
		trace = current.from + "." + trace
		if current.err != nil && (current.parent == nil || current.parent.err == nil) {
			rootCause = current.err.Error()
			pos = len(trace) - 1
		}
	}

	trace = "source" + trace

	if pos == 0 {
		return trace + ": success"
	}

	padding := strings.Repeat(" ", len(trace)-pos)

	var debug strings.Builder

	debug.WriteString(trace)
	debug.WriteByte('\n')
	debug.WriteString(padding)
	debug.WriteString("x\n")
	debug.WriteString(padding)
	debug.WriteString("|\n")
	debug.WriteString(padding)
	debug.WriteString(rootCause)

	return debug.String()
}

func (m *MapRetriever) at(index int) *MapRetriever {
	child := m.newChild(fmt.Sprintf("at(%d)", index))

	if m.raw == nil {
		child.err = fmt.Errorf("cannot get value at %d: parent is nil", index)
		return child
	}

	raw := reflect.ValueOf(m.raw)
	rawKind := raw.Type().Kind()

	if rawKind != reflect.Slice && rawKind != reflect.Array {
		child.err = fmt.Errorf("cannot get value at %d: parent is not a slice or array", index)
		return child
	}

	length := raw.Len()

	if index < 0 {
		index += length
	}

	if index >= length {
		child.err = fmt.Errorf("%d out of bounds: parent length is %d", index, length)
	} else {
		child.raw = raw.Index(index).Interface()
	}

	return child
}

func (m *MapRetriever) get(key any) *MapRetriever {
	child := m.newChild(fmt.Sprintf("get(%v)", key))

	if m.raw == nil {
		child.err = fmt.Errorf("cannot get value for key %v: parent is nil", key)
		return child
	}

	raw := reflect.ValueOf(m.raw)
	rawType := raw.Type()

	if rawType.Kind() != reflect.Map {
		child.err = fmt.Errorf("cannot get value for key %v: parent is not a map", key)
		return child
	}

	theKey := reflect.ValueOf(key)
	theKeyKind := theKey.Kind()
	rawKeyKind := rawType.Key().Kind()

	if rawKeyKind != reflect.Interface && rawKeyKind != theKeyKind {
		child.err = fmt.Errorf("cannot get value for key %v: key type %v does not match map key type %v", key, theKeyKind, rawKeyKind)
		return child
	}

	if candidate := raw.MapIndex(theKey); candidate.IsValid() {
		child.raw = candidate.Interface()
	} else {
		child.err = fmt.Errorf("cannot get value for key %v: key not found in map", key)
	}

	return child
}
