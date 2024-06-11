package stack

import (
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
)

var stateToName = map[types.State]string{
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

func LogStack() {
	for i, s := range STACK {
		fmt.Printf("%s", stateToName[s])
		if i != len(STACK)-1 {
			fmt.Printf("->")
		} else {
			fmt.Printf("\n")
		}
	}
}
