package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"],
        beego.ControllerComments{
            Method: "GetTipos",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"],
        beego.ControllerComments{
            Method: "GetByTipo",
            Router: "/:tipo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"],
        beego.ControllerComments{
            Method: "GetByTipoAndID",
            Router: "/:tipo/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
