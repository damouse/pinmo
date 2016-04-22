package mantle

import (
	"encoding/json"
	"fmt"
	"reflect"

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
	target string      // A type, function, variable, or constant
	cb     uint64      // The callback id to deliver the result on
	eb     uint64      // errback id to deliver failure on
	args   interface{} // Arguments to pass to the target
}

func Handle(line string) {
	n, err := deserialize(line)

	if err != nil {
		fmt.Printf("Ignoring message %b. Error: %s\n", line, err.Error())
		return
	}

	fmt.Println("Received: ", n)

	var result interface{}
	var resultingId = n.cb

	if m, ok := core.Variables[n.target]; ok {
		result = handleVariable(m, n.args)
	} else if m, ok := core.Types[n.target]; ok {
		fmt.Printf("Types %v\n", m)
	} else if m, ok := core.Consts[n.target]; ok {
		result = m.Elem()
	} else if m, ok := core.Functions[n.target]; ok {
		fmt.Printf("Functions %v\n", m)
		// r := m.Call(n.args)
		// fmt.Printf("Result: %v\n", r)
	} else {
		fmt.Println("Unknown!")
	}

	dispatch(resultingId, result)
}

// Assign the given value to a variable and return its value. If we are passed "nil" as a
// new value this is just a read-- dont try and set the value. Obviously this means nil is
// not allowed as a variable value.
func handleVariable(v reflect.Value, n interface{}) interface{} {
	// Note: this will panic on bad type conversions, please check!
	if n != nil {
		c := reflect.ValueOf(n).Convert(v.Elem().Type())
		v.Elem().Set(c)
	}

	return v.Elem()
}

func handleType() (interface{}, error) {
	return nil, nil
}

func handleFunction() (interface{}, error) {
	return nil, nil
}

func handleObject() (interface{}, error) {
	return nil, nil
}

// Dispatch a callback to the given session
func dispatch(id uint64, arg interface{}) {
	fmt.Printf("Dispatching %d %v\n", id, arg)
}

func deserialize(j string) (*invocation, error) {
	var d []interface{}
	if e := json.Unmarshal([]byte(j), &d); e != nil {
		return nil, fmt.Errorf("Unable to unmarshall data: %s\n", e)
	}

	n := &invocation{}

	if s, ok := d[0].(string); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 0. Got %v", d[0])
	} else {
		n.target = s
	}

	if s, ok := d[1].(float64); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 1. Got %v", d[1])
	} else {
		n.cb = uint64(s)
	}

	if s, ok := d[2].(float64); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 2. Got %v", d[2])
	} else {
		n.eb = uint64(s)
	}

	if s, ok := d[3].(interface{}); !ok {
		if d[3] == nil {
			n.args = nil
		} else {
			return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 3. Got %v", d[3])
		}
	} else {
		n.args = s
	}

	return n, nil
}
