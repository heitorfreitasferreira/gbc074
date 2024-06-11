package view

import (
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/dfa.go"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/internal/handlers"
)

var finalStateView map[types.State]func() = map[types.State]func(){
	types.State(3):  handlers.NewCreateUserUseCase,
	types.State(11): nil,
	types.State(12): nil,
	types.State(5):  nil,
	types.State(6):  nil,
	types.State(7):  nil,
	types.State(9):  nil,
	types.State(10): nil,
}

func View(sm *dfa.StateMachine) {
	s := sm.Q
	fmt.Print("\033[H\033[2J") // Clear the console screen
	if sm.IsInFinalState() {
		finalStateView[s]()
		err := sm.Step('b')
		if err != nil {
			panic("can't return back after use-case view")
		}
		fmt.Println("Press Enter to continue...")
		fmt.Scanln()
		View(sm)

	} else {
		fmt.Printf("%s\n", sm.Stack.ToString("->", sm.StateName))
		avTransitions := ""
		for _, tr := range sm.AvaliableTransitions() {
			avTransitions += fmt.Sprintf("%c ", tr)
		}

		fmt.Printf("Avaliable transitions: %s\n", avTransitions)
	}
}
