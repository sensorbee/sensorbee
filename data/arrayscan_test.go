package data

import (
	"strings"
	"testing"
)

// TODO: skip if key doesn't start with "foo["
func TestArraySlicingWithArray(t *testing.T) {
	data := Array{scanTestElem0, scanTestElem1, scanTestElem2}

	t.Run("valid cases", func(t *testing.T) {
		for input, expected := range arraySlicingExamples {
			if !strings.HasPrefix(input, "foo[") {
				continue
			}

			input = input[3:]
			t.Run(input, func(t *testing.T) {
				path, err := CompilePath(input)
				if err != nil {
					t.Fatalf("Cannot compile a path %v: %v", input, err)
				}
				actual, err := data.Get(path)
				if expected == nil {
					if err == nil {
						t.Errorf("Get should fail with %v", input)
					}
				} else {
					if err != nil {
						t.Fatalf("Get should succeed with %v: %v", input, err)
					}
					if !Equal(expected, actual) {
						t.Errorf("Expected: %v, Actual: %v", expected, actual)
					}
				}
			})
		}
	})

	t.Run("invalid cases", func(t *testing.T) {
		cases := []string{
			"foo",
			"foo[0]",
			`["0"]`,
		}

		for _, c := range cases {
			t.Run(c, func(t *testing.T) {
				path, err := CompilePath(c)
				if err != nil {
					t.Fatalf("Cannot compile a path %v: %v", c, err)
				}
				_, err = data.Get(path)
				if err == nil {
					t.Errorf("Get should fail with %v", c)
				}
			})
		}
	})
}
