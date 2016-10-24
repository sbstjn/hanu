package allot

import (
	"errors"
	"regexp"
	"strings"
)

// CommandInterface is the interface
type CommandInterface interface {
	Expression() *regexp.Regexp
	Has(name ParameterInterface) bool
	Match(req string) (MatchInterface, error)
	Matches(req string) bool
	Parameters() []Parameter
	Position(param ParameterInterface) int
	Text() string
}

// Command is the struct
type Command struct {
	text string
}

// Text returns the command text
func (c Command) Text() string {
	return c.text
}

// Expression returns the regular expression
func (c Command) Expression() *regexp.Regexp {
	expr := c.Text()

	for _, param := range c.Parameters() {
		expr = strings.Replace(expr, "<"+param.Name()+":"+param.Data()+">", "("+param.Expression().String()+")", -1)
		expr = strings.Replace(expr, "<"+param.Name()+">", "("+param.Expression().String()+")", -1)
	}

	return regexp.MustCompile("^" + expr + "$")
}

// Parameters returns the list of defined parameters
func (c Command) Parameters() []Parameter {
	var list []Parameter
	re := regexp.MustCompile("<(.*?)>")
	result := re.FindAllStringSubmatch(c.Text(), -1)

	for _, p := range result {
		if len(p) != 2 {
			continue
		}

		pType := ""
		if !strings.Contains(p[1], ":") {
			pType = ":string"
		}

		list = append(list, Parse(p[1]+pType))
	}

	return list
}

// Has checks if the parameter is found in the command
func (c Command) Has(param ParameterInterface) bool {
	return c.Position(param) != -1
}

// Position returns the position of a parameter
func (c Command) Position(param ParameterInterface) int {
	for index, item := range c.Parameters() {
		if item.Equals(param) {
			return index
		}
	}

	return -1
}

// Match returns matches command
func (c Command) Match(req string) (MatchInterface, error) {
	if c.Matches(req) {
		return Match{c, req}, nil
	}

	return nil, errors.New("Request does not match Command.")
}

// Matches checks if a comand definition matches a request
func (c Command) Matches(req string) bool {
	// fmt.Println(c.Text(), "\n => "+c.Expression().String())
	return c.Expression().MatchString(req)
}

// New returns a new command
func New(command string) Command {
	return Command{command}
}
