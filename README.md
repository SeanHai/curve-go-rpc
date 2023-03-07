## curve-go-rpc
curve-go-rpc是一个与curve交互的golang实现的 rpc client 功能集。

## usage

`git submodule add https://github.com/SeanHai/curve-go-rpc destPath`

`git submodule update --remote`

```
import bsrpc github.com/SeanHai/curve-go-rpc

var mdsClient bsrpc.NewMdsClient = bsrpc.NewMdsClient(bsrpc.MdsClientOption{
	TimeoutMs: 500,
	RetryTimes: 3,
	Addrs: []string{"127.0.0.1:6666","127.0.0.2:6666","127.0.0.3:6666"},
})

func ListPhysicalPool() (interface{}, error) {
    return mdsClient.ListPhysicalPool()
}
```
