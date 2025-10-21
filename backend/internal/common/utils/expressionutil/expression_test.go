package expressionutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVariables(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		expected   []string
		expectErr  bool
	}{
		{
			name:       "Simple expression",
			expression: "a + b > 5 && c",
			expected:   []string{"a", "b", "c"},
			expectErr:  false,
		},
		{
			name:       "With function calls",
			expression: "SQRT(x*x + y*y)",
			expected:   []string{"SQRT", "x", "y"},
			expectErr:  false,
		},
		{
			name:       "Duplicate variables",
			expression: "var1 + var2 - var1",
			expected:   []string{"var1", "var2"},
			expectErr:  false,
		},
		{
			name:       "Empty expression",
			expression: "",
			expected:   nil,
			expectErr:  true,
		},
		{
			name:       "Invalid expression",
			expression: "a + b *",
			expected:   nil,
			expectErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			vars, err := GetVariables(tc.expression)

			if tc.expectErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.Equal(tc.expected, vars)
		})
	}
}

func TestReplaceExpression(t *testing.T) {
	testCases := []struct {
		name         string
		expression   string
		replacements map[string]string
		expected     string
	}{
		{
			name:         "Simple replacement",
			expression:   "a + b",
			replacements: map[string]string{"a": "x", "b": "y"},
			expected:     "x + y",
		},
		{
			name:         "Partial replacement",
			expression:   "var1 > 10 && var2 < 5",
			replacements: map[string]string{"var1": "temp"},
			expected:     "temp > 10 && var2 < 5",
		},
		{
			name:         "No matching variables",
			expression:   "c - d",
			replacements: map[string]string{"a": "x", "b": "y"},
			expected:     "c - d",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result, err := ReplaceExpression(tc.expression, tc.replacements)
			assert.NoError(err)
			assert.Equal(tc.expected, result)
		})
	}
}

func TestCustomFunctions(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		vars       map[string]any
		expected   any
	}{
		{name: "INT", expression: "INT(5.7)", expected: int64(5)},
		{name: "SUMPRODUCT", expression: "SUMPRODUCT(arr1, arr2)", vars: map[string]any{
			"arr1": []any{1.0, 2.0, 3.0},
			"arr2": []any{4.0, 5.0, 6.0},
		}, expected: 32.0},
		{name: "SQRT", expression: "SQRT(16)", expected: 4.0},
		{name: "SIN", expression: "SIN(0)", expected: 0.0},
		{name: "COS", expression: "COS(0)", expected: 1.0},
		{name: "TAN", expression: "TAN(0)", expected: 0.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			result, err := EvalExpression(tc.expression, tc.vars)
			assert.NoError(err, "unexpected error executing '%s'", tc.expression)

			if fExpected, ok := tc.expected.(float64); ok {
				assert.InDelta(fExpected, result, 1e-9)
			} else {
				assert.Equal(tc.expected, result)
			}
		})
	}
}
