package tipos

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	TercerosCrudModels "github.com/udistrital/terceros_crud/models"
	TercerosHelper "github.com/udistrital/terceros_mid/helpers/crud/terceros"
	"github.com/udistrital/terceros_mid/models"
	"github.com/udistrital/utils_oas/request"
)

// GetContratista trae la lista de contratistas registrados en Terceros, con opción de filtrar por ID
func GetContratista(idTercero int, query string) (terceros []map[string]interface{}, outputError map[string]interface{}) {
	const funcion = "GetContratista - "
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

	// PARTE 1. Traer los ID de los parámetros asociados a contratistas

	// Los siguientes son los códigos de los registros de la tabla "parametro" de la API
	// de parámetros, cuyo tipo_parámetro_id sea el asociado a Cargos.
	// Específicamente, los códigos de parámetros que se asignen a contratistas
	codigosParametroContratista := []string{"CPS"} // , "OPS", "PS"
	codigoTipoParamVinculacion := "TV"
	parametroContratistaID := make(map[string]int)

	var respBody models.RespuestaAPI1Arr
	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro?limit=-1"
	urlParametros += "&fields=Id,CodigoAbreviacion"
	urlParametros += "&query=Activo:true,TipoParametroId__Activo:true,TipoParametroId__CodigoAbreviacion:" + codigoTipoParamVinculacion
	// logs.Debug("urlParametros:", urlParametros)
	step = "1"
	if resp, err := request.GetJsonTest(urlParametros, &respBody); err == nil && resp.StatusCode == 200 {
		step = "1.1"
		if respBody.Data == nil || len(respBody.Data) == 0 || len(respBody.Data[0]) == 0 {
			err := fmt.Errorf("No están registrados los parámetros asociados a contratistas")
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "GetContratista - respBody.Data == nil || len(respBody.Data) == 0 || len(respBody.Data[0]) == 0",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
		step = "1.2"
		for _, paramVinculacion := range respBody.Data {
			// fmt.Printf("Param #%d: %#v\n", k, paramVinculacion)
			var codParam string
			if v, ok := paramVinculacion["CodigoAbreviacion"].(string); ok {
				codParam = v
			} else {
				continue
			}
			// fmt.Printf("codParam (%T): %v\n", codParam, codParam)
			for _, codigoContratista := range codigosParametroContratista {
				if id, ok := paramVinculacion["Id"].(float64); ok && codigoContratista == codParam {
					// fmt.Printf("P=P %v - T(id):%T - v:%f\n", paramVinculacion, paramVinculacion["Id"], paramVinculacion["Id"])
					parametroContratistaID[codigoContratista] = int(id)
					break
				}
			}
		}
		step = "1.3"
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetContratista - request.GetJsonTest(urlParametros, &respBody)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// logs.Debug("parametroContratistaID:", parametroContratistaID)

	// PARTE 2. Traer los terceros que tengan los ID anteriores en la tabla vinculacion
	empty := []string{}

	// vinculacionesMap := make(map[int]models.Vinculacion)
	var vinculos = []string{}
	step = "2"
	for _, id := range parametroContratistaID {
		vinculos = append(vinculos, fmt.Sprint(id))
	}
	var vinculacionesTerceros []TercerosCrudModels.Vinculacion
	fullQueryVinculaciones := "Activo:true"
	if idTercero > 0 {
		fullQueryVinculaciones += ",TerceroPrincipalId:" + fmt.Sprint(idTercero)
	}
	fullQueryVinculaciones += ",TipoVinculacionId__in:" + strings.Join(vinculos, "|")
	limit := -1
	offset := 0
	fieldsVinculaciones := []string{"Id", "TerceroPrincipalId"}
	step = "3"
	if err := TercerosHelper.GetVinculaciones(&vinculacionesTerceros, fullQueryVinculaciones, limit, offset, fieldsVinculaciones, empty, empty); err != nil {
		outputError = err
		return
	}

	tercerosMap := make(map[int]TercerosCrudModels.Tercero)
	for _, v := range vinculacionesTerceros {
		tercerosMap[v.TerceroPrincipalId.Id] = *v.TerceroPrincipalId
	}
	// logs.Debug("tercerosMap:", tercerosMap)
	// logs.Debug("vinculacionesTerceros:", vinculacionesTerceros)

	documentosMap := make(map[int]TercerosCrudModels.DatosIdentificacion)
	for terceroId := range tercerosMap {
		fullQueryDocumentos := "Activo:true,TerceroId__Activo:true,TerceroId__Id:" + fmt.Sprint(terceroId)
		fieldsDocumentos := []string{"Id", "TipoDocumentoId", "Numero"}
		var documentosTerceros []TercerosCrudModels.DatosIdentificacion
		step = "4"
		if err := TercerosHelper.GetDatosIdentificacion(&documentosTerceros, fullQueryDocumentos, limit, offset, fieldsDocumentos, empty, empty); err != nil {
			outputError = err
			return
		}
		step = "5"
		// logs.Debug("documentosTerceros:", fmt.Sprintf("%+v", documentosTerceros))
		fin := len(documentosTerceros)
		for k, v := range documentosTerceros {
			step = fmt.Sprintf("5.%d/%d", k, fin)
			var completo TercerosCrudModels.DatosIdentificacion = v
			if tercero, ok := tercerosMap[terceroId]; ok {
				completo.TerceroId = &tercero
			}
			documentosMap[v.Id] = completo
		}
	}

	current := 1
	fin := len(documentosMap)
	for k, v := range documentosMap {
		step = fmt.Sprint("6.%d/%d(idDoc:%d)", current, fin, k)
		terceros = append(terceros, map[string]interface{}{
			"Tercero": map[string]interface{}{
				"Id":             v.TerceroId.Id,
				"NombreCompleto": v.TerceroId.NombreCompleto,
				// "UsuarioWSO2":    vincul.TerceroPrincipalId.UsuarioWSO2,
			},
			"Identificacion": map[string]interface{}{
				// "TipoDocumentoId": dataModel.TipoDocumentoId,
				"TipoDocumentoId": map[string]interface{}{
					"Id":                v.TipoDocumentoId.Id,
					"Nombre":            v.TipoDocumentoId.Nombre, // TODO: Revisar otras APIs y eliminar este campo
					"CodigoAbreviacion": v.TipoDocumentoId.CodigoAbreviacion,
				},
				"Numero": v.Numero,
			},
		})
		current++
	}

	return
}
