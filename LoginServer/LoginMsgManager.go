package main

import (
    //"fmt"
    "log"
)

type MsgHandler func(int, int)

type Test struct{
    a string
}

type User struct {
    Name string
    Age  int8
}

var s_msg2HandlerMap = make(map[int]MsgHandler)
//msg2HandlerMap := make(map[int]Test)

func registerMsgHandler(msgType int, handler MsgHandler){
    mm, ok := s_msg2HandlerMap[msgType]
    if !ok {
        s_msg2HandlerMap[msgType] = handler
    } else {
        log.Println("duplicate register msg handler")
        mm(1, 2)
    }
}
