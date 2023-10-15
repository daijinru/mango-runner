package main
import (
  "fmt"
  "net/rpc"
)
type RunReqOption struct {
  Path string
  Tag string
  Filename string
}
type PipelineReply struct {
  Tag string
  Total int
  Filenames []string
  Running bool
  Content string
}
type Reply struct {
  Status int
  Message string
  Data PipelineReply
}
func main() {
  client, err := rpc.Dial("tcp", "localhost:1234")
  if err != nil {
    fmt.Println("无法连接到 RPC 服务:", err)
    return
  }
  defer client.Close()
  var reply = &Reply{}
  reqOption := &RunReqOption{
    Path: "/datas/mango-runner",
    Tag: "0c539af420a54f60a55f6d7a0c4be1ec",
    Filename: "0c539af420a54f60a55f6d7a0c4be1ec_20231014_010235.txt",
  }
  // err = client.Call("CiService.CreatePipeline", reqOption, reply)
  // err = client.Call("CiService.ReadPipeline", reqOption, reply)
  err = client.Call("CiService.ReadPipelineStatus", reqOption, reply)
  // err = client.Call("CiService.ReadPipelines", reqOption, reply)
  if err != nil {
    fmt.Println("call fail: ", err)
    return
  } else if reply.Status == 1 {
    fmt.Printf("ant its message: %v\n", reply.Message)
    fmt.Printf("ant its running status: %v\n", reply.Data.Running)
    fmt.Printf("ant its content: %v\n", reply.Data.Content)
    
    fmt.Printf("and its tag: %v\n", reply.Data.Tag)
    fmt.Printf("and its filenames(count %v): %v\n", reply.Data.Total, reply.Data.Filenames)
  } else {
    fmt.Printf("its message: %v", reply.Message)
  }
}