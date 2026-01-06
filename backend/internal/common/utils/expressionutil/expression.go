package expressionutil

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/vm"
)

// --- Static Analysis Visitor & Patcher ---

// variableVisitor implements ast.Visitor to collect identifiers from an expression.
type variableVisitor struct {
	variables map[string]struct{}
}

func (v *variableVisitor) Visit(node *ast.Node) {
	if identifier, ok := (*node).(*ast.IdentifierNode); ok {
		v.variables[identifier.Value] = struct{}{}
	}
}

// newVariableVisitor creates a new visitor for collecting variables.
func newVariableVisitor() *variableVisitor {
	return &variableVisitor{
		variables: make(map[string]struct{}),
	}
}

// GetVariables extracts all unique variable names from an expression string.
func GetVariables(expression string) ([]string, error) {
	visitor := newVariableVisitor()
	// We only need to parse, not fully compile, to get the AST.
	tree, err := parser.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression for variables: %w", err)
	}
	ast.Walk(&tree.Node, visitor)

	vars := make([]string, 0, len(visitor.variables))
	for v := range visitor.variables {
		vars = append(vars, v)
	}
	sort.Strings(vars) // Return a consistent order
	return vars, nil
}

// variablePatcher implements ast.Visitor to replace identifiers in an expression.
type variablePatcher struct {
	replacer func(string) string
}

func (p *variablePatcher) Visit(node *ast.Node) {
	if identifier, ok := (*node).(*ast.IdentifierNode); ok {
		replacement := p.replacer(identifier.Value)
		if replacement != identifier.Value {
			ast.Patch(node, &ast.IdentifierNode{Value: replacement})
		}
	}
}

// ReplaceExpressionFunc replaces variable names in an expression using a replacement function.
func ReplaceExpressionFunc(expression string, varReplacer func(string) string) (string, error) {
	tree, err := parser.Parse(expression)
	if err != nil {
		return "", fmt.Errorf("failed to parse expression: %w", err)
	}

	patcher := &variablePatcher{replacer: varReplacer}
	ast.Walk(&tree.Node, patcher)

	return tree.Node.String(), nil
}

// ReplaceExpression replaces variable names in an expression using a replacement map.
func ReplaceExpression(expression string, varReplacer map[string]string) (string, error) {
	replacerFunc := func(name string) string {
		if replacement, ok := varReplacer[name]; ok {
			return replacement
		}
		return name
	}
	return ReplaceExpressionFunc(expression, replacerFunc)
}

// EncloseExpressionVars encloses all variable names in an expression with a specified escape character.
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

// --- Expression Execution ---

// customFunctions defines custom functions that are not built into the 'expr' library.
var customFunctions = []expr.Option{
	expr.Function("AVERAGE",
		func(args ...any) (any, error) {
			if len(args) == 0 {
				return 0.0, nil
			}
			sum, err := variadicToFloat64Slice(args)
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
		func(args ...any) (any, error) {
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
		func(args ...any) (any, error) {
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
		func(args ...any) (any, error) {
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
		func(args ...any) (any, error) {
			if len(args) == 0 {
				return 1.0, nil
			}
			nums, err := variadicToFloat64Slice(args)
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
	expr.Function("INT",
		func(args ...any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("INT expects 1 argument")
			}
			num, ok := toFloat64(args[0])
			if !ok {
				return nil, fmt.Errorf("invalid argument for INT")
			}
			return int64(num), nil
		},
		new(func(float64) int64),
	),
	expr.Function("SUMPRODUCT",
		func(args ...any) (any, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("SUMPRODUCT expects 2 arguments (slice1, slice2)")
			}
			slice1, err1 := argToFloat64Slice(args[0])
			slice2, err2 := argToFloat64Slice(args[1])
			if err1 != nil {
				return nil, fmt.Errorf("first argument to SUMPRODUCT must be a slice of numbers: %w", err1)
			}
			if err2 != nil {
				return nil, fmt.Errorf("second argument to SUMPRODUCT must be a slice of numbers: %w", err2)
			}
			if len(slice1) != len(slice2) {
				return nil, fmt.Errorf("arrays must have the same length")
			}
			var sum float64
			for i := range slice1 {
				sum += slice1[i] * slice2[i]
			}
			return sum, nil
		},
		new(func([]any, []any) float64),
	),
	// Add common Math functions for compatibility
	expr.Function("SQRT",
		func(args ...any) (any, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("SQRT expects 1 argument")
			}
			num, ok := toFloat64(args[0])
			if !ok {
				return nil, fmt.Errorf("invalid argument for SQRT")
			}
			return math.Sqrt(num), nil
		},
		new(func(float64) float64),
	),
	expr.Function("SIN", func(args ...any) (any, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("SIN expects 1 argument")
		}
		num, ok := toFloat64(args[0])
		if !ok {
			return nil, fmt.Errorf("invalid argument for SIN")
		}
		return math.Sin(num), nil
	}, new(func(float64) float64)),
	expr.Function("COS", func(args ...any) (any, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("COS expects 1 argument")
		}
		num, ok := toFloat64(args[0])
		if !ok {
			return nil, fmt.Errorf("invalid argument for COS")
		}
		return math.Cos(num), nil
	}, new(func(float64) float64)),
	expr.Function("TAN", func(args ...any) (any, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("TAN expects 1 argument")
		}
		num, ok := toFloat64(args[0])
		if !ok {
			return nil, fmt.Errorf("invalid argument for TAN")
		}
		return math.Tan(num), nil
	}, new(func(float64) float64)),
	// Note: 'expr' has built-in 'sum', 'max', 'min', 'len' (for COUNT), 'abs'.
	// POWER can be done with the '**' operator.
	// Logical ops (AND, OR, NOT, XOR) and IF are built-in operators/syntax.
}

