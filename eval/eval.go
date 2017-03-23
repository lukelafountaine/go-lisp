package eval

import (
	"errors"
	"fmt"

	"github.com/lukelafountaine/go-lisp/types"
)

func Eval(exp types.Expression, env *types.Scope) (result types.Expression, err error) {

	switch exp := exp.(type) {

	// variable reference
	case types.Symbol:
		var scope *types.Scope
		scope, err = getSymbol(types.Symbol(exp), env)

		if err == nil {
			result = scope.Symbols[types.Symbol(exp)]
		}

	case []types.Expression:

		// make sure we have something look at
		if len(exp) == 0 {
			break
		}

		switch t := exp[0].(type) {

		// switch on the first word
		case types.Symbol:

			switch t {

			case "quote":
				result = exp[1]

			case "set!":
				if len(exp) != 3 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'set!'")
					break
				}

				key, ok := exp[1].(types.Symbol)
				if !ok {
					err = errors.New("Syntax Error: Cannot assign to a literal")
					break
				}

				var scope *types.Scope
				scope, err = getSymbol(key, env)
				if err != nil {
					break
				}

				var value types.Expression
				value, err = Eval(exp[2], env)
				if err != nil {
					break
				}

				scope.Symbols[key] = value

			case "begin":
				for _, i := range exp[1:] {
					result, err = Eval(i, env)

					if err != nil {
						break
					}
				}

			case "if":
				if len(exp) != 4 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'if'")
					break
				}

				var condition types.Expression
				condition, err = Eval(exp[1], env)

				if err != nil {
					break
				}

				consequence := exp[2]
				alternative := exp[3]
				condResult := false

				switch condition := condition.(type) {
				case bool:
					if condition {
						condResult = true
					}

				case types.Number:
					if condition != 0 {
						condResult = true
					}

				case []types.Expression:
					if len(condition) > 0 {
						condResult = true
					}

				default:
					if condition != nil {
						condResult = true
					}
				}

				if condResult {
					result, err = Eval(consequence, env)
				} else {
					result, err = Eval(alternative, env)
				}

			case "lambda":
				if len(exp) != 3 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'lambda'")
					break
				}
				result = types.Function{exp[1], exp[2], env}

			case "define":
				if len(exp) != 3 {
					err = errors.New("Syntax Error: Wrong number of arguments to 'define'")
					break
				}

				key, ok := exp[1].(types.Symbol)
				if !ok {
					err = errors.New("Syntax Error: Cannot assign to a literal")
					break
				}

				var value types.Expression
				value, err = Eval(exp[2], env)

				if err != nil {
					break
				}

				env.Symbols[key] = value

			// otherwise its a function
			default:
				result, err = applyFn(exp, env)
			}

		// otherwise its *probably* a function literal
		default:
			result, err = applyFn(exp, env)
		}

	// constant literal
	default:
		result = exp
	}

	return result, err
}

func getSymbol(symbol types.Symbol, env *types.Scope) (*types.Scope, error) {

	// get the symbol value if its there
	if _, ok := env.Symbols[symbol]; ok {
		return env, nil
	}

	// otherwise check the next scope
	if env.Outer != nil {
		return getSymbol(symbol, env.Outer)
	}

	// otherwise the symbol is not found
	return nil, errors.New("'" + string(symbol) + "' is not defined")
}

func applyFn(fn []types.Expression, env *types.Scope) (result types.Expression, err error) {

	args := fn[1:]
	evaluated_args := make([]types.Expression, len(args))

	// get the function body
	var body types.Expression
	body, err = Eval(fn[0], env)

	// evaluate the arguments
	for i, op := range args {
		evaluated_args[i], err = Eval(op, env)
		if err != nil {
			return nil, err
		}
	}

	// call the function
	switch f := body.(type) {

	// built in functions
	case func(...types.Expression) (types.Expression, error):
		result, err = f(evaluated_args...)

	case func(...types.Expression) types.Expression:
		result = f(evaluated_args...)

	// user defined functions
	case types.Function:

		// make new environment with outer scope
		scope := &types.Scope{make(map[types.Symbol]types.Expression), f.Env}

		switch params := f.Params.(type) {

		case []types.Expression:
			if len(params) != len(evaluated_args) {
				err = errors.New(fmt.Sprintf("Wrong number of arguments to function. Expecting %s, got %s", len(params), len(args)))
			}

			for i, key := range params {
				scope.Symbols[key.(types.Symbol)] = evaluated_args[i]
			}

		default:
			if len(evaluated_args) != 1 {
				err = errors.New(fmt.Sprintf("Wrong number of arguments to function. Expecting 1, got %s", len(evaluated_args)))
			}
			scope.Symbols[params.(types.Symbol)] = evaluated_args[0]
		}

		result, err = Eval(f.Body, scope)

	default:
		err = errors.New(fmt.Sprintf("%s is not callable", fn))
	}

	return result, err
}
