package controllers

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/terceros_mid/helpers/tipos"
	e "github.com/udistrital/utils_oas/errorctrl"
)

// TercerosController operations for Terceros
type TiposController struct {
	beego.Controller
}

// URLMapping ...
func (c *TiposController) URLMapping() {
	c.Mapping("GetTipos", c.GetTipos)
	c.Mapping("GetByTipo", c.GetByTipo)
	c.Mapping("GetByTipoAndId", c.GetByTipoAndID)
}

// GetTipos ...
// @Title GetAll
// @Description List the Tercero types that can be used to gather Terceros by {tipo}
// @Success 200 {object} []string
// @router / [get]
func (c *TiposController) GetTipos() {

	// Puede que ni sea necesario en este controlador, pero se coloca por lineamiento...
	defer e.ErrorControlController(c.Controller, "TiposController")

	if v, err := tipos.GetTipos(); err == nil {
		if len(v) > 0 {
			c.Data["json"] = v
		} else {
			c.Data["json"] = []interface{}{}
		}
		c.Data["json"] = v
		c.ServeJSON()
	} else {
		panic(err)
	}
}

// GetByTipo ...
// @Title GetAll
// @Description get Terceros with the specified {tipo}
// @Param	tipo	path 	string	true		"Tercero type available from /tipo/"
// @Param	query	query	string	false		"Query param"
// @Success 200 {object} []map[string]interface{}
// @Failure 500 Internal Error
// @Failure 501 {tipo} Not Implemented
// @Failure 502 Error with external API
// @router /:tipo [get]
func (c *TiposController) GetByTipo() {
	defer e.ErrorControlController(c.Controller, "TiposController")

	tipo := c.Ctx.Input.Param(":tipo")
	query := c.GetString("query")
	logs.Debug("query:", query)

	if helper, err := tipos.GetHelperTipo(tipo); err == nil {
		if v, err := helper(0, query); err == nil {
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

// GetByTipoAndId ...
// @Title GetAll
// @Description get Terceros with the specified {tipo} and {id} of a record in terceros table from Terceros CRUD API
// @Param	tipo path  string true		"Tercero type available from /tipo/"
// @Param	id   path  uint   true		"Tercero ID. Must be positive when specified"
// @Success 200 {object} []map[string]interface{}
// @Failure 400 Wrong ID
// @Failure 404 ID with {tipo} Not Found
// @Failure 500 Internal Error
// @Failure 501 {tipo} Not Implemented
// @Failure 502 Error with external API
// @router /:tipo/:id [get]
func (c *TiposController) GetByTipoAndID() {
	const funcion = "GetByTipoAndID - "
	defer e.ErrorControlController(c.Controller, "TiposController")

	tipo := c.Ctx.Input.Param(":tipo")
	var (
		id  int
		err error
	)
	if id, err = c.GetInt(":id"); err != nil || id < 0 {
		if err == nil {
			err = fmt.Errorf("ID MUST be greater than 0 - Got: %d", id)
		}
		panic(e.Error(funcion+`c.GetInt(":id")`, err, fmt.Sprint(http.StatusBadRequest)))
	}

	if helper, err := tipos.GetHelperTipo(tipo); err == nil {
		if v, err := helper(id, ""); err == nil {
			if len(v) == 0 {
				err := fmt.Errorf("no se encontró un Tercero tipo '%s' con id '%d'", tipo, id)
				panic(e.Error(funcion+"len(v) == 0", err, fmt.Sprint(http.StatusNotFound)))
			}
			c.Data["json"] = v
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	c.ServeJSON()
}
