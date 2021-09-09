package propiedades

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

var diccionarioPropiedadesHelper = map[string](func(string) ([]map[string]interface{}, map[string]interface{})){
	"dependencia": GetDependencia,
}

// GetPropiedades retorna la lista de propiedades que pueden ser usados con GetHelperPropiedades
func GetPropiedades() (propiedad []string, outputError map[string]interface{}) {

	// Puede que ni sea necesario en este helper, pero se coloca por lineamiento...
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetTipos - Unhandled Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	for k := range diccionarioPropiedadesHelper {
		propiedad = append(propiedad, k)
	}
	return propiedad, nil
}

// GetHelperTipo trae los terceros con el criterio especificado.
// El criterio debe ser alguno de los valores retornados por GetTipos
func GetHelperPropiedades(propiedad string) (helper func(string) ([]map[string]interface{}, map[string]interface{}), outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetHelperTipo - Unhandled Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	if helper, found := diccionarioPropiedadesHelper[propiedad]; found {
		return helper, nil
	}

	err := fmt.Errorf("\"%s\" not implemented", propiedad)
	logs.Error(err)

	return nil, map[string]interface{}{
		"funcion": "GetHelperTipo - found := diccionarioTipoHelper[tipo]",
		"err":     err,
		"status":  "404",
	}
}
