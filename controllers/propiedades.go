package controllers

import (
	"github.com/astaxie/beego"

	"github.com/udistrital/terceros_mid/helpers/propiedades"
	e "github.com/udistrital/utils_oas/errorctrl"
)

// PropiedadesController operations for propiedades
type PropiedadesController struct {
	beego.Controller
}

// URLMapping ...
func (c *PropiedadesController) URLMapping() {
	c.Mapping("GetPropiedades", c.GetPropiedades)
	c.Mapping("GetPropiedadesDeUnTerceroId", c.GetPropiedadesDeUnTerceroId)
}

// GetPropiedades ...
// @Title GetAll
// @Description List the Propiedades types that can be used to gather Propiedades by {propiedad}
// @Success 200 {object} []string
// @router / [get]
func (c *PropiedadesController) GetPropiedades() {
	// Puede que ni sea necesario en este controlador, pero se coloca por lineamiento...
	defer e.ErrorControlController(c.Controller, "PropiedadesController")

	if v, err := propiedades.GetPropiedades(); err == nil {
		if len(v) > 0 {
			c.Data["json"] = v
		} else {
			c.Data["json"] = []interface{}{}
		}
		c.ServeJSON()
	} else {
		panic(err)
	}
}

// GetPropiedadesDeUnTerceroId ...
// @Title GetAll
// @Description get Dependencia with the specified {idTercero}
// @Param	propiedad	path 	string	true		"type propiedad of Terceros"
// @Param	idTercero	path 	string	true		"Tercero ID from terceros_crud"
// @Success 200 {object} []map[string]interface{}
// @Failure 500 Internal Error
// @Failure 501 {user} Not Implemented
// @Failure 502 Error with external API
// @router /:propiedad/:idTercero [get]
func (c *PropiedadesController) GetPropiedadesDeUnTerceroId() {
	defer e.ErrorControlController(c.Controller, "PropiedadesController")

	idTercero := c.Ctx.Input.Param(":idTercero")
	propiedad := c.Ctx.Input.Param(":propiedad")

	if helper, err := propiedades.GetHelperPropiedades(propiedad); err == nil {
		if v, err := helper(idTercero); err == nil {
			if len(v) > 0 {
				c.Data["json"] = v
			} else {
				c.Data["json"] = []interface{}{}
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
