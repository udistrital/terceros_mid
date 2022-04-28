package tipos

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/terceros_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetJefeDependencia obtiene la información del tercero y el cargo del jefe de la dependencia indicada
func GetJefeDependencia(dependenciaId int) (jefeDependencia []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetJefeDependencia - Uncaught Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	var dependencia models.Dependencia
	var terceros []models.Vinculacion
	var tercero models.Vinculacion
	var respParam models.RespuestaAPI1Obj

	url := "http://" + beego.AppConfig.String("oikos2Service") + "dependencia/" + fmt.Sprint(dependenciaId)
	if resp, err := request.GetJsonTest(url, &dependencia); err == nil && resp.StatusCode == 200 {
		urlTercero := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=-1"
		urlTercero += "&fields=Id,TerceroPrincipalId,DependenciaId,CargoId"
		urlTercero += "&query=Activo:true,DependenciaId:" + fmt.Sprint(dependenciaId)

		if resp2, err2 := request.GetJsonTest(urlTercero, &terceros); err2 == nil && resp2.StatusCode == 200 {
			if len(terceros) > 0 {
				tercero = terceros[0]
				urlCargo := "http://" + beego.AppConfig.String("parametrosService") + "parametro/" + fmt.Sprint(tercero.CargoId)
				if resp3, err3 := request.GetJsonTest(urlCargo, &respParam); err3 == nil && resp3.StatusCode == 200 {
					jefeDependencia = append(jefeDependencia, map[string]interface{}{
						"TerceroPrincipal": tercero.TerceroPrincipalId.NombreCompleto,
						"Cargo":            respParam.Data["Nombre"],
						"DependenciaId":    dependencia.Nombre,
					})
				} else {
					if err3 == nil {
						err3 = fmt.Errorf("Undesired status code - Got:%d", resp3.StatusCode)
					}
					logs.Error(err3)
					outputError = map[string]interface{}{
						"funcion": "/GetJefeDependencia - request.GetJsonTest(urlCargo, &respParam)",
						"err":     err3,
						"status":  "502",
					}
					return nil, outputError
				}
			}
		} else {
			if err2 == nil {
				err2 = fmt.Errorf("Undesired status code - Got:%d", resp2.StatusCode)
			}
			logs.Error(err2)
			outputError = map[string]interface{}{
				"funcion": "/GetJefeDependencia - request.GetJsonTest(urlTercero, &terceros)",
				"err":     err2,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetJefeDependencia - request.GetJsonTest(url, &dependencia)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}

	return jefeDependencia, outputError
}
