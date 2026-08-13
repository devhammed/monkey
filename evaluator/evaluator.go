package evaluator

import (
	"fmt"
	"math"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/typing"
	"os"
)

// Builtin singletons
var (
	NULL     = &object.Null{}
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	VERSION  = &object.String{Value: "v0.2.7"}
	ARGV     = &object.Array{}
	STDIN    = &object.Resource{Kind: "stdin", Handle: os.Stdin}
	STDOUT   = &object.Resource{Kind: "stdout", Handle: os.Stdout}
	STDERR   = &object.Resource{Kind: "stderr", Handle: os.Stderr}
	E        = &object.Float{Value: math.E}
	PI       = &object.Float{Value: math.Pi}
	PHI      = &object.Float{Value: math.Phi}
	SQRT2    = &object.Float{Value: math.Sqrt2}
	SQRTE    = &object.Float{Value: math.SqrtE}
	SQRTPI   = &object.Float{Value: math.SqrtPi}
	SQRTPHI  = &object.Float{Value: math.SqrtPhi}
	LN2      = &object.Float{Value: math.Ln2}
	LOG2E    = &object.Float{Value: math.Log2E}
	LN10     = &object.Float{Value: math.Ln10}
	LOG10E   = &object.Float{Value: math.Log10E}
	builtins = map[string]*object.Builtin{}
)

func init() {
	for _, arg := range os.Args {
		ARGV.Elements = append(ARGV.Elements, &object.String{Value: arg})
	}
}

