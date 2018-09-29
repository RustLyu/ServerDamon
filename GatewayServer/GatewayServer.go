package main

import (
    //"bufio"
    "fmt"
    "net"
    "../Config"
    //"os"
    //"time"
    //"../Cmd"
    //"github.com/golang/protobuf/proto"
)

func readMessage(conn net.Conn) {
    defer conn.Close()
        //buf := make([]byte, 4096, 4096)
        for {
            //cnt, err := conn.Read(buf)
            //    if err != nil {
            //        panic(err)
            //    }
            // stReceive := &Cmd.LoginCmd{}
            // pData := buf[:cnt]

            // err = proto.Unmarshal(pData, stReceive)
            // if err != nil {
            //     panic(err)
            // }

            //fmt.Println("receive:", conn.RemoteAddr())//, stReceive)
                // if *stReceive.Message == "stop" {
                //     os.Exit(1)
                //     }
                // }
    }
}

func main() {
    // server Login => Gateway
    config := Config.GetInstance()
    go func() {
        listener, err := net.Listen("tcp", config.GetGateway4LoginAddr())
            if err != nil {
                panic(err)
            }
        for {
            conn, err := listener.Accept()
                if err != nil {
                    panic(err)
                }
            fmt.Println("new Login connect to Gateway success:", conn.RemoteAddr())
                go readMessage(conn)
        }

    }()
    //  server GameServer => Gateway
    go func(){
        listener, err := net.Listen("tcp", config.GetGateway4GameAddr())
            if err != nil {
                panic(err)
            }
        for {
            conn, err := listener.Accept()
                if err != nil {
                    panic(err)
                }
            fmt.Println("new connect", conn.RemoteAddr())
                go readMessage(conn)
        }
    }()
    for {}
}
