package transition

import (
	"fmt"
	"os"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/stack"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
)

func Transition(s types.State) types.State {
	var attempthTransition types.Transition

	var err error

	var auxState types.State
	fmt.Scanf("%c \n", &attempthTransition)
	auxState, err = nextState(s, attempthTransition)
	if err == nil {
		return auxState
	}

	fmt.Printf("Invalid Transition [%c]\n", attempthTransition)
	return Transition(s)
}

func nextState(s types.State, t types.Transition) (types.State, error) {
	if t == 'q' {
		os.Exit(0)
		fmt.Printf("Quitting\n")
	}
	if t == '0' {
		stack.STACK = []types.State{}
		return 0, nil
	}

	switch s {
	case 0:
		switch t {
		case 'u':
			stack.STACK = append(stack.STACK, s)
			return 1, nil
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 2, nil
		}
	case 1:
		switch t {
		case 'c':
			stack.STACK = append(stack.STACK, s)
			return 3, nil
		case 'r':
			stack.STACK = append(stack.STACK, s)
			return 4, nil
		case 'u':
			stack.STACK = append(stack.STACK, s)
			return 5, nil
		case 'd':
			stack.STACK = append(stack.STACK, s)
			return 6, nil
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 0, nil
		}
	case 2:
		switch t {
		case 'c':
			stack.STACK = append(stack.STACK, s)
			return 7, nil
		case 'r':
			stack.STACK = append(stack.STACK, s)
			return 8, nil
		case 'u':
			stack.STACK = append(stack.STACK, s)
			return 9, nil
		case 'd':
			stack.STACK = append(stack.STACK, s)
			return 10, nil
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 0, nil
		}
	case 3:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 1, nil
		}
	case 4:
		switch t {
		case '1':
			stack.STACK = append(stack.STACK, s)
			return 11, nil
		case 'a':
			stack.STACK = append(stack.STACK, s)
			return 12, nil
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 1, nil
		}
	case 5:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 1, nil
		}
	case 6:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 1, nil
		}
	case 7:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 2, nil
		}
	case 8:
		switch t {
		case '1':
			stack.STACK = append(stack.STACK, s)
			return 13, nil
		case 'a':
			stack.STACK = append(stack.STACK, s)
			return 14, nil
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 2, nil
		}
	case 9:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 2, nil
		}
	case 10:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 2, nil
		}
	case 11:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 4, nil
		}
	case 12:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 4, nil
		}
	case 13:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 8, nil
		}
	case 14:
		switch t {
		case 'b':
			stack.STACK = stack.STACK[:len(stack.STACK)-1]
			return 8, nil
		}
	}
	return s, fmt.Errorf("invalid transition")
}
