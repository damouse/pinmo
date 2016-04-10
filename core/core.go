// This is the "server"
// It should respond to the Java "client" as if it were native code.
// This means the go could be compiled as a local library, or served
// remotely. This should manage memory and objects
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"reflect"

	"github.com/exis-io/core"
)

func main() {
	listener, _ := net.Listen("tcp", ":9876")
	fmt.Println("core started")

	methods := make(map[string]map[string]reflect.Method)

	// Build the list of methods available
	for k, v := range core.Types {
		e := make(map[string]reflect.Method)

		for i := 0; i < v.NumMethod(); i++ {
			m := v.Method(i)
			e[m.Name] = m
		}

		methods[k] = e
	}

	var sessions []*session

	for {
		conn, _ := listener.Accept()
		writer := bufio.NewWriter(conn)
		reader := bufio.NewReader(conn)

		s := &session{
			mem:    make(map[uint64]interface{}),
			conn:   &conn,
			writer: writer,
			reader: reader,
		}

		sessions = append(sessions, s)
		go s.run()
	}
}

// Provides a generalized interface into all Core functions
// Relies on github.com/damouse/pgkreflect to dynamically check calls and types

type invocation struct {
	target string
	cb     float64
	eb     float64
	args   []interface{}
}

type session struct {
	mem    map[uint64]interface{} // "heap" space for this session
	conn   *net.Conn
	writer *bufio.Writer
	reader *bufio.Reader
}

func (s *session) run() {
	for {
		line, _ := s.reader.ReadBytes('\n')

		var d []interface{}
		if e := json.Unmarshal(line, &d); e != nil {
			fmt.Printf("Unable to unmarshall data: %s\n", e)
		}

		n := &invocation{
			target: d[0].(string),
			cb:     d[1].(float64),
			eb:     d[2].(float64),
			args:   d[3].([]interface{}),
		}

		fmt.Println("Received: ", n)
		// s.writer.WriteString("Go: Hey client!\n")
		// s.writer.Flush()

		if m, ok := core.Variables[n.target]; ok {
			fmt.Println("Variable", m)
		} else if m, ok := core.Types[n.target]; ok {
			fmt.Println("Types", m)
		} else if m, ok := core.Consts[n.target]; ok {
			fmt.Println("Consts", m)
		} else if m, ok := core.Functions[n.target]; ok {
			fmt.Printf("Functions %v\n", m)
			r := m.Call(n.args)
			fmt.Printf("Result: %v\n", r)
		} else {
			fmt.Println("Unknown!")
		}
	}
}

// Setting variables
func (s *session) Set(key string, value []interface{}) {

}

func (s *session) Get(key string) {

}

// Invoke the given function with args and callback ids. If the function can and does
// succeed, cb is called back with the results. If the function returns a non-nil error,
// eb is invoked with a string
func (s *session) InvokeDeferred(name string, cb int64, eb int64, args []interface{}) {

}

func (s *session) Invoke(name string, args []interface{}) {

}

// Maps a JSON to a struct matched by a key?
// DeferredReturn: uint64 pointer
func (s *session) Alloc(key string, args []interface{}) {

}

func (s *session) Free(ptr uint64) {

}
