package enums

import "strings"

// FunctionCode represents Modbus function codes
// FC1->Coil, FC2->Input, FC3->HoldingRegister, FC4->InputRegister
type FunctionCode int

const (
	FCCoil            FunctionCode = 1
	FCInput           FunctionCode = 2
	FCHoldingRegister FunctionCode = 3
	FCInputRegister   FunctionCode = 4
)

// Code returns the function code value
func (fc FunctionCode) Code() int {
	return int(fc)
}

// GetFunctionCodeByNameIgnoreCase returns function code from name
func GetFunctionCodeByNameIgnoreCase(name string) FunctionCode {
	switch strings.ToLower(name) {
	case "coil":
		return FCCoil
	case "input":
		return FCInput
	case "holdingregister":
		return FCHoldingRegister
	case "inputregister":
		return FCInputRegister
	default:
		return FCHoldingRegister // Default to 3
	}
}

// String returns the string representation
func (fc FunctionCode) String() string {
	switch fc {
	case FCCoil:
		return "Coil"
	case FCInput:
		return "Input"
	case FCHoldingRegister:
		return "HoldingRegister"
	case FCInputRegister:
		return "InputRegister"
	default:
		return "HoldingRegister"
	}
}
