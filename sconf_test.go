package sconf

import (
	"errors"
	"testing"
)

var dummyError = errors.New("dummy")

func TestSetReadLimit(t *testing.T) {

	type setReadLimitTest struct {
		limit int64
		err   error
	}

	var setReadLimitTests = []setReadLimitTest{
		{-1, dummyError},
		{0, dummyError},
		{20 * 1024 * 1024, dummyError},
		{1024, nil},
	}

	for _, test := range setReadLimitTests {
		err := SetReadLimit(test.limit)
		if test.err != nil && err == nil {
			t.Fatalf("SetReadLimit(%d): Failure expected, but got OK", test.limit)
		}
		if test.err == nil && err != nil {
			t.Fatalf("SetReadLimit(%d): OK expected, but got failure", test.limit)
		}
	}

}

func TestRead(t *testing.T) {

	type readTest struct {
		path string
		err  error
	}

	// test strict mode first
	ToggleStrictMode()
	readTests := []readTest{
		{"tests/strict-scopes.config", dummyError},
		{"tests/strict-keys.config", dummyError},
	}

	for _, test := range readTests {
		err := Read(test.path)
		if test.err != nil && err == nil {
			t.Fatalf("Read(%s): Failure expected, but got OK", test.path)
		}
		if test.err == nil && err != nil {
			t.Fatalf("Read(%s): OK expected, but got failure (%s)", test.path, err)
		}
	}

	// other tests
	Clear()
	ToggleStrictMode()
	readTests = []readTest{
		{"", dummyError},
		{"tests/no such file", dummyError},
		{"tests/big.config", dummyError},
		{"tests/bad.config", dummyError},
		{"tests/good.config", nil},
	}

	for _, test := range readTests {
		err := Read(test.path)
		if test.err != nil && err == nil {
			t.Fatalf("Read(%s): Failure expected, but got OK", test.path)
		}
		if test.err == nil && err != nil {
			t.Fatalf("Read(%s): OK expected, but got failure", test.path)
		}
	}

}

func TestStr(t *testing.T) {

	type strTest struct {
		scope, key, val string
		err             error
	}

	var strTests = []strTest{
		{"no such scope", "key", "", dummyError},
		{"scope #1", "str key", "value1", nil},
		{"scope #2", "str key", "value2", nil},
		{"scope #1", "int key", "1", nil},
		{"end", "test key", "test value", nil},
	}

	for _, test := range strTests {

		val, err := Str(test.scope, test.key)

		if err != nil && test.err == nil {
			t.Fatalf("Str(%s, %s): OK expected, but got failure", test.scope, test.key)
		}

		if err == nil && test.err != nil {
			t.Fatalf("Str(%s, %s): Failure expected, but got OK", test.scope, test.key)
		}

		if err == nil && val != test.val {
			t.Fatalf("Str(%s, %s): %s != %s", test.scope, test.key, val, test.val)
		}

	}

}

func TestInt(t *testing.T) {

	type intTest struct {
		scope, key string
		val        int64
		err        error
	}

	var intTests = []intTest{
		{"scope #1", "str key", 0, dummyError},
		{"scope #1", "int key", 1, nil},
		{"scope #2", "int key", 2, nil},
		{"scope #1", "float key", 0, dummyError},
	}

	for _, test := range intTests {

		val, err := Int(test.scope, test.key)

		if err != nil && test.err == nil {
			t.Fatalf("Int(%s, %s): OK expected, but got failure", test.scope, test.key)
		}

		if err == nil && test.err != nil {
			t.Fatalf("Int(%s, %s): Failure expected, but got OK", test.scope, test.key)
		}

		if err == nil && val != test.val {
			t.Fatalf("Int(%s, %s): %d != %d", test.scope, test.key, val, test.val)
		}

	}

}

func TestFloat(t *testing.T) {

	type floatTest struct {
		scope, key string
		val        float64
		err        error
	}

	var floatTests = []floatTest{
		{"scope #1", "str key", 0, dummyError},
		{"scope #1", "float key", 1.1, nil},
		{"scope #2", "float key", 2.2, nil},
		{"scope #1", "bool key", -1, dummyError},
	}

	for _, test := range floatTests {

		val, err := Float(test.scope, test.key)

		if err != nil && test.err == nil {
			t.Fatalf("Float(%s, %s): OK expected, but got failure", test.scope, test.key)
		}

		if err == nil && test.err != nil {
			t.Fatalf("Float(%s, %s): Failure expected, but got OK", test.scope, test.key)
		}

		if err == nil && val != test.val {
			t.Fatalf("Float(%s, %s): %f != %f", test.scope, test.key, val, test.val)
		}

	}

}

func TestBool(t *testing.T) {

	type boolTest struct {
		scope, key string
		val        bool
		err        error
	}

	var boolTests = []boolTest{
		{"scope #1", "str key", true, nil},
		{"scope #1", "float key", true, nil},
		{"scope #1", "bool key", false, nil},
		{"scope #2", "bool key", true, nil},
	}

	for _, test := range boolTests {

		val, err := Bool(test.scope, test.key)

		if err != nil && test.err == nil {
			t.Fatalf("Bool(%s, %s): OK expected, but got failure", test.scope, test.key)
		}

		if err == nil && test.err != nil {
			t.Fatalf("Bool(%s, %s): Failure expected, but got OK", test.scope, test.key)
		}

		if err == nil && val != test.val {
			t.Fatalf("Bool(%s, %s): %v != %v", test.scope, test.key, val, test.val)
		}

	}

}

func TestDump(t *testing.T) {
	Dump("tests/generated-dump.conf")
}
