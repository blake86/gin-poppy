# gin-poppy

gin-poppy 是为gin框架提供的一个统计中间件，可以非常方便的嵌入到gin框架中实现对各个uri的请求数量、http返回码、处理延时的记录统计


##  使用方法

### Import
```go
import "git.sogou-inc.com/lihao/gin-poppy"
```

### 典型用例

[embedmd]:# (example/main.go go)
```go
package main

import (
	"github.com/blake86/gin-poppy"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
    
    	gin_poppy.Register(r)
	
	//r.GET("/ping", handler)
	
	r.Run(":8080")
}
```

### 查看数据统计

```
http://localhost:8999/debug/poppy/
```

返回数据格式为：
```json
{{
     "poppy":{
         "UriResult":{

             "/ping":{
                 "Result":{
                     "TotalCount":8,
                     "TotalCodeCount":null,
                     "Count":8,
                     "CodeCount":{
                         "200":8
                     },
                     "Min":10429,
                     "Max":166490,
                     "Average":30551.75,
                     "Mean":30551.75,
                     "Stdev":54931.488195361526,
                     "P90":166490,
                     "P95":166490,
                     "P99":166490,
                     "StartTime":1585639529,
                     "LastTime":1585639530
                 }
             }
         },
         "GlobalResult":{
             "StartTime":1585639529,
             "LastTime":1585641580,
             "TotalCount":8,
             "TotalCodeCount":{
                 "200":8
             }
         }
     }
 }
```