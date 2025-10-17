package utils

import (
	"fmt"
	"math"
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// TODO: ExpressionUtils.java contains logic for static analysis of expressions (parsing variables, function names).
// The 'expr' library does not provide a simple public API for this. This part of the migration is postponed.
// The focus here is on executing expressions, which corresponds to ExpressionFunctions.java.

// customFunctions defines custom functions that are not built into the 'expr' library but were present in ExpressionFunctions.java.
var customFunctions = []expr.Option{
	expr.Function("AVERAGE",
		func(args ...interface{}) (interface{}, error) {
			if len(args) == 0 {
				return 0.0, nil
			}
			sum, err := toFloat64Slice(args)
			if err != nil {
				return nil, err
			}
			var total float64
			for _, v := range sum {
				total += v
			}
			return total / float64(len(sum)), nil
		},
		new(func(...float64) float64),
	),
	expr.Function("FIXED",
		func(args ...interface{}) (interface{}, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("FIXED expects 2 arguments (number, scale)")
			}
			num, ok1 := toFloat64(args[0])
			scale, ok2 := toInt(args[1])
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("invalid arguments for FIXED")
			}
			// Simple rounding, not exactly like BigDecimal but sufficient for many cases.
			return math.Round(num*math.Pow10(scale)) / math.Pow10(scale), nil
		},
		new(func(float64, int) float64),
	),
	expr.Function("LOG",
		func(args ...interface{}) (interface{}, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("LOG expects 2 arguments (number, base)")
			}
			num, ok1 := toFloat64(args[0])
			base, ok2 := toFloat64(args[1])
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("invalid arguments for LOG")
			}
			return math.Log(num) / math.Log(base), nil
		},
		new(func(float64, float64) float64),
	),
	expr.Function("MOD",
		func(args ...interface{}) (interface{}, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("MOD expects 2 arguments")
			}
			a, ok1 := toFloat64(args[0])
			b, ok2 := toFloat64(args[1])
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("invalid arguments for MOD")
			}
			return math.Mod(a, b), nil
		},
		new(func(float64, float64) float64),
	),
	expr.Function("PRODUCT",
		func(args ...interface{}) (interface{}, error) {
			if len(args) == 0 {
				return 1.0, nil
			}
			nums, err := toFloat64Slice(args)
			if err != nil {
				return nil, err
			}
			var product float64 = 1
			for _, v := range nums {
				product *= v
			}
			return product, nil
		},
		new(func(...float64) float64),
	),
	// Note: 'expr' has built-in 'sum', 'max', 'min', 'len' (for COUNT).
	// 'abs' is also built-in. POWER can be done with the '**' operator.
}

// CompileExpression compiles an expression string and returns a runnable program.
// Corresponds to ExpressionFunctions.compileExpression.
func CompileExpression(expression string) (*vm.Program, error) {
	// Add custom functions to the options
	program, err := expr.Compile(expression, customFunctions...)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}
	return program, nil
}

// ExecuteExpression executes a compiled program with the given environment variables.
// Corresponds to ExpressionFunctions.executeExpression.
func ExecuteExpression(program *vm.Program, vars map[string]interface{}) (interface{}, error) {
	result, err := expr.Run(program, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to execute expression: %w", err)
	}
	return result, nil
}

// EvalExpression compiles and runs an expression in one step.
func EvalExpression(expression string, vars map[string]interface{}) (interface{}, error) {
	opts := append([]expr.Option{expr.Env(vars)}, customFunctions...)
	output, err := expr.Compile(expression, opts...)
	if err != nil {
		return nil, err
	}
	return expr.Run(output, vars)
}

// --- Type conversion helpers ---

func toFloat64(v interface{}) (float64, bool) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), true
	case reflect.Float32, reflect.Float64:
		return val.Float(), true
	default:
		return 0, false
	}
}

func toInt(v interface{}) (int, bool) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int()), true
	default:
		return 0, false
	}
}

