package main

import (
	"fmt"
)

type transition rune
type state int

var stack []state = []state{0}

func nextState(s state, t transition) (state, error) {
	if t == 'q' {
		panic("quitting")
	}
	if t == '0' {
		stack = []state{}
		return 0, nil
	}

	switch s {
	case 0:
		switch t {
		case 'u':
			stack = append(stack, s)
			return 1, nil
		case 'b':
			stack = stack[:len(stack)-1]
			return 2, nil
		}
	case 1:
		switch t {
		case 'c':
			stack = append(stack, s)
			return 3, nil
		case 'r':
			stack = append(stack, s)
			return 4, nil
		case 'u':
			stack = append(stack, s)
			return 5, nil
		case 'd':
			stack = append(stack, s)
			return 6, nil
		case 'b':
			stack = stack[:len(stack)-1]
			return 0, nil
		}
	case 2:
		switch t {
		case 'c':
			stack = append(stack, s)
			return 7, nil
		case 'r':
			stack = append(stack, s)
			return 8, nil
		case 'u':
			stack = append(stack, s)
			return 9, nil
		case 'd':
			stack = append(stack, s)
			return 10, nil
		case 'b':
			stack = stack[:len(stack)-1]
			return 0, nil
		}
	case 3:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 1, nil
		}
	case 4:
		switch t {
		case '1':
			stack = append(stack, s)
			return 11, nil
		case 'a':
			stack = append(stack, s)
			return 12, nil
		case 'b':
			stack = stack[:len(stack)-1]
			return 1, nil
		}
	case 5:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 1, nil
		}
	case 6:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 1, nil
		}
	case 7:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 2, nil
		}
	case 8:
		switch t {
		case '1':
			stack = append(stack, s)
			return 13, nil
		case 'a':
			stack = append(stack, s)
			return 14, nil
		case 'b':
			stack = stack[:len(stack)-1]
			return 2, nil
		}
	case 9:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 2, nil
		}
	case 10:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 2, nil
		}
	case 11:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 4, nil
		}
	case 12:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 4, nil
		}
	case 13:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 8, nil
		}
	case 14:
		switch t {
		case 'b':
			stack = stack[:len(stack)-1]
			return 8, nil
		}
	}
	return s, fmt.Errorf("invalid transition")
}
func readInput() transition {
	var t transition
	fmt.Printf("Choose a transition: ")
	fmt.Scan(&t)
	return t
}

func main() {
	var state state
	for {
		view(state)
		state = state.transition()
	}
}
func (state state) transition() state {
	attempthTransition := readInput()
	var err error

	auxState, err := nextState(state, attempthTransition)
	for {
		if err == nil {
			break
		}
		// Se tentou uma transição errada, imprime o erro e pede para tentar novamente
		fmt.Printf("Invalid transition: %v\n", err)
		attempthTransition = readInput()
		auxState, err = nextState(state, attempthTransition)
	}
	return auxState
}

func q0()  {}
func q1()  {}
func q2()  {}
func q3()  {}
func q4()  {}
func q5()  {}
func q6()  {}
func q7()  {}
func q8()  {}
func q9()  {}
func q10() {}
func q11() {}
func q12() {}
func q13() {}
func q14() {}

var states = []func(){
	q0,
	q1,
	q2,
	q3,
	q4,
	q5,
	q6,
	q7,
	q8,
	q9,
	q10,
	q11,
	q12,
	q13,
	q14,
}
var stateToName = map[state]string{
	0:  "start",
	1:  "user",
	2:  "book",
	3:  "create",
	4:  "read",
	5:  "update",
	6:  "delete",
	7:  "create",
	8:  "read",
	9:  "update",
	10: "delete",
	11: "term",
	12: "all",
	13: "term",
	14: "all",
}

func view(s state) {
	logStack()
	states[s]()

	var avaliableTransitions = []rune{
		'u',
		'b',
		'q',
		'0',
		'c',
		'r',
		'd',
		'1',
		'a',
	}
	fmt.Printf("Avaliable transitions: ")
	for _, tr := range avaliableTransitions {
		fmt.Printf("%c ", tr)
	}
	fmt.Println()
}
func logStack() {
	for i, s := range stack {
		fmt.Printf("->%s  ", stateToName[s])
		if i != len(stack)-1 {
			fmt.Printf("->")
		} else {
			fmt.Printf("\n")
		}
	}
}
