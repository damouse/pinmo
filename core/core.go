// This is the "server"
// It should respond to the Java "client" as if it were native code.
// This means the go could be compiled as a local library, or served
// remotely. This should manage memory and objects
package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/exis-io/core"
)

func main() {
	listener, _ := net.Listen("tcp", ":9876")
	fmt.Println("Server started. Reflected package interface: ", main.Types)

	fmt.Println(core.Types)

	for {
		conn, _ := listener.Accept()
		fmt.Println("Someone joined!")

		writer := bufio.NewWriter(conn)
		reader := bufio.NewReader(conn)

		go func() {
			line, _ := reader.ReadString('\n')

			// func Unmarshal(data []byte, v interface{}) error

			fmt.Println("Read: ", line)
			fmt.Println("Writing: Hey, client!")

			writer.WriteString("Hey, client!\n")
			writer.WriteString("null")
		}()
	}
}
