package evaluator

import (
	"math"
	"monkey/object"
)

var (
	E       = &object.Float{Value: math.E}
	PI      = &object.Float{Value: math.Pi}
	PHI     = &object.Float{Value: math.Phi}
	SQRT2   = &object.Float{Value: math.Sqrt2}
	SQRTE   = &object.Float{Value: math.SqrtE}
	SQRTPI  = &object.Float{Value: math.SqrtPi}
	SQRTPHI = &object.Float{Value: math.SqrtPhi}
	LN2     = &object.Float{Value: math.Ln2}
	LOG2E   = &object.Float{Value: math.Log2E}
	LN10    = &object.Float{Value: math.Ln10}
	LOG10E  = &object.Float{Value: math.Log10E}
)

func init() {
	superGlobals["E"] = E

	superGlobals["PI"] = PI

	superGlobals["PHI"] = PHI

	superGlobals["SQRT2"] = SQRT2

	superGlobals["SQRT2"] = SQRT2

	superGlobals["SQRTE"] = SQRTE

	superGlobals["SQRTPI"] = SQRTPI

	superGlobals["SQRTPHI"] = SQRTPHI

	superGlobals["LN2"] = LN2

	superGlobals["LOG2E"] = LOG2E

	superGlobals["LN10"] = LN10

	superGlobals["LOG10E"] = LOG10E
}
