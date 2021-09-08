package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/terceros_mid/helpers/propiedades"
)

// PropiedadesController operations for propiedades
type PropiedadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *PropiedadesController) URLMapping() {
	c.Mapping("GetDependenciaByDocumento", c.GetDependenciaByDocumento)
}

// GetDependenciaByDocumento ...
// @Title GetAll
// @Description get Dependencia with the specified {document}
// @Param	document	path 	string	true		"Tercero type available from /document/"
// @Param	typeDoc  	path 	string	true		"typeDoc of document"
// @Success 200 {object} []map[string]interface{}
// @Failure 500 Internal Error
// @Failure 501 {user} Not Implemented
// @Failure 502 Error with external API
// @router /:document/:typeDoc [get]
func (c *PropiedadesController) GetDependenciaByDocumento() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "PropiedadesController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	document := c.Ctx.Input.Param(":document")
	typeDoc := c.Ctx.Input.Param(":typeDoc")
	if dependencia, err := propiedades.GetHelperPropierdad(document, typeDoc); err == nil {
		c.Data["json"] = dependencia
	} else {
		panic(err)
	}
	c.ServeJSON()
}
