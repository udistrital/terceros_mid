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
	c.Mapping("GetPropiedades", c.GetPropiedades)
	c.Mapping("GetDependenciaById", c.GetDependenciaById)
}

// GetPropiedades ...
// @Title GetAll
// @Description List the Propiedades types that can be used to gather Propiedades by {propiedad}
// @Success 200 {object} []string
// @router / [get]
func (c *PropiedadesController) GetPropiedades() {

	// Puede que ni sea necesario en este controlador, pero se coloca por lineamiento...
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			localError := err.(map[string]interface{})
			c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + "TercerosController" + "/" + (localError["funcion"]).(string))
			c.Data["data"] = (localError["err"])
			if status, ok := localError["status"]; ok {
				c.Abort(status.(string))
			} else {
				c.Abort("500") // Error no manejado!
			}
		}
	}()

	if v, err := propiedades.GetPropiedades(); err == nil {
		c.Data["json"] = v
		c.ServeJSON()
	} else {
		panic(err)
	}
}

// GetDependenciaById ...
// @Title GetAll
// @Description get Dependencia with the specified {idTercero}
// @Param	propiedad	path 	string	true		"type propiedad of Terceros"
// @Param	idTercero	path 	string	true		"Identify Dependencia by IdTercero"
// @Success 200 {object} []map[string]interface{}
// @Failure 500 Internal Error
// @Failure 501 {user} Not Implemented
// @Failure 502 Error with external API
// @router /:propiedad/:idTercero [get]
func (c *PropiedadesController) GetDependenciaById() {
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

	idTercero := c.Ctx.Input.Param(":idTercero")
	propiedad := c.Ctx.Input.Param(":propiedad")

	if helper, err := propiedades.GetHelperPropiedades(propiedad); err == nil {
		if v, err := helper(idTercero); err == nil {
			c.Data["json"] = v
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
