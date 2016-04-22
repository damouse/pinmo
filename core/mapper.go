package mantle

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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
	target  string      // A type, function, variable, or constant
	cb      uint64      // The callback id to deliver the result on
	eb      uint64      // errback id to deliver failure on
	address uint64      // multiple-use field: "Pointer" when dealing with types, handler id for domain operations
	args    interface{} // Arguments to pass to the target
}

func Handle(line string) {
	n, err := deserialize(line)

	if err != nil {
		fmt.Printf("Ignoring message %b. Error: %s\n", line, err.Error())
		return
	}

	var result interface{}
	var resultingId = n.cb

	// Note: Types are used only for the sake of method reading-- instantiate types through their
	// public constructors!

	if m, ok := core.Variables[n.target]; ok {
		result = handleVariable(m, n.args)
	} else if m, ok := core.Consts[n.target]; ok {
		result = m.Interface()
	} else if m, ok := core.Functions[n.target]; ok {
		if ret, err := handleFunction(m, n.args); err != nil {
			resultingId = n.eb
			result = err.Error()
		} else {
			if handleConstructor(n.target, n.address, ret) {
				result = n.address
			}

			result = ret
			resultingId = n.cb
		}
	} else if m, ok := memory[n.address]; ok {
		v := reflect.ValueOf(m).MethodByName(n.target)

		var err error
		if result, err = handleFunction(v, n.args); err != nil {
			resultingId = n.eb
			result = err.Error()
		} else {
			resultingId = n.cb
		}
	} else {
		fmt.Println("Unknown!")
	}

	dispatch(resultingId, result)
}

// Assign the given value to a variable and return its value. If we are passed "nil" as a
// new value this is just a read-- dont try and set the value. Obviously this means nil is
// not allowed as a variable value.
// TODO: handle bad type conversions
func handleVariable(v reflect.Value, n interface{}) interface{} {
	if n != nil {
		c := reflect.ValueOf(n).Convert(v.Elem().Type())
		v.Elem().Set(c)
	}

	return v.Elem()
}

func handleFunction(fn reflect.Value, args interface{}) ([]interface{}, error) {
	if argsList, ok := args.([]interface{}); !ok {
		return nil, fmt.Errorf("Function invocations require a list of arguments!")
	} else {
		return core.Cumin(fn.Interface(), argsList)
	}
}

// Checks to see if a function invocation instantiated an object by checking the string of the target.
// By convention constructors must be named "New[TypeName]" and return pointers.
// If found and memory has been allocated for the given pointer, return true
func handleConstructor(target string, address uint64, invocationResult []interface{}) bool {
	if len(invocationResult) != 1 {
		return false
	}

	if strings.Index(target, "New") != -1 {
		split := strings.Split(target, "New")

		if len(split) == 2 && split[0] == "" {
			memory[address] = invocationResult[0]
			return true
		}
	}

	return false
}

func handleObject(t reflect.Value) (interface{}, error) {
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

	if s, ok := d[3].(float64); !ok {
		return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 3. Got %v", d[2])
	} else {
		n.address = uint64(s)
	}

	if s, ok := d[4].(interface{}); !ok {
		if d[4] == nil {
			n.args = nil
		} else {
			return nil, fmt.Errorf("Couldn't parse message-- incorrect type at position 4. Got %v", d[3])
		}
	} else {
		n.args = s
	}

	return n, nil
}
