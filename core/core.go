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
