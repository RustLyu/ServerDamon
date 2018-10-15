package util

import (
    "bytes"
    "encoding/binary"
    "github.com/golang/protobuf/proto"
    // "log"
)
// msg title length
var cmdLength = 4
// Msg handler func type
type MsgHandler func([]byte)

// int => byte[]
func IntToBytes(n int) []byte {
    tmp := int32(n)
    bytesBuffer := bytes.NewBuffer([]byte{})
    binary.Write(bytesBuffer, binary.BigEndian, tmp)
    return bytesBuffer.Bytes()
}

// byte[] => int
func BytesToInt(buf []byte) int {
    bytesBuffer := bytes.NewBuffer(buf)
    var tmp int32
    binary.Read(bytesBuffer, binary.BigEndian, &tmp)
    return int(tmp)
}

func Serialize(id int, pb proto.Message) []byte {
	ret, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	count := len(ret)
	buf := make([]byte, count + cmdLength)
	for i := cmdLength; i < count + cmdLength; i++{
		buf[i] = ret[i - cmdLength]
	}
	cmd := IntToBytes(id)
	for i:= 0; i < cmdLength; i++{
		buf[i] = cmd[i]
	}
	return buf
}

// get msg num
func UnSerialize(buf []byte) int {
    ret := make([]byte, cmdLength, cmdLength)
    for i := 0; i < len(ret); i++{
        ret[i] = buf[i]
    }
    return BytesToInt(ret)
}

// func main(){
//     var input = 1000
//     t_0 := IntToBytes(input)
//     for i := 0; i < len(t_0); i++{
//         log.Println(t_0[i])
//     }
// 
//     t_1 := BytesToInt(t_0)
//     log.Println(t_1)
// }
