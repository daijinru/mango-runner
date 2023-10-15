# Mango CLI

## Build

```bash
$ go build
```
## Use Guider

Add it the content below to your `[project-root]/.mango/mango-ci.yaml`.
```yaml
Version: "abc"
Stages:
  - start
  - build

job-dev:
  stage: start
  scripts:
    - echo "dev success"

build-job:
  stage: build
  scripts:
    - echo "build success"
```

Execute at the command line.
```bash
$ mango-cli rpc start 1234
```

### How to call by RPC

```go
func main() {
  client, err := rpc.Dial("tcp", "localhost:1234")
  // ...
  defer client.Close()
  // ...
  var reply = &Reply{}
  // ...
  err = client.Call("CiService.Run", reqOption, reply)
  if err != nil {
    fmt.Println("call fail: ", err)
  }
}
```