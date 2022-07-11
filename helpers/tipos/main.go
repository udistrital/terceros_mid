package tipos

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"

	e "github.com/udistrital/utils_oas/errorctrl"
)

var diccionarioTipoHelper = map[string](func(int) ([]map[string]interface{}, map[string]interface{})){
	"funcionarioPlanta": GetFuncionariosPlanta,
	"ordenadoresGasto":  GetOrdenadores,
	"contratista":       GetContratista,
	"proveedor":         GetProveedor,
	"funcionarios":      GetFuncionarios,
	"jefeDependencia":   GetJefeDependencia,
}

// GetTipos retorna la lista de tipos que pueden ser usados con GetHelperTipo
func GetTipos() (tercero []string, outputError map[string]interface{}) {
	// Puede que ni sea necesario en este helper, pero se coloca por lineamiento...
	const funcion = "GetTipos - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	for k := range diccionarioTipoHelper {
		tercero = append(tercero, k)
	}
	return tercero, nil
}

// GetHelperTipo trae los terceros con el criterio especificado.
// El criterio debe ser alguno de los valores retornados por GetTipos
func GetHelperTipo(tipo string) (helper func(int) ([]map[string]interface{}, map[string]interface{}), outputError map[string]interface{}) {
	const funcion = "GetHelperTipo"
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	if helper, found := diccionarioTipoHelper[tipo]; found {
		return helper, nil
	}

	err := fmt.Errorf("\"%s\" not implemented", tipo)
	logs.Error(err)

	return nil, e.Error(funcion+"helper, found := diccionarioTipoHelper[tipo]",
		err, fmt.Sprint(http.StatusNotFound))
}
