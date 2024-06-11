package stack

import "github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"

// A pilha do automato para visualização de onde está no grafo até agora
var STACK []types.State = []types.State{}
