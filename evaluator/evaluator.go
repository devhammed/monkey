package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

// Builtin singletons
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// the "init" function is necessary to prevent initialization loop error.
// see https://stackoverflow.com/questions/51667411/initialization-loop-golang#51667738

var builtins = map[string]*object.Builtin{}

func init() {
	builtins = map[string]*object.Builtin{
		"len": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				switch arg := args[0].(type) {
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				default:
					return newError("argument to `len` not supported, got %s",
						args[0].Type())
				}
			},
		},
		"print": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			},
		},
		"array_first": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `first` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)

				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}

				return NULL
			},
		},
		"array_last": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `last` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)

				if length > 0 {
					return arr.Elements[length-1]
				}

				return NULL
			},
		},
		"array_rest": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `rest` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)

				if length > 0 {
					newElements := make([]object.Object, length-1, length-1)
					copy(newElements, arr.Elements[1:length])

					return &object.Array{Elements: newElements}
				}

				return NULL
			},
		},
		"array_push": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `push` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				newElements := make([]object.Object, length+1, length+1)

				copy(newElements, arr.Elements)

				newElements[length] = args[1]

				arr.Elements = newElements

				return NULL
			},
		},
		"array_map": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 2 {
					return newError("wrong number of arguments. got=%d, expected at least=2",
						len(args))
				}

				if len(args) > 3 {
					return newError("wrong number of arguments. got=%d, expected max=3",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("first argument to `map` must be ARRAY, got %s",
						args[0].Type())
				}

				if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
					return newError("second argument to `map` must be FUNCTION, got %s",
						args[1].Type())
				}

				arr := args[0].(*object.Array)
				elements := arr.Elements
				length := len(elements)
				newElements := make([]object.Object, length, length)

				for i := 0; i < length; i++ {
					fnArgs := []object.Object{elements[i], &object.Integer{Value: int64(i)}}

					if len(args) == 3 {
						fnArgs = append(fnArgs, arr)
					}

					newElements[i] = applyFunction(args[1], fnArgs)
				}

				return &object.Array{Elements: newElements}
			},
		},
		"array_copy": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `push` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				length := len(arr.Elements)
				newElements := make([]object.Object, length, length)

				copy(newElements, arr.Elements)

				return &object.Array{Elements: newElements}
			},
		},
	}
}

// Eval evaluates the AST passed
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// Expressions
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)

		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)

		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)

		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)

		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)

		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	}

	return nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		argsLength := len(args)
		parametersLength := len(fn.Parameters)

		if argsLength != parametersLength {
			return newError("wrong number of arguments passed to function. got=%d, expected=%d",
				argsLength, parametersLength)
		}

		extendedEnv := extendFunctionEnv(fn, args)
		extendedEnv.Set("arguments", &object.Array{Elements: args})
		evaluated := Eval(fn.Body, extendedEnv)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	if val, ok := env.Get(node.Value); ok {
		return val
	}

	return newError("identifier not found: " + node.Value)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()

			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}
