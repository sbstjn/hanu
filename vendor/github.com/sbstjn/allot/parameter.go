package allot

import (
	"regexp"
	"strings"
)

var regexpMapping = map[string]string{
	"string":  "[^\\s]+",
	"integer": "[0-9]+",
}

// Expression returns the regexp for a data type
func Expression(data string) *regexp.Regexp {
	if exp, ok := regexpMapping[data]; ok {
		return regexp.MustCompile(exp)
	}

	return nil
}

// ParameterInterface is the interface
type ParameterInterface interface {
	Equals(param ParameterInterface) bool
	Expression() *regexp.Regexp
	Name() string
	Data() string
}

// Parameter is the struct
type Parameter struct {
	name string
	data string
	expr *regexp.Regexp
}

// Expression returns the regexp behind the type
func (p Parameter) Expression() *regexp.Regexp {
	return p.expr
}

// Name returns the Parameter name
func (p Parameter) Name() string {
	return p.name
}

// Data returns the Parameter name
func (p Parameter) Data() string {
	return p.data
}

// Equals checks if two parameter are equal
func (p Parameter) Equals(param ParameterInterface) bool {
	return p.Name() == param.Name() && p.Expression().String() == param.Expression().String()
}

// NewParameterWithType returns
func NewParameterWithType(name string, data string) Parameter {
	return Parameter{name, data, Expression(data)}
}

// Parse parses parameter info
func Parse(text string) Parameter {
	var splits []string
	var name, data string

	name = strings.Replace(text, "<", "", -1)
	name = strings.Replace(name, ">", "", -1)
	data = "string"

	if strings.Contains(name, ":") {
		splits = strings.Split(name, ":")

		name = splits[0]
		data = splits[1]
	}

	return NewParameterWithType(name, data)
}
