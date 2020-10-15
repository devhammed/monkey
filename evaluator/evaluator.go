package evaluator

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
	"os/user"
)

// Builtin singletons
var (
	NULL           = &object.Null{}
	TRUE           = &object.Boolean{Value: true}
	FALSE          = &object.Boolean{Value: false}
	MONKEY_VERSION = &object.String{Value: "v0.2.5"}
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
		"type": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				return &object.String{Value: string(args[0].Type())}
			},
		},
		"puts": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}

				return NULL
			},
		},
		"sys_exit": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.INTEGER_OBJ {
					return newError("argument to `sys_exit` must be INTEGER, got %s",
						args[0].Type())
				}

				code := args[0].(*object.Integer)

				os.Exit(int(code.Value))

				return NULL
			},
		},
		"sys_user": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				user, err := user.Current()

				if err != nil {
					return newError("failed to get user: %s", err.Error())
				}

				return &object.String{Value: user.Username}
			},
		},
		"sys_user_name": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				user, err := user.Current()

				if err != nil {
					return newError("failed to get user name: %s", err.Error())
				}

				return &object.String{Value: user.Name}
			},
		},
		"sys_user_gid": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				user, err := user.Current()

				if err != nil {
					return newError("failed to get user GID: %s", err.Error())
				}

				return &object.String{Value: user.Gid}
			},
		},
		"sys_user_uid": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				user, err := user.Current()

				if err != nil {
					return newError("failed to get user UID: %s", err.Error())
				}

				return &object.String{Value: user.Uid}
			},
		},
		"sys_user_home": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				user, err := user.Current()

				if err != nil {
					return newError("failed to get user home: %s", err.Error())
				}

				return &object.String{Value: user.HomeDir}
			},
		},
		"sys_user_groups": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				user, err := user.Current()

				if err != nil {
					return newError("failed to get user: %s", err.Error())
				}

				groups, err := user.GroupIds()

				if err != nil {
					return newError("failed to get user groups: %s", err.Error())
				}

				var groupsArray []object.Object

				for _, group := range groups {
					groupsArray = append(groupsArray, &object.String{Value: group})
				}

				return &object.Array{Elements: groupsArray}
			},
		},
		"require": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `require` must be STRING, got %s",
						args[0].Type())
				}

				file := args[0].Inspect()
				data, err := ioutil.ReadFile(file)

				if err != nil {
					return newError("failed to require file: %s", err.Error())
				}

				env := object.NewEnvironment()

				evaluated := Run(string(data), file, FALSE, env, os.Stdout)

				if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
					return newError(
						"error in required file (%s):\n %s",
						file,
						evaluated.Inspect(),
					)
				}

				return env.ExportedHash()
			},
		},
		"file_read": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `file_read` must be STRING, got %s",
						args[0].Type())
				}

				fileName := args[0].Inspect()

				data, err := ioutil.ReadFile(fileName)

				if err != nil {
					return newError(err.Error())
				}

				return &object.String{Value: string(data)}
			},
		},
		"file_readlines": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return newError("argument to `file_read` must be STRING, got %s",
						args[0].Type())
				}

				fileName := args[0].Inspect()

				f, err := os.Open(fileName)

				if err != nil {
					return newError(err.Error())
				}

				s := bufio.NewScanner(f)

				var lines []object.Object

				for s.Scan() {
					lines = append(lines, &object.String{Value: s.Text()})
				}

				err = s.Err()

				if err != nil {
					return newError(err.Error())
				}

				if err = f.Close(); err != nil {
					return newError(err.Error())
				}

				return &object.Array{Elements: lines}
			},
		},
		"file_write": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("wrong number of arguments. got=%d, want=3",
						len(args))
				}

				if args[0].Type() != object.STRING_OBJ {
					return newError("first argument to `file_write` must be STRING, got %s",
						args[0].Type())
				}

				if args[1].Type() != object.STRING_OBJ {
					return newError("second argument to `file_write` must be STRING, got %s",
						args[1].Type())
				}

				if args[2].Type() != object.INTEGER_OBJ {
					return newError("third argument to `file_write` must be INTEGER, got %s",
						args[2].Type())
				}

				fileName := args[0].Inspect()
				filePerms := args[2].(*object.Integer)

				err := ioutil.WriteFile(
					fileName,
					[]byte(args[1].Inspect()),
					os.FileMode(filePerms.Value),
				)

				if err != nil {
					return newError(err.Error())
				}

				return NULL
			},
		},
		"range": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 2 {
					return newError("wrong number of arguments. got=%d, want=2",
						len(args))
				}

				if args[0].Type() != object.INTEGER_OBJ {
					return newError("first argument to `range` must be INTEGER, got %s",
						args[0].Type())
				}

				if args[1].Type() != object.INTEGER_OBJ {
					return newError("second argument to `range` must be INTEGER, got %s",
						args[1].Type())
				}

				if len(args) == 3 && args[2].Type() != object.INTEGER_OBJ {
					return newError("third argument to `range` must be INTEGER, got %s",
						args[1].Type())
				}

				step := int64(1)
				end := args[1].(*object.Integer)
				start := args[0].(*object.Integer)

				if len(args) == 3 {
					s := args[2].(*object.Integer)
					step = s.Value
				}

				i := start.Value
				arr := []object.Object{}

				for i < end.Value {
					arr = append(arr, &object.Integer{Value: i})
					i = i + step
				}

				return &object.Array{Elements: arr}
			},
		},
		"array_first": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `array_first` must be ARRAY, got %s",
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
					return newError("argument to `array_last` must be ARRAY, got %s",
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
					return newError("argument to `array_rest` must be ARRAY, got %s",
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
					return newError("argument to `array_push` must be ARRAY, got %s",
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
					return newError("first argument to `array_map` must be ARRAY, got %s",
						args[0].Type())
				}

				if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
					return newError("second argument to `array_map` must be FUNCTION, got %s",
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
		"array_each": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				hasThis := len(args) == 3

				if len(args) < 2 {
					return newError("wrong number of arguments. got=%d, expected at least=2",
						len(args))
				}

				if len(args) > 3 {
					return newError("wrong number of arguments. got=%d, expected max=3",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("first argument to `array_each` must be ARRAY, got %s",
						args[0].Type())
				}

				if args[1].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
					return newError("second argument to `array_each` must be FUNCTION, got %s",
						args[1].Type())
				}

				if hasThis && args[2].Type() != object.ARRAY_OBJ {
					return newError("third argument to `array_each` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				elements := arr.Elements

				for i := 0; i < len(elements); i++ {
					fnArgs := []object.Object{elements[i], &object.Integer{Value: int64(i)}}

					if hasThis {
						fnArgs = append(fnArgs, arr)
					}

					applyFunction(args[1], fnArgs)
				}

				return NULL
			},
		},
		"array_reduce": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				maxArgs := 4
				minArgs := 3
				hasThis := len(args) == maxArgs

				if len(args) < minArgs {
					return newError("wrong number of arguments. got=%d, expected at least=%d",
						len(args), minArgs)
				}

				if len(args) > maxArgs {
					return newError("wrong number of arguments. got=%d, expected max=%d",
						len(args), maxArgs)
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("first argument to `array_reduce` must be ARRAY, got %s",
						args[0].Type())
				}

				if args[2].Type() != object.FUNCTION_OBJ && args[1].Type() != object.BUILTIN_OBJ {
					return newError("third argument to `array_reduce` must be FUNCTION, got %s",
						args[1].Type())
				}

				if hasThis && args[3].Type() != object.ARRAY_OBJ {
					return newError("fourth argument to `array_reduce` must be ARRAY, got %s",
						args[0].Type())
				}

				arr := args[0].(*object.Array)
				elements := arr.Elements
				acc := args[1]

				for i := 0; i < len(elements); i++ {
					fnArgs := []object.Object{acc, elements[i], &object.Integer{Value: int64(i)}}

					if hasThis {
						fnArgs = append(fnArgs, arr)
					}

					acc = applyFunction(args[2], fnArgs)
				}

				return acc
			},
		},
		"array_copy": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1",
						len(args))
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError("argument to `array_copy` must be ARRAY, got %s",
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

// Run lexes, parses and evaluates code
func Run(
	code string,
	file string,
	isMain object.Object,
	env *object.Environment,
	out io.Writer,
) object.Object {
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()
	errors := p.Errors()

	if len(errors) != 0 {
		io.WriteString(out, "Woops! We ran into some monkey business here!\n")
		io.WriteString(out, " parser errors:\n")

		for _, msg := range errors {
			io.WriteString(out, "\t"+msg+"\n")
		}

		return nil
	}

	env.Set("MAIN", isMain)
	env.Set("MONKEY_VERSION", MONKEY_VERSION)
	env.Set("FILE", &object.String{Value: file})

	return Eval(program, env)
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
	case *ast.Null:
		return NULL
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
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
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
	case *ast.AssignmentExpression:
		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		value := Eval(node.Value, env)

		if isError(value) {
			return value
		}

		if ident, ok := node.Left.(*ast.Identifier); ok {
			env.Set(ident.Value, value)

			return NULL
		}

		if ie, ok := node.Left.(*ast.IndexExpression); ok {
			obj := Eval(ie.Left, env)

			if isError(obj) {
				return obj
			}

			if array, ok := obj.(*object.Array); ok {
				index := Eval(ie.Index, env)

				if isError(index) {
					return index
				}

				if idx, ok := index.(*object.Integer); ok {
					array.Elements[idx.Value] = value
				} else {
					return newError("cannot index array with %#v", index)
				}

				return NULL
			}

			if hash, ok := obj.(*object.Hash); ok {
				key := Eval(ie.Index, env)

				if isError(key) {
					return key
				}

				if hashKey, ok := key.(object.Hashable); ok {
					hashed := hashKey.HashKey()
					hash.Pairs[hashed] = object.HashPair{Key: key, Value: value}

					return NULL
				}

				return newError("cannot index hash with %T", key)
			}

			return newError("object type %T does not support item assignment", obj)
		}

		return newError("expected identifier or index expression got=%T", left)
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

		if argsLength < parametersLength {
			return newError(
				"number of arguments passed to function is lesser than expected. got=%d, expected=%d",
				argsLength,
				parametersLength,
			)
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
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)

		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)

		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)

		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()

		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
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

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]

	if !ok {
		return NULL
	}

	return pair.Value
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