// CompileExpression compiles an expression string and returns a runnable program.
func CompileExpression(expression string) (*vm.Program, error) {
	// Add custom functions to the options
	program, err := expr.Compile(expression, customFunctions...)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}
	return program, nil
}

// ExecuteExpression executes a compiled program with the given environment variables.
func ExecuteExpression(program *vm.Program, vars map[string]any) (any, error) {
	result, err := expr.Run(program, vars)
	if err != nil {
		return nil, fmt.Errorf("failed to execute expression: %w", err)
	}
	return result, nil
}

// EvalExpression compiles and runs an expression in one step.
func EvalExpression(expression string, vars map[string]any) (any, error) {
	opts := append([]expr.Option{expr.Env(vars)}, customFunctions...)
	output, err := expr.Compile(expression, opts...)
	if err != nil {
		return nil, err
	}
	return expr.Run(output, vars)
}

// --- Type conversion helpers ---

func toFloat64(v any) (float64, bool) {
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

func toInt(v any) (int, bool) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(val.Int()), true
	default:
		return 0, false
	}
}

// variadicToFloat64Slice converts variadic arguments (...any) to a slice of float64.
func variadicToFloat64Slice(args []any) ([]float64, error) {
	nums := make([]float64, len(args))
	for i, arg := range args {
		num, ok := toFloat64(arg)
		if !ok {
			return nil, fmt.Errorf("all arguments must be numbers")
		}
		nums[i] = num
	}
	return nums, nil
}

// argToFloat64Slice converts a single argument which is a slice into a slice of float64.
func argToFloat64Slice(arg any) ([]float64, error) {
	val := reflect.ValueOf(arg)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("argument is not a slice")
	}

	nums := make([]float64, val.Len())
	for i := 0; i < val.Len(); i++ {
		num, ok := toFloat64(val.Index(i).Interface())
		if !ok {
			return nil, fmt.Errorf("all items in slice must be numbers")
		}
		nums[i] = num
	}
	return nums, nil
}

// --- Expression Static Analysis (Legacy Stubs Removed) ---

// ParseResult represents the result of parsing an expression to extract metadata.
type ParseResult struct {
	Vars               []string // Variable names found in the expression
	Functions          []string // Function names found in the expression
	AggregateFunctions []string // Aggregate function names (e.g., SUM, AVG)
	IsBooleanResult    bool     // Whether the expression returns a boolean
}

// ParseExpression parses an expression and extracts variables, functions, and other metadata.
func ParseExpression(expression string) (*ParseResult, error) {
	vars, err := GetVariables(expression)
	if err != nil {
		return nil, err
	}
	// TODO: Implement function extraction and boolean result detection.
	return &ParseResult{
		Vars:               vars,
		Functions:          []string{},
		AggregateFunctions: []string{},
		IsBooleanResult:    false, // This would require type-checking the AST, which is more complex.
	}, nil
}
