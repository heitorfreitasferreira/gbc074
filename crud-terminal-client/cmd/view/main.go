package view

import (
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/stack"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/internal/handlers"
)

func View(s types.State) {

	fmt.Print("\033[H\033[2J") // Clear the console screen
	stack.LogStack()
	//--------------------------------------------------------------------------
	[]func(){
		func() { // q0
			avaliableTransitions := []rune{'u', 'l'}
			fmt.Println("Welcome to the CRUD terminal client!")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q1
			avaliableTransitions := []rune{'c', 'r', 'u', 'd', 'b'}
			fmt.Println("User state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q2
			avaliableTransitions := []rune{'c', 'r', 'u', 'd', 'b'}
			fmt.Println("Book state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q3
			handlers.NewCreateUserUseCase()
			avaliableTransitions := []rune{'b'}
			fmt.Println("Create state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q4
			avaliableTransitions := []rune{'1', 'a', 'b'}
			fmt.Println("Read state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q5
			avaliableTransitions := []rune{'b'}
			fmt.Println("Update state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q6
			avaliableTransitions := []rune{'b'}
			fmt.Println("Delete state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q7
			avaliableTransitions := []rune{'b'}
			fmt.Println("Create state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q8
			avaliableTransitions := []rune{'1', 'a', 'b'}
			fmt.Println("Read state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q9
			avaliableTransitions := []rune{'b'}
			fmt.Println("Update state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q10
			avaliableTransitions := []rune{'b'}
			fmt.Println("Delete state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q11
			avaliableTransitions := []rune{'b'}
			fmt.Println("Term state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q12
			avaliableTransitions := []rune{'b'}
			fmt.Println("All state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q13
			avaliableTransitions := []rune{'b'}
			fmt.Println("Term state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
		func() { // q14
			avaliableTransitions := []rune{'b'}
			fmt.Println("All state")
			fmt.Printf("Avaliable transitions: ")
			for _, tr := range avaliableTransitions {
				fmt.Printf("%c ", tr)
			}
			fmt.Println()
		},
	}[s]()

	fmt.Println()

}
