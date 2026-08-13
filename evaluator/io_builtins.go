package evaluator

import (
	"encoding/json"
	"fmt"
	"io"
	"monkey/object"
	"monkey/typing"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {
	builtins["print"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"println",
				args,
				typing.MinimumArgs(1),
			); err != nil {
				return newError("%s", err.Error())
			}

			stdout, ok := env.Get("STDOUT")
			if !ok {
				return newError("println: STDOUT is not defined")
			}

			resource, ok := stdout.Value.(*object.Resource)
			if !ok {
				return newError("println: STDOUT is not a resource")
			}

			writer, ok := resource.Handle.(io.Writer)
			if !ok {
				return newError("println: STDOUT is not writable")
			}

			for _, arg := range args {
				if _, err := io.WriteString(writer, arg.Inspect()); err != nil {
					return newError("println: %s", err)
				}
			}

			if _, err := io.WriteString(writer, "\n"); err != nil {
				return newError("println: %s", err)
			}

			return NULL
		},
	}

	builtins["input"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"input",
				args,
				typing.RangeOfArgs(0, 1),
				typing.WithTypes(object.STRING_OBJ),
			); err != nil {
				return newError("%s", err.Error())
			}

			if len(args) == 1 {
				stdout, ok := env.Get("STDOUT")
				if !ok {
					return newError("input: STDOUT is not defined")
				}

				resource, ok := stdout.Value.(*object.Resource)
				if !ok {
					return newError("input: STDOUT is not a resource")
				}

				writer, ok := resource.Handle.(io.Writer)
				if !ok {
					return newError("input: STDOUT is not writable")
				}

				if _, err := io.WriteString(writer, args[0].(*object.String).Value); err != nil {
					return newError("input: %s", err)
				}
			}

			stdin, ok := env.Get("STDIN")
			if !ok {
				return newError("input: STDIN is not defined")
			}

			resource, ok := stdin.Value.(*object.Resource)
			if !ok {
				return newError("input: STDIN is not a resource")
			}

			reader, ok := resource.Handle.(io.Reader)
			if !ok {
				return newError("input: STDIN is not readable")
			}

			var line strings.Builder
			var buffer [1]byte
			for {
				n, err := reader.Read(buffer[:])
				if n > 0 {
					if buffer[0] == '\n' {
						break
					}
					line.WriteByte(buffer[0])
				}
				if err != nil {
					if err != io.EOF {
						return newError("input: %s", err)
					}
					break
				}
			}

			return &object.String{Value: strings.TrimSuffix(line.String(), "\r")}
		},
	}

	builtins["kind"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("kind", args, typing.ExactArgs(1), typing.WithTypes(object.RESOURCE_OBJ)); err != nil {
				return newError("%s", err.Error())
			}

			return &object.String{Value: args[0].(*object.Resource).Kind}
		},
	}

	builtins["open"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("open", args, typing.RangeOfArgs(1, 2), typing.AllOfType(object.STRING_OBJ)); err != nil {
				return newError("%s", err.Error())
			}

			path := args[0].(*object.String).Value

			if filepath.IsAbs(path) || filepath.VolumeName(path) != "" {
				path = "file://" + filepath.ToSlash(path)
			}

			u, err := url.Parse(path)
			if err != nil {
				return newError("%s", err.Error())
			}

			if u.Scheme == "" {
				u.Scheme = "file"
				u.Path = path
			}

			mode := ""

			if len(args) == 2 {
				mode = args[1].(*object.String).Value
			}

			switch u.Scheme {
			case "file":
				var flags int

				if mode == "" {
					mode = "r"
				}

				switch mode {
				case "r":
					flags = os.O_RDONLY
				case "r+":
					flags = os.O_RDWR
				case "w":
					flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
				case "w+":
					flags = os.O_RDWR | os.O_CREATE | os.O_TRUNC
				case "a":
					flags = os.O_WRONLY | os.O_CREATE | os.O_APPEND
				case "a+":
					flags = os.O_RDWR | os.O_CREATE | os.O_APPEND
				default:
					return newError("invalid mode: %q", mode)
				}

				return &object.Resource{Kind: "FILE", Handle: os.OpenFile(filepath.FromSlash(u.Path), flags, 0644)}
			default:
				return newError("unsupported protocol: %q", u.Scheme)
			}
		},
	}

	builtins["json_encode"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check("json_encode", args, typing.ExactArgs(1)); err != nil {
				return newError("%s", err.Error())
			}

			value, err := objectToJson(args[0])
			if err != nil {
				return newError("json_encode: %s", err)
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				return newError("json_encode: %s", err)
			}
			return &object.String{Value: string(encoded)}
		},
	}

	builtins["json_decode"] = &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if err := typing.Check(
				"json_decode",
				args,
				typing.ExactArgs(1),
				typing.WithTypes(object.STRING_OBJ),
			); err != nil {
				return newError("%s", err.Error())
			}

			decoder := json.NewDecoder(strings.NewReader(args[0].(*object.String).Value))
			decoder.UseNumber()
			var value any
			if err := decoder.Decode(&value); err != nil {
				return newError("json_decode: %s", err)
			}
			// Decode rejects trailing non-whitespace JSON, matching json.Unmarshal.
			var trailing any
			if err := decoder.Decode(&trailing); err != io.EOF {
				if err == nil {
					return newError("json_decode: invalid trailing data")
				}
				return newError("json_decode: %s", err)
			}

			decoded, err := jsonToObject(value)
			if err != nil {
				return newError("json_decode: %s", err)
			}
			return decoded
		},
	}
}

