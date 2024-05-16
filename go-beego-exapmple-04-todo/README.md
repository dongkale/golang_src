# BeegoTodo

## Steps [Only API]

- Get Bee Cli Tool
- Generate Project Using `bee api [appname]`
- Generate Model using `bee generate model [modelname] [-fields="name:type"]`
- Generate Controller `bee generate controller [controllerfile]
`
- Add Controller To `routers/router.go`

```go=
package routers

import (
	"github.com/nafeem-evatix/beegotodo/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/todo",
			beego.NSInclude(
				&controllers.TodoController{},
			),
		),
	)

	beego.AddNamespace(ns)
}
```

- Run `bee run -downdoc=true -gendoc=true` this will run and generate and serve API Documentation
- Go to `localhost:8080/swagger` use to do CRUD Operations

![](https://i.imgur.com/ghej4OQ.png)

## Document

- https://hackmd.io/@haquenafeem/ryFHIFXKY
- https://github.com/nafeem-evatix/beegotodo
