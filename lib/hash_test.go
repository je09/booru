package lib

import (
	"reflect"
	"testing"
)

func TestAuthHash(t *testing.T) {
	s := AuthHash()
	s1 := AuthHash()

	if reflect.TypeOf(s).String() != "string" {
		t.Errorf("Expected type is string got %s", reflect.TypeOf(s).String())
	}
	if s == s1 {
		t.Errorf("Expected different values of resuls, got the same: %s", s)
	}
}

func Test_randomString(t *testing.T) {
	type testCase struct {
		name string
		args int
		want int
	}

	tests := []testCase{
		{
			"test size 1", 1, 1,
		},
		{
			"test size 32", 32, 32,
		},
		{
			"test size 512", 512, 512,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := randomString(tt.args)
			if reflect.TypeOf(s).String() != "string" {
				t.Errorf("Expected type is string got %s", reflect.TypeOf(s).String())
			}
			if len(s) != tt.want {
				t.Errorf("Expected length is %d string got %d", tt.want, len(s))
			}
		})
	}
}
