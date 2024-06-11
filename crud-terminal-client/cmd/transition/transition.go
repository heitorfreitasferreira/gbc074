package transition

import (
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/dfa.go"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
)

func Transition(sm *dfa.StateMachine) {
	var attempthTransition types.Transition

	var err error

	fmt.Scanf("%c \n", &attempthTransition)
	err = sm.Step(attempthTransition)
	if err != nil {
		fmt.Printf("Invalid Transition [%c]\n", attempthTransition)
		Transition(sm)
	}
}