func toFloat64Slice(args []interface{}) ([]float64, error) {
	var nums []float64
	for _, arg := range args {
		if slice, ok := arg.([]interface{}); ok {
			// Handle case where a slice is passed as a single argument
			for _, item := range slice {
				num, ok := toFloat64(item)
				if !ok {
					return nil, fmt.Errorf("all arguments must be numbers")
				}
				nums = append(nums, num)
			}
		} else {
			num, ok := toFloat64(arg)
			if !ok {
				return nil, fmt.Errorf("all arguments must be numbers")
			}
			nums = append(nums, num)
		}
	}
	return nums, nil
}

// --- Expression Static Analysis (ExpressionUtils.java) ---

// ParseResult represents the result of parsing an expression to extract metadata.
// Corresponds to ExpressionUtils.ParseResult
type ParseResult struct {
	Vars               []string // Variable names found in the expression
	Functions          []string // Function names found in the expression
	AggregateFunctions []string // Aggregate function names (e.g., SUM, AVG)
	IsBooleanResult    bool     // Whether the expression returns a boolean
}

// TODO: The Java implementation uses Alibaba Druid SQL Parser to statically analyze expressions.
// This allows parsing variable names, function names, and determining if the expression is boolean.
//
// The 'expr' library used for Go does not expose a simple public API for AST traversal or static analysis.
// To fully implement this, we would need one of the following approaches:
//
// 1. Use a SQL parser library for Go (e.g., github.com/xwb1989/sqlparser, github.com/pingcap/parser)
//    and adapt it to parse expressions as SQL SELECT items.
//
// 2. Use the internal AST nodes from 'expr' library (if exposed) to traverse the parsed tree.
//
// 3. Implement a custom lexer/parser for the expression language subset we need.
//
// For now, these functions are stubs that return empty results or errors.

// ParseExpression parses an expression and extracts variables, functions, and other metadata.
// Corresponds to ExpressionUtils.parseExpression
func ParseExpression(expression string) (*ParseResult, error) {
	// TODO: Implement static analysis using a SQL parser or expression parser.
	// The Java version uses Druid SQL Parser to parse "SELECT <expression> FROM `__table`"
	// and visits the AST nodes to extract:
	// - SQLIdentifierExpr for variables
	// - SQLMethodInvokeExpr for functions
	// - SQLAggregateExpr for aggregate functions
	// - SQLBinaryOpExpr to determine if result is boolean (relational/logical operators)

	return &ParseResult{
		Vars:               []string{},
		Functions:          []string{},
		AggregateFunctions: []string{},
		IsBooleanResult:    false,
	}, fmt.Errorf("ParseExpression not implemented: requires SQL parser integration")
}

// ReplaceExpression replaces variable names in an expression using a replacement map.
// Corresponds to ExpressionUtils.replaceExpression(String, Map<String, String>)
func ReplaceExpression(expression string, varReplacer map[string]string) (string, error) {
	// TODO: Implement variable replacement using AST manipulation.
	// The Java version parses the expression, visits identifier nodes, and replaces them.

	return "", fmt.Errorf("ReplaceExpression not implemented: requires SQL parser integration")
}

// ReplaceExpressionFunc replaces variable names in an expression using a replacement function.
// Corresponds to ExpressionUtils.replaceExpression(String, Function<String, String>)
func ReplaceExpressionFunc(expression string, varReplacer func(string) string) (string, error) {
	// TODO: Implement variable replacement using AST manipulation.

	return "", fmt.Errorf("ReplaceExpressionFunc not implemented: requires SQL parser integration")
}

// EncloseExpressionVars encloses all variable names in an expression with a specified escape character.
// Corresponds to ExpressionUtils.encloseExpressionVars
func EncloseExpressionVars(expression string, escape rune) (string, error) {
	// This is a special case of ReplaceExpressionFunc
	replacer := func(name string) string {
		if len(name) > 0 && rune(name[0]) != escape {
			return string(escape) + name + string(escape)
		}
		return name
	}
	return ReplaceExpressionFunc(expression, replacer)
}