// Run lexes, parses, and evaluates code
func Run(
	code string,
	file string,
	dir string,
	isMain bool,
	env *object.Environment,
) object.Object {
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()
	errors := p.Errors()

	if len(errors) != 0 {
		fmt.Println("Woops! We ran into some monkey business here!")
		fmt.Println(" parser errors:")

		for _, msg := range errors {
			fmt.Println("\t" + msg)
		}

		return nil
	}

	if _, ok := env.Get("ARGV"); !ok {
		env.Set("ARGV", ARGV, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("STDIN"); !ok {
		env.Set("STDIN", STDIN, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("STDOUT"); !ok {
		env.Set("STDOUT", STDOUT, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("STDERR"); !ok {
		env.Set("STDERR", STDERR, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("VERSION"); !ok {
		env.Set("VERSION", VERSION, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("E"); !ok {
		env.Set("E", E, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("PI"); !ok {
		env.Set("PI", PI, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("PHI"); !ok {
		env.Set("PHI", PHI, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("SQRT2"); !ok {
		env.Set("SQRT2", SQRT2, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("SQRTE"); !ok {
		env.Set("SQRTE", SQRTE, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("SQRTPI"); !ok {
		env.Set("SQRTPI", SQRTPI, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("SQRTPHI"); !ok {
		env.Set("SQRTPHI", SQRTPHI, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("LN2"); !ok {
		env.Set("LN2", LN2, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("LOG2E"); !ok {
		env.Set("LOG2E", LOG2E, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("LN10"); !ok {
		env.Set("LN10", LN10, object.BindingOptions{SuperGlobal: true})
	}
	if _, ok := env.Get("LOG10E"); !ok {
		env.Set("LOG10E", LOG10E, object.BindingOptions{SuperGlobal: true})
	}

	env.Set("MAIN", nativeBoolToBooleanObject(isMain), object.BindingOptions{SuperGlobal: true})
	env.Set("FILE", &object.String{Value: file}, object.BindingOptions{SuperGlobal: true})
	env.Set("DIR", &object.String{Value: dir}, object.BindingOptions{SuperGlobal: true})

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
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
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
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		name := ""

		if ident, ok := node.Function.(*ast.Identifier); ok {
			name = ident.Value
		}

		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)

		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, env, name)
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
		value := Eval(node.Value, env)

		if isError(value) {
			return value
		}

		if ident, ok := node.Left.(*ast.Identifier); ok {
			if v, ok := env.Get(ident.Value); ok && v.SuperGlobal {
				return newError("cannot reassign a superglobal")
			}

			if immutable, ok := value.(object.Immutable); ok {
				env.Set(ident.Value, immutable.Clone(), object.BindingOptions{})
			} else {
				env.Set(ident.Value, value, object.BindingOptions{})
			}

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
					hashed, err := hashKey.HashKey()

					if err != nil {
						return newError("hash key error: %s", err.Error())
					}

					hash.Pairs[hashed] = object.HashPair{Key: key, Value: value}

					return NULL
				}

				return newError("cannot index hash with %T", key)
			}

			return newError("object type %s does not support item assignment", obj.Type())
		}

		left := Eval(node.Left, env)

		if isError(left) {
			return left
		}

		return newError("expected identifier or index expression got=%T", left)
	}

	return nil
}

func newError(format string, a ...any) *object.Error {
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

func applyFunction(fn object.Object, args []object.Object, env *object.Environment, name string) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if name == "" {
			name = "(anonymous)"
		}

		parametersLength := len(fn.Parameters)

		if err := typing.Check(
			name,
			args,
			typing.MinimumArgs(parametersLength),
		); err != nil {
			return newError("%s", err.Error())
		}

		extendedEnv := extendFunctionEnv(fn, args)

		extendedEnv.Set("arguments", &object.Array{Elements: args}, object.BindingOptions{})

		evaluated := Eval(fn.Body, extendedEnv)

		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(env, args...)
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
		env.Set(param.Value, args[paramIdx], object.BindingOptions{})
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
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
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

		hashed, err := hashKey.HashKey()

		if err != nil {
			return newError("hash key error: %s", err.Error())
		}

		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hash{Pairs: pairs}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	maxIndex := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > maxIndex {
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

	hashKey, err := key.HashKey()
	if err != nil {
		return newError("hash key error: %s", err.Error())
	}

	pair, ok := hashObject.Pairs[hashKey]

	if !ok {
		return NULL
	}

	return pair.Value
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	maxIndex := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > maxIndex {
		return &object.String{Value: ""}
	}

	return &object.String{Value: string(stringObject.Value[idx])}
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
		return val.Value
	}

	return newError("identifier not found: %s", node.Value)
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
	}

	return NULL
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
	case isNumeric(left) && isNumeric(right):
		return evalNumericInfixExpression(operator, left, right)
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

func isCallable(value object.Object) bool {
	return value.Type() == object.FUNCTION_OBJ || value.Type() == object.BUILTIN_OBJ
}

func isNumeric(value object.Object) bool {
	return value.Type() == object.INTEGER_OBJ || value.Type() == object.FLOAT_OBJ
}

func numericValue(value object.Object) (float64, error) {
	switch value := value.(type) {
	case *object.Integer:
		return float64(value.Value), nil
	case *object.Float:
		return value.Value, nil
	default:
		return 0, fmt.Errorf("%s is not a number", value.Type())
	}
}

func evalNumericInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal, err := numericValue(left)
	if err != nil {
		return newError("%s", err)
	}
	rightVal, err := numericValue(right)
	if err != nil {
		return newError("%s", err)
	}

	leftInteger, leftIsInteger := left.(*object.Integer)
	rightInteger, rightIsInteger := right.(*object.Integer)
	bothIntegers := leftIsInteger && rightIsInteger

	switch operator {
	case "+":
		if bothIntegers {
			return &object.Integer{Value: leftInteger.Value + rightInteger.Value}
		}
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		if bothIntegers {
			return &object.Integer{Value: leftInteger.Value - rightInteger.Value}
		}
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		if bothIntegers {
			return &object.Integer{Value: leftInteger.Value * rightInteger.Value}
		}
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightInteger.Value == 0 {
			return newError("division by zero")
		}
		if bothIntegers && leftInteger.Value%rightInteger.Value == 0 {
			return &object.Integer{Value: leftInteger.Value / rightInteger.Value}
		}
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
	if right.Type() == object.FLOAT_OBJ {
		return &object.Float{Value: -right.(*object.Float).Value}
	}
	if right.Type() == object.INTEGER_OBJ {
		return &object.Integer{Value: -right.(*object.Integer).Value}
	}

	return newError("unknown operator: -%s", right.Type())
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
