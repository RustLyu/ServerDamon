package main

import (
    "log"
    //"github.com/golang/protobuf/descriptor"
    "../util"
)

var s_msg2HandlerMap = make(map[int]util.MsgHandler)

func RegisterMsgHandler(msgType int, handler util.MsgHandler){
    mm, ok := s_msg2HandlerMap[msgType]
    if !ok {
        s_msg2HandlerMap[msgType] = handler
        log.Println("RegisterMsgHandler ok")
    } else {
        log.Println("duplicate register msg handler")
        log.Println(mm)
    }
}

func handler(buf []byte){
    log.Println("handler test")
}

func Init(){
    RegisterMsgHandler(100, handler)
   // hand, find := s_msg2HandlerMap["666"]
   // if find == true{
   //     hand(4, 5)
   // }else{
   //     log.Println("duplicate register msg handler 11111111")
   // }
}
