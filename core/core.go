// This is the "server"
// It should respond to the Java "client" as if it were native code.
// This means the go could be compiled as a local library, or served
// remotely. This should manage memory and objects
package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listener, _ := net.Listen("tcp", ":9876")
	fmt.Println("Server started. Reflected package interface: ")

	for {
		conn, _ := listener.Accept()
		fmt.Println("Someone connected")
		writer := bufio.NewWriter(conn)
		reader := bufio.NewReader(conn)

		go func() {
			line, _ := reader.ReadString('\n')
			fmt.Println("Received: ", line)
			writer.WriteString("Go: Hey client!\n")
			writer.Flush()
		}()
	}
}
package core

// Provides a generalized interface into all Core functions
// Relies on github.com/damouse/pgkreflect to dynamically check calls and types

type session struct {
    mem map[uint64]interface{} // "heap" space for this session
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

// Maps a JSON to a struct matched by a key?
// DeferredReturn: uint64 pointer
func (s *session) Alloc(key string, args []interface{}) {

}

func (s *session) Free(ptr uint64) {

}

// Wraps a function in a deferred. If the function returns a non-nil error, callback eb
// with the string resut. If the function returns anything else then call it back
func (s *session) WrapDeferred(cb uint64, eb uint64) {

}
