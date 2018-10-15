package main

import (
    "fmt"
    "net"
    "net/http"
    "time"
    "bufio"
    "os"
    "../Config"
    "log"
    //"os"
    "github.com/golang/protobuf/proto"
    "../Cmd"
)

// func readMessage(conn net.Conn) {
//     defer conn.Close()
//         buf := make([]byte, 4096, 4096)
//         for {
//             cnt, err := conn.Read(buf)
//                 if err != nil {
//                     panic(err)
//                 }
//             stReceive := &Cmd.LoginCmd{}
//             pData := buf[:cnt]
//             // log.Print(type(pData))
// 
//             err = proto.Unmarshal(pData, stReceive)
//             if err != nil {
//                 panic(err)
//             }
// 
//             fmt.Println("receive:", conn.RemoteAddr())//, stReceive)
//             if *stReceive.Message == "stop" {
//                 os.Exit(1)
//             }
//         }
// }

type timeHandler struct {
    format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tm := time.Now().Format(th.format)
    w.Write([]byte("The time is: " + tm))
}

func Handler4Client(w http.ResponseWriter, r * http.Request){
    fmt.Fprintln(w, "hello world")

    handler, find := s_msg2HandlerMap[100]
    if find == true{
        buf := make([]byte, 4, 4)
        handler(buf)
    }else{
        log.Println("can not find msg handler")
    }
}

func main() {
    // server client => Login
    config := Config.GetInstance()
    Init()
    // RegisterMsgHandler(100, handler)
    // handler, find := s_msg2HandlerMap[100]
    // if find == true{
    //     handler(1, 2)
    // }else{
    //     log.Println("can not find msg handler")
    // }
    // go func() {
    //     //listener, err := net.Listen("tcp", "localhost:6603")
    //     listener, err := net.Listen("tcp",config.GetLogin4ClientAddr())
    //         if err != nil {
    //             panic(err)
    //         }
    //     for {
    //         conn, err := listener.Accept()
    //             if err != nil {
    //                 panic(err)
    //             }
    //         fmt.Println("new client connect Login:", conn.RemoteAddr())
    //             go readMessage(conn)
    //     }

    // }()
    go func(){
        http.HandleFunc("/handler4Client/", Handler4Client)
        http.ListenAndServe(config.GetLogin4ClientAddr(), nil)
    }()
    //go func(){
    //    mux := http.NewServeMux()

    //    th := &timeHandler{format: time.RFC1123}
    //    mux.Handle("/time", th)
    //     //mux.HandleFunc("/time", timeHandler)

    //     //log.Println("Listening...")
    //     http.ListenAndServe("127.0.0.1:18080", mux)
    //}()
    
    // client Login => Gateway
    go func(){
        //strIP := "localhost:6600"
        strIP := config.GetLogin2GatewayAddr()
        var conn net.Conn
        var err error
        for conn, err = net.Dial("tcp", strIP); err != nil; conn, err = net.Dial("tcp", strIP) {
            fmt.Println("connect", strIP, "fail")
            time.Sleep(time.Second)
            fmt.Println("reconnect...")
        }
        fmt.Println("connect", strIP, "success")
        defer conn.Close()
        cnt := 0
        sender := bufio.NewScanner(os.Stdin)
        inputStr := sender.Text()
        for sender.Scan() {
            cnt++
            stSend := &Cmd.LoginCmd{
                Msg: &inputStr,
                Length:  proto.Int(len(sender.Text())),
                Cnt:     proto.Int(cnt),
            }
            //var t string
            // stSend.AppendToString(t)
            pData, err := proto.Marshal(stSend)
            if err != nil {
                panic(err)
            }

            conn.Write(pData)
            if sender.Text() == "stop" {
                return
            }
        }
    }()

    for{}

}
