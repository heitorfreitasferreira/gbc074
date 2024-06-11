package dfa

import (
	"fmt"
	"os"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/utils"
)

type StateMachine struct {
	Qs        []types.State
	Sigma     []types.Transition
	Q         types.State
	Fs        []types.State
	Delta     map[types.State]map[types.Transition]types.State
	StateName map[types.State]string
	Stack     *utils.Stack
}

func (sm *StateMachine) Step(t types.Transition) error {
	if t == 'q' {
		fmt.Printf("Quitting...\n")
		os.Exit(0)
	}
	transitions, ok := sm.Delta[sm.Q]
	if !ok {
		stateList := ""
		for _, state := range sm.Qs {
			stateList += fmt.Sprintf("%d ", state)
		}
		panic(fmt.Errorf("state [%d] not found, avaliable states: [%s]", sm.Q, stateList))
	}
	newState, ok := transitions[t]
	if !ok {
		transitionList := ""
		for _, transition := range sm.Sigma {
			transitionList += fmt.Sprintf("%c ", transition)
		}
		return fmt.Errorf("transition [%c] not found, avaliable transitions: [%s]", t, transitionList)
	}
	sm.Q = newState

	switch t {
	case 'b':
		sm.Stack.Pop()
	case 0:
		sm.Stack.Clear()
	default:
		sm.Stack.Push(sm.Q)
	}

	return nil
}

func (sm *StateMachine) AvaliableTransitions() []types.Transition {
	transitions, ok := sm.Delta[sm.Q]
	if !ok {
		stateList := ""
		for _, state := range sm.Qs {
			stateList += fmt.Sprintf("%d ", state)
		}
		panic(fmt.Errorf("state [%d] not found, avaliable states: [%s]", sm.Q, stateList))
	}
	var avaliableTransitions []types.Transition
	for transition := range transitions {
		avaliableTransitions = append(avaliableTransitions, transition)
	}
	return avaliableTransitions
}

func (sm StateMachine) IsInFinalState() bool {
	for _, s := range sm.Fs {
		if sm.Q == s {
			return true
		}
	}
	return false
}

func NewViewPicker() *StateMachine {
	return &StateMachine{
		Qs: []types.State{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		StateName: map[types.State]string{
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
			11: "by cpf",
			12: "all",
			13: "by isbn",
			14: "all",
			15: "quit",
		},
		Sigma: []types.Transition{'c', 'r', 'u', 'd', 'a', '1', 'q', 'b', '0'},
		Q:     0,
		Fs:    []types.State{3, 11, 12, 5, 6, 7, 13, 14, 9, 10},
		Delta: map[types.State]map[types.Transition]types.State{
			0: {
				'u': 1,
				'l': 2,
				'q': 15,
				'b': 0,
			},
			1: {
				'c': 3,
				'r': 4,
				'u': 5,
				'd': 6,
				'q': 15,
				'b': 0,
			},
			2: {
				'c': 7,
				'r': 8,
				'u': 9,
				'd': 10,
				'b': 0,
			},
			3: {
				'b': 1,
			},
			4: {
				'1': 11,
				'a': 12,
				'b': 1,
			},
			5: {
				'b': 1,
			},
			6: {
				'b': 1,
			},
			7: {
				'b': 2,
			},
			8: {
				'1': 13,
				'a': 14,
				'b': 2,
			},
			9: {
				'b': 2,
			},
			10: {
				'b': 2,
			},
			11: {
				'b': 4,
			},
			12: {
				'b': 4,
			},
			13: {
				'b': 8,
			},
			14: {
				'b': 8,
			},
		},
		Stack: utils.NewStack(),
	}
}
