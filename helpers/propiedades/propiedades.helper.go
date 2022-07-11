package propiedades

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"

	e "github.com/udistrital/utils_oas/errorctrl"
)

var diccionarioPropiedadesHelper = map[string](func(string) ([]map[string]interface{}, map[string]interface{})){
	"dependencia": GetDependencia,
	"documento":   GetDocumento,
	"cargo":       GetCargo,
}

// GetPropiedades retorna la lista de propiedades que pueden ser usados con GetHelperPropiedades
func GetPropiedades() (propiedad []string, outputError map[string]interface{}) {
	// Puede que ni sea necesario en este helper, pero se coloca por lineamiento...
	const funcion = "GetPropiedades - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	for k := range diccionarioPropiedadesHelper {
		propiedad = append(propiedad, k)
	}
	return propiedad, nil
}

// GetHelperTipo trae los terceros con el criterio especificado.
// El criterio debe ser alguno de los valores retornados por GetTipos
func GetHelperPropiedades(propiedad string) (helper func(string) ([]map[string]interface{}, map[string]interface{}), outputError map[string]interface{}) {
	const funcion = "GetHelperPropiedades - "
	defer e.ErrorControlFunction(funcion+"Unhandled Error!", fmt.Sprint(http.StatusInternalServerError))

	if helper, found := diccionarioPropiedadesHelper[propiedad]; found {
		return helper, nil
	}

	err := fmt.Errorf("\"%s\" not implemented", propiedad)
	logs.Error(err)

	return nil, e.Error(funcion+"helper, found := diccionarioPropiedadesHelper[propiedad]",
		err, fmt.Sprint(http.StatusNotFound))
}
