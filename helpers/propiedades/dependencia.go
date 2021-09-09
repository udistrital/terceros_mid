package propiedades

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

func GetDependencia(idTercero string) (dependencia []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetDependencia - Unhandled Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()
	// PARTE 1. Consultar el Id del tercero con el tipo de documento y documento

	// Document es el numero de documento con el que se consultara el terero.
	// Tipo es el valor del tipo de documento para filtar el documento con su tipo,
	// para evitar errores de documentos duplicados con diferentes tipos.

	// Esto nos devolvera el id del tercero asociado para consultar las dependencias luego.
	//var idTercero []map[string]interface{}
	//var tercero models.Tercero
	//urlTercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion?limit=-1"
	//urlTercero += "&fields=TerceroId"
	//urlTercero += "&query=TipoDocumentoId__CodigoAbreviacion%3A" + tipo + "%2CNumero%3A" + document
	//if resp, err := request.GetJsonTest(urlTercero, &idTercero); err == nil && resp.StatusCode == 200 {
	//	if len(idTercero) == 0 || len(idTercero[0]) == 0 {
	//		err := fmt.Errorf("No hay terceros registrados con el documento " + document)
	//		logs.Error(err)
	//		outputError = map[string]interface{}{
	//			"funcion": "GetDependencia - request.GetJsonTest(urlTercero, &idTercero)",
	//			"err":     err,
	//			"status":  "502",
	//		}
	//		return nil, outputError
	//	}
	//	jsonString, _ := json.Marshal(idTercero[0]["TerceroId"])
	//	bytes := []byte(jsonString)
	//	json.Unmarshal(bytes, &tercero)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//} else {
	//	if err == nil {
	//		err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
	//	}
	//	logs.Error(err)
	//	outputError = map[string]interface{}{
	//		"funcion": "GetDependencia - request.GetJsonTest(urlTercero, &idTercero)",
	//		"err":     err,
	//		"status":  "502",
	//	}
	//	return nil, outputError
	//}
	// PARTE 2. Consultar el tercero con el Id del paso anterior
	// Se obtendra la dependencia que estan relacionados al tercero
	var idDependencia []map[string]interface{}
	urlDependencia := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=-1"
	urlDependencia += "&fields=DependenciaId"
	urlDependencia += "&query=TerceroPrincipalId%3A" + idTercero
	if resp, err := request.GetJsonTest(urlDependencia, &idDependencia); err == nil && resp.StatusCode == 200 {
		if len(idDependencia) == 0 || len(idDependencia[0]) == 0 {
			err := fmt.Errorf("No hay dependencias registradas con el id " + idTercero)
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetDependencia - request.GetJsonTest(urlDependencia, &idDependencia)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	} else {
		if err == nil {
			err = fmt.Errorf("undesired status code - got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDependencia - request.GetJsonTest(urlDependencia, &idDependencia)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// PARTE 3. Consultar las dependencias a partir del id identificada en vinculacion
	// Se obtendra las dependencia que estan relacionados al tercero por medio del id
	// del tercero y se adjuntaran a un array que sera la respuesta que se devolvera

	// Este proceso se repetira la cantidad de id Dependencias se hallan encontrado
	var temp map[string]interface{}
	var Dependencias []map[string]interface{}
	for _, idDep := range idDependencia[0] {
		urlDependencias := "http://" + beego.AppConfig.String("oikos2Service") + "dependencia/" + fmt.Sprintf("%v", idDep) + "?limit=-1"
		if resp, err := request.GetJsonTest(urlDependencias, &temp); err == nil && resp.StatusCode == 200 {
			if len(temp) == 0 {
				err := fmt.Errorf("No hay dependencias registradas con el id " + fmt.Sprintf("%v", idDep))
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetDependencia - request.GetJsonTest(urlDependencias, &Dependencias)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
			Dependencias = append(Dependencias, temp)
		} else {
			if err == nil {
				err = fmt.Errorf("undesired status code - got:%d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetDependencia - request.GetJsonTest(urlDependencias, &Dependencias)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
	}
	return Dependencias, nil
}
