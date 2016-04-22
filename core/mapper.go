package mantle

import (
	"encoding/json"
	"fmt"

	"github.com/exis-io/core"
)

// Reflects the core, handles invocations, and manages memory

// All connected sessions. Note-- moving this to the server
// var sessions = make(map[uint64]*session)

// func OpenSession() uint64 {
// 	s := &session{make(map[uint64]interface{})}
// 	id := core.NewID()
// 	sessions[id] = s

// 	return id
// }

// func CloseSession(id uint64) {
// 	delete(sessions, id)
// }

// type session struct {
// 	memory map[uint64]interface{} // "heap" space for this session
// }

var memory = make(map[uint64]interface{})

type invocation struct {
	target string        // A type, function, variable, or constant
	cb     float64       // The callback id to deliver the result on
	eb     float64       // errback id to deliver failure on
	args   []interface{} // Arguments to pass to the target
}

func Handle(line []byte) {
	n, err := deserialize(line)

	if err != nil {
		fmt.Printf("Ignoring message %b. Error: %s", line, err.Error())
	}

	fmt.Println("Received: ", n)

	if m, ok := core.Variables[n.target]; ok {
		fmt.Println("Variable", m)
	} else if m, ok := core.Types[n.target]; ok {
		fmt.Println("Types", m)
	} else if m, ok := core.Consts[n.target]; ok {
		fmt.Println("Consts", m)
	} else if m, ok := core.Functions[n.target]; ok {
		fmt.Printf("Functions %v\n", m)
		// r := m.Call(n.args)
		// fmt.Printf("Result: %v\n", r)
	} else {
		fmt.Println("Unknown!")
	}
}

func deserialize(j []byte) (*invocation, error) {
	var d []interface{}
	if e := json.Unmarshal(j, &d); e != nil {
		return nil, fmt.Errorf("Unable to unmarshall data: %s\n", e)
	}

	n := &invocation{}

	if s, ok := d[0].(string); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 0")
	} else {
		n.target = s
	}

	if s, ok := d[1].(float64); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 1")
	} else {
		n.cb = s
	}

	if s, ok := d[2].(float64); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 2")
	} else {
		n.eb = s
	}

	if s, ok := d[3].([]interface{}); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 3")
	} else {
		n.args = s
	}

	return n, nil
}
