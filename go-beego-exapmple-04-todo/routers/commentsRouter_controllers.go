package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"] = append(beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"] = append(beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"] = append(beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"] = append(beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"] = append(beego.GlobalControllerRouter["go-beego-exapmple-04-todo/controllers:TodoController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
