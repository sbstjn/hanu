package allot

import (
	"errors"
	"fmt"
	"strconv"
)

// MatchInterface is the interface
type MatchInterface interface {
	String(name string) (string, error)
	Integer(name string) (int, error)
	Match(position int) (string, error)

	Parameter(param ParameterInterface) (string, error)
}

// Match is the struct
type Match struct {
	Command CommandInterface
	Request string
}

// String returns the value for a string parameter
func (m Match) String(name string) (string, error) {
	return m.Parameter(NewParameterWithType(name, "string"))
}

// Integer returns the value for an integer parameter
func (m Match) Integer(name string) (int, error) {
	str, err := m.Parameter(NewParameterWithType(name, "integer"))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

// Parameter returns the value for a parameter
func (m Match) Parameter(param ParameterInterface) (string, error) {
	pos := m.Command.Position(param)
	if pos == -1 {
		return "", errors.New("Unknonw parameter \"" + param.Name() + "\"")
	}

	matches := m.Command.Expression().FindAllStringSubmatch(m.Request, -1)[0][1:]
	return matches[m.Command.Position(param)], nil
}

// Match returns the match at given position
func (m Match) Match(position int) (string, error) {
	matches := m.Command.Expression().FindAllStringSubmatch(m.Request, -1)

	if len(matches) != 1 {
		return "", errors.New("Unable to parse request")
	}

	if position >= len(matches[0]) {
		return "", fmt.Errorf("No parameter at position %d", position)
	}

	return matches[0][position+1], nil
}
