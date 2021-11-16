package tipos

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/mapstructure"
	"github.com/udistrital/terceros_mid/helpers/propiedades"
	"github.com/udistrital/terceros_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetFuncionarios trae los terceros que tienen un registro en la tabla vinculacion del api terceros_crud
func GetFuncionarios(idTercero int) (terceros []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetFuncionarios - Uncaught Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	var vinculaciones []models.Vinculacion
	urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=-1"
	urlTerceros += "&fields=Id,TerceroPrincipalId,TipoVinculacionId,DependenciaId"
	urlTerceros += "&query=Activo:true"
	if idTercero > 0 {
		urlTerceros += ",TerceroPrincipalId__Id:" + fmt.Sprint(idTercero)
	}
	if resp, err := request.GetJsonTest(urlTerceros, &vinculaciones); err == nil && resp.StatusCode == 200 {

		if len(vinculaciones) == 0 || vinculaciones[0].Id == 0 {
			return nil, nil
		}
		// fmt.Println("paramId:", paramID, "#vinculaciones: ", len(vinculaciones))

		// Lo siguiente es para que no se vuelva a agregar un tercero
		// cuando el tercero tenga más de una vinculación
		for _, vincul := range vinculaciones {
			add := true
			for _, tercero := range terceros {
				if mTercero := tercero["Tercero"].(*models.Tercero); vincul.TerceroPrincipalId.Id == mTercero.Id {
					add = false
					break
				}
			}
			if add {
				terceros = append(terceros, map[string]interface{}{
					"Tercero":         vincul.TerceroPrincipalId,
					"TipoVinculacion": vincul.TipoVinculacionId,
					"DependenciaId":   vincul.DependenciaId,
				})
			}
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetFuncionarios - request.GetJsonTest(urlTerceros, &vinculaciones)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// fmt.Println("#terceros:", len(terceros))

	// PARTE 3 - Agregar Información complementaria de Sede y Dependencia (si la hay)

	var sedesDependencias []models.AsignacionEspacioFisicoDependencia
	for _, tercero := range terceros {
		// fmt.Println("k:", k, "tercero:", tercero)

		// 3.1 traer los registros necesarios/disponibles
		var terceroModelo models.Tercero
		if err := mapstructure.Decode(tercero["Tercero"], &terceroModelo); err != nil {
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetContratista - mapstructure.Decode(tercero[\"Tercero\"], &terceroModelo)",
				"err":     err,
				"status":  "500",
			}
			return nil, outputError
		} else {
			if identificacion, err := propiedades.GetDocumento(fmt.Sprint(terceroModelo.Id)); err != nil && len(identificacion) > 0 {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetContratista - mapstructure.Decode(tercero[\"Tercero\"], &terceroModelo)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			} else {
				if len(identificacion) > 0 {
					tercero["Identificacion"] = identificacion[0]
				}
			}
		}

		consultar := true
		for _, seDependencia := range sedesDependencias {
			if seDependencia.DependenciaId.Id == tercero["Dependencia"] {
				consultar = false
			}
		}
		if consultar {

			var resBody []models.AsignacionEspacioFisicoDependencia
			urlOikos := "http://" + beego.AppConfig.String("oikos2Service") + "asignacion_espacio_fisico_dependencia?limit=-1"
			urlOikos += "&fields=Id,EspacioFisicoId,DependenciaId&query=Activo:true"
			urlOikos += ",EspacioFisicoId__TipoEspacioFisicoId__CodigoAbreviacion:Tipo_1"
			urlOikos += ",DependenciaId__Id:" + fmt.Sprint(tercero["DependenciaId"])
			if resp, err := request.GetJsonTest(urlOikos, &resBody); err == nil && resp.StatusCode == 200 {

				if len(resBody) == 0 || resBody[0].Id == 0 {
					// No se encontró relación sede-dependencia para el tercero actual
					continue
				}

				for _, v := range resBody {
					sedesDependencias = append(sedesDependencias, v)
				}

			} else {
				if err == nil {
					err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
				}
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetFuncionarios - request.GetJsonTest(urlOikos, &resBody)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
		}

		// 3.2 asignar la información disponible
		for _, seDep := range sedesDependencias {
			if tercero["DependenciaId"] == seDep.DependenciaId.Id {
				tercero["Sede"] = seDep.EspacioFisicoId
				tercero["Dependencia"] = seDep.DependenciaId
				break
			}
		}
	}

	return terceros, nil
}
