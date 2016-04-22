package mantle

// Serve the core through JSON-RPC. Deferred for now.

// func main() {
//     listener, _ := net.Listen("tcp", ":9876")
//     fmt.Println("core started")

//     methods := make(map[string]map[string]reflect.Method)

//     // Build the list of methods available
//     for k, v := range core.Types {
//         e := make(map[string]reflect.Method)

//         for i := 0; i < v.NumMethod(); i++ {
//             m := v.Method(i)
//             e[m.Name] = m
//             fmt.Printf("%v\n", m)
//         }

//         methods[k] = e
//     }

//     return
//     var sessions []*session

//     for {
//         conn, _ := listener.Accept()
//         writer := bufio.NewWriter(conn)
//         reader := bufio.NewReader(conn)

//         s := &session{
//             mem:    make(map[uint64]interface{}),
//             conn:   &conn,
//             writer: writer,
//             reader: reader,
//         }

//         sessions = append(sessions, s)
//         go s.run()
//     }
// }

// func (s *session) run() {
//     for {
//         line, _ := s.reader.ReadBytes('\n')

//         var d []interface{}
//         if e := json.Unmarshal(line, &d); e != nil {
//             fmt.Printf("Unable to unmarshall data: %s\n", e)
//         }

//         n := &invocation{
//             target: d[0].(string),
//             cb:     d[1].(float64),
//             eb:     d[2].(float64),
//             args:   d[3].([]interface{}),
//         }

//         fmt.Println("Received: ", n)
//         // s.writer.WriteString("Go: Hey client!\n")
//         // s.writer.Flush()

//         if m, ok := core.Variables[n.target]; ok {
//             fmt.Println("Variable", m)
//         } else if m, ok := core.Types[n.target]; ok {
//             fmt.Println("Types", m)
//         } else if m, ok := core.Consts[n.target]; ok {
//             fmt.Println("Consts", m)
//         } else if m, ok := core.Functions[n.target]; ok {
//             fmt.Printf("Functions %v\n", m)
//             // r := m.Call(n.args)
//             // fmt.Printf("Result: %v\n", r)
//         } else {
//             fmt.Println("Unknown!")
//         }
//     }
// }
