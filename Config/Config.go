package Config;

import (
    "encoding/xml"
    "io/ioutil"
    "os"
    "fmt"
)

type ConfigInfo struct {
    XMLName     xml.Name `xml:"config"`
    Version     string   `xml:"version,attr"`
    Svs         []gateWayAddr `xml:"GatewayServer"`
    LoginInfo   []loginAddr `xml:"LoginServer"`
    Description string   `xml:",innerxml"`
}


type gateWayAddr struct {
    XMLName    xml.Name `xml:"GatewayServer"`
    Gateway4LoginAddr   string   `xml:"Gateway4LoginAddr"`
    Gateway4GameAddr    string   `xml:"Gateway4GameAddr"`
    Gateway4ClientAddr    string   `xml:"Gateway4ClientAddr"`
}

type loginAddr struct {
    XMLName             xml.Name `xml:"LoginServer"`
    Login4ClientAddr    string   `xml:"Login4ClientAddr"`
    Login2GatewayAddr   string   `xml:"Login2GatewayAddr"`
}

// one instance Config handler
var Info *ConfigInfo

func GetInstance() *ConfigInfo{
    if Info == nil{
        Info = &ConfigInfo{}
        Init()
    }
    return Info
}

func Init() {
    file, err := os.Open("./config.xml") // For read access.     
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }
    defer file.Close()
    data, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }
    err = xml.Unmarshal(data, &Info)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }
}

func (info ConfigInfo) GetGateway4LoginAddr() string{
    return info.Svs[0].Gateway4LoginAddr
}

func (info ConfigInfo) GetGateway4GameAddr() string{
    return info.Svs[0].Gateway4GameAddr
}

func (info ConfigInfo) GetGateway4ClientAddr() string{
    return info.Svs[0].Gateway4ClientAddr
}

func (info ConfigInfo) GetLogin4ClientAddr() string{
    return info.LoginInfo[0].Login4ClientAddr
}

func (info ConfigInfo) GetLogin2GatewayAddr() string{
    return info.LoginInfo[0].Login2GatewayAddr
}
