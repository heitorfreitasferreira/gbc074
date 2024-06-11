package utils

import (
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
)

type Stack struct {
	Items []types.State
}

func NewStack() *Stack {
	return &Stack{
		Items: make([]types.State, 0),
	}
}

func (s *Stack) Push(state types.State) {
	s.Items = append(s.Items, state)
}

func (s *Stack) Pop() (types.State, error) {
	if len(s.Items) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}
	state := s.Items[len(s.Items)-1]
	s.Items = s.Items[:len(s.Items)-1]
	return state, nil
}

func (s *Stack) Peek() (types.State, error) {
	if len(s.Items) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}
	return s.Items[len(s.Items)-1], nil
}

func (s *Stack) Clear() {
	s.Items = make([]types.State, 0)
}

func (s Stack) ToString(separator string, stateName map[types.State]string) string {
	str := ""
	for i, item := range s.Items {
		space := ""
		if i == len(s.Items)-1 {
			space = "."
		} else {
			space = separator
		}
		str += fmt.Sprintf("%s%s", stateName[item], space)
	}
	return str
}
