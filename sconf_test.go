package sconf

import (
	"errors"
	"testing"
)

var dummyError = errors.New("dummy")

func TestRead(t *testing.T) {

	type readTest struct {
		path string
		err  error
	}

	var readTests = []readTest{
		{"", dummyError},
		{"no such file", dummyError},
		{"bad-test.config", dummyError},
		{"good-test.config", nil},
	}

	for _, test := range readTests {
		err := Read(test.path)
		if test.err != nil && err == nil {
			t.Errorf("Read(%s): Failure expected, but got OK", test.path)
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
		{"scope #1", "str key", "value", nil},
		{"scope #1", "int key", "123", nil},

		{"end", "test key", "test value", nil},
	}

	for _, test := range strTests {

		val, err := Str(test.scope, test.key)

		if err != nil && test.err == nil {
			t.Errorf("Str(%s, %s): OK expected, but got failure", test.scope, test.key)
		}

		if err == nil && test.err != nil {
			t.Errorf("Str(%s, %s): Failure expected, but got OK", test.scope, test.key)
		}

		if err == nil && val != test.val {
			t.Errorf("Str(%s, %s): %s != %s", test.scope, test.key, val, test.val)
		}

	}

}
