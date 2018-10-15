package main

import (
    //"bufio"
    "fmt"
    "net"
    "../Config"
    //"os"
    //"time"
    //"../Cmd"
    "github.com/golang/protobuf/proto"
    "../util"
    "log"
)

// loginID => ip + port
var login2ConnMap = make(map[int]net.Conn)

// gameID => ip + port
var gameServer2ConnMap = make(map[int]net.Conn)

// cid => id + port
var cid2ConnMap = make(map[int]net.Conn)

func readLoginMessage(conn net.Conn) {
    defer conn.Close()
    buf := make([]byte, 4096, 4096)
    for {
        cnt, err := conn.Read(buf)
        if err != nil {
            panic(err)
        }
        log.Println(cnt)
        cmd := util.UnSerialize(buf)
        handler, find := s_msg2HandlerMap[cmd]
        if find == true{
            handler(buf)
        }
    }
}

func sendToLogin(pb proto.Message){
    conn := login2ConnMap[1]
    pData := util.Serialize(1, pb)
    conn.Write(pData)
}

func readClientMessage(conn net.Conn) {
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
            go readLoginMessage(conn)
            // first msg get cid and register if not have automake one give it to client
            login2ConnMap[1] = conn
       }

   }()
    // server Client => Gate
   go func() {
       listener, err := net.Listen("tcp", config.GetGateway4ClientAddr())
           if err != nil {
               panic(err)
           }
       for {
           conn, err := listener.Accept()
               if err != nil {
                   panic(err)
               }
           fmt.Println("new Client connect to Gateway success:", conn.RemoteAddr())
            // first msg get cid and register if not have automake one give it to client
            go readClientMessage(conn)
            cid2ConnMap[1] = conn
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
                        fmt.Println("new game server connect", conn.RemoteAddr())
                        gameServer2ConnMap[1] = conn
                        go readLoginMessage(conn)
                    }
                }()
            for {}
}
