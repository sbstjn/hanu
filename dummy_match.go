package hanu

import "github.com/sbstjn/allot"

type dummyMatch struct {
}

func (m dummyMatch) String(name string) (string, error) {
	return "", nil
}

func (m dummyMatch) Integer(name string) (int, error) {
	return 0, nil
}

func (m dummyMatch) Parameter(param allot.ParameterInterface) (string, error) {
	return "", nil
}

func (m dummyMatch) Match(position int) (string, error) {
	return "", nil
}