func objectToJson(value object.Object) (any, error) {
	switch value := value.(type) {
	case *object.Null:
		return nil, nil
	case *object.Boolean:
		return value.Value, nil
	case *object.Integer:
		return value.Value, nil
	case *object.Float:
		return value.Value, nil
	case *object.String:
		return value.Value, nil
	case *object.Array:
		result := make([]any, len(value.Elements))
		for i, element := range value.Elements {
			converted, err := objectToJson(element)
			if err != nil {
				return nil, err
			}
			result[i] = converted
		}
		return result, nil
	case *object.Hash:
		result := make(map[string]any, len(value.Pairs))
		for _, pair := range value.Pairs {
			converted, err := objectToJson(pair.Value)
			if err != nil {
				return nil, err
			}
			result[pair.Key.Inspect()] = converted
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported object type %s", value.Type())
	}
}

func jsonToObject(value any) (object.Object, error) {
	switch value := value.(type) {
	case nil:
		return NULL, nil
	case bool:
		return nativeBoolToBooleanObject(value), nil
	case string:
		return &object.String{Value: value}, nil
	case json.Number:
		if !strings.ContainsAny(string(value), ".eE") {
			if integer, err := strconv.ParseInt(string(value), 10, 64); err == nil {
				return &object.Integer{Value: integer}, nil
			}
		}
		if number, err := strconv.ParseFloat(string(value), 64); err == nil {
			return &object.Float{Value: number}, nil
		}
		return nil, fmt.Errorf("JSON number %q is not a valid number", value)
	case []any:
		elements := make([]object.Object, len(value))
		for i, element := range value {
			converted, err := jsonToObject(element)
			if err != nil {
				return nil, err
			}
			elements[i] = converted
		}
		return &object.Array{Elements: elements}, nil
	case map[string]any:
		pairs := make(map[object.HashKey]object.HashPair, len(value))
		for key, element := range value {
			converted, err := jsonToObject(element)
			if err != nil {
				return nil, err
			}
			keyObject := &object.String{Value: key}
			hashKey, err := keyObject.HashKey()
			if err != nil {
				return nil, err
			}
			pairs[hashKey] = object.HashPair{Key: keyObject, Value: converted}
		}
		return &object.Hash{Pairs: pairs}, nil
	default:
		return nil, fmt.Errorf("unsupported JSON value %T", value)
	}
}
