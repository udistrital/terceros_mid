package propiedades

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

func GetCargo(idTercero string) (cargo []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetCargo - Unhandled Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	// PARTE 1. Consultar el tercero con el Id
	// Se obtiene el CargoId del tercero en la tabla vinculacion
	var idCargo []map[string]interface{}
	urlVinculacion := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=1"
	urlVinculacion += "&fields=CargoId&sortby=Id&order=desc"
	urlVinculacion += "&query=Activo%3Atrue,TerceroPrincipalId%3A" + idTercero

	if resp, err := request.GetJsonTest(urlVinculacion, &idCargo); err == nil && resp.StatusCode == 200 {
		if len(idCargo) == 0 || len(idCargo[0]) == 0 {
			err := fmt.Errorf("El tercero" + idTercero + "No tiene una vinculación activa")
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetCargo - request.GetJsonTest(urlVinculacion, &idCargo)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		} else if idCargo[0]["CargoId"].(float64) == 0 {
			l := []map[string]interface{}{}
			return l, nil
		}
	} else {
		if err == nil {
			err = fmt.Errorf("undesired status code - got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCargo - equest.GetJsonTest(urlVinculacion, &idCargo)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// PARTE 2. Se consulta el CargoId en el api parametros_crud

	var parametro map[string]interface{}
	urlParametro := "http://" + beego.AppConfig.String("parametrosService") + "parametro/" + fmt.Sprintf("%v", idCargo[0]["CargoId"])
	if resp, err := request.GetJsonTest(urlParametro, &parametro); err == nil && resp.StatusCode == 200 {
		if len(parametro) == 0 {
			err := fmt.Errorf("No existe el parametro con el id " + fmt.Sprintf("%v", idCargo[0]["CargoId"]))
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetCargo - request.GetJsonTest(urlParametro, &parametro)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
		cargo = append(cargo, parametro["Data"].(map[string]interface{}))
	} else {
		if err == nil {
			err = fmt.Errorf("undesired status code - got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetCargo - request.GetJsonTest(urlParametro, &parametro)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return cargo, nil
}
