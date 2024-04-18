package tipos

import (
	"fmt"
	"strings"

	ParametrosCrudModels "github.com/udistrital/parametros_crud/models"
	ParametrosHelper "github.com/udistrital/terceros_mid/helpers/crud/parametros"
	TercerosHelper "github.com/udistrital/terceros_mid/helpers/crud/terceros"
)

// GetFuncionarios trae los terceros que tienen un registro en la tabla vinculacion del api terceros_crud
func GetFuncionarios(idTercero int, query string) (terceros []map[string]interface{}, outputError map[string]interface{}) {
	const funcion = "GetFuncionarios - "
	step := "0"
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": funcion + "uncaught error after step:" + step,
				"err":     err,
				"status":  "500", // Uncaught error!
			}
			panic(outputError)
		}
	}()

	if query != "" {
		err := errors.New("query no implementado")
		return nil, e.Error(funcion+`query != ""`, err, fmt.Sprint(http.StatusNotImplemented))
	}

	var vinculaciones []models.Vinculacion
	urlTerceros := "http://" + beego.AppConfig.String("TercerosService") + "vinculacion?limit=-1"
	urlTerceros += "&fields=Id,TerceroPrincipalId,TipoVinculacionId,DependenciaId"
	urlTerceros += "&query=Activo:true"
	if idTercero > 0 {
		urlTerceros += ",TerceroPrincipalId__Id:" + fmt.Sprint(idTercero)
	}
	if resp, err := request.GetJsonTest(urlTerceros, &vinculaciones); err == nil && resp.StatusCode == 200 {

		if len(vinculaciones) == 0 || vinculaciones[0].TerceroPrincipalId == nil {
			var tercero models.Tercero
			urlTerceros = "http://" + beego.AppConfig.String("TercerosService") + "tercero/" + fmt.Sprint(idTercero)
			if resp, err := request.GetJsonTest(urlTerceros, &tercero); err == nil && resp.StatusCode == 200 {
				terceros = append(terceros, map[string]interface{}{
					"Tercero": tercero,
				})
				return terceros, nil
			} else {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "/GetFuncionarios - request.GetJsonTest(urlTerceros, &tercero)",
					"err":     err,
					"status":  "502",
				}
				return nil, outputError
			}
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
	empty := []string{}

	// PARTE 1. Traer los ID de los parámetros asociados a funcionarios

	// Los siguientes son los códigos de los registros de la tabla "parametro" de la API
	// de parámetros, cuyo tipo_parámetro_id sea el asociado a Cargos.
	// Específicamente, los códigos de parámetros que se asignen a funcionarios
	codigosParametroFuncionarios := []string{"P", "AP", "DCTC", "DCMT", "A", "O", "TO", "PD", "PA", "DP", "AS", "E", "D", "EP", "PU", "T"}
	codigoTipoParamVinculacion := "TV"

	fieldsParametros := []string{"Id", "CodigoAbreviacion"}
	queryParametros := "Activo:true,TipoParametroId__Activo:true"
	queryParametros += ",TipoParametroId__CodigoAbreviacion:" + codigoTipoParamVinculacion
	queryParametros += ",CodigoAbreviacion__in:" + strings.Join(codigosParametroFuncionarios, "|")

			var resBody []models.AsignacionEspacioFisicoDependencia
			urlOikos := "http://" + beego.AppConfig.String("OikosService") + "asignacion_espacio_fisico_dependencia?limit=-1"
			urlOikos += "&fields=Id,EspacioFisicoId,DependenciaId&query=Activo:true"
			urlOikos += ",EspacioFisicoId__TipoEspacioFisicoId__CodigoAbreviacion:Tipo_1"
			urlOikos += ",DependenciaId__Id:" + fmt.Sprint(tercero["DependenciaId"])
			if resp, err := request.GetJsonTest(urlOikos, &resBody); err == nil && resp.StatusCode == 200 {
	var parametros []ParametrosCrudModels.Parametro
	step = "1"

	outputError = ParametrosHelper.GetParametros(&parametros, queryParametros, -1, 0, fieldsParametros, empty, empty)
	if outputError != nil {
		return
	}

	var vinculos = []string{}
	step = "2"
	for _, parametro := range parametros {
		vinculos = append(vinculos, fmt.Sprint(parametro.Id))
	}

	terceros, outputError = TercerosHelper.GetTrVinculacionIdentificacion(query, strings.Join(vinculos, ","), "", "")

	return

}
