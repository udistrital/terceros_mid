package tipos

import (
	"fmt"
	"strings"

	ParametrosCrudModels "github.com/udistrital/parametros_crud/models"
	TercerosCrudModels "github.com/udistrital/terceros_crud/models"
	ParametrosHelper "github.com/udistrital/terceros_mid/helpers/crud/parametros"
	TercerosHelper "github.com/udistrital/terceros_mid/helpers/crud/terceros"
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

	empty := []string{}

	// PARTE 1. Traer los ID de los parámetros asociados a contratistas

	// Los siguientes son los códigos de los registros de la tabla "parametro" de la API
	// de parámetros, cuyo tipo_parámetro_id sea el asociado a Cargos.
	// Específicamente, los códigos de parámetros que se asignen a contratistas
	codigosParametroContratista := []string{"CPS"} // , "OPS", "PS"
	codigoTipoParamVinculacion := "TV"
	parametroContratistaID := make(map[string]int)

	fieldsParametros := []string{"Id", "CodigoAbreviacion"}
	queryParametros := "Activo:true,TipoParametroId__Activo:true"
	queryParametros += ",TipoParametroId__CodigoAbreviacion:" + codigoTipoParamVinculacion
	// Descomentar la siguiente línea una vez se tenga soporte __in en los queries de parametros_crud...
	// queryParametros += ",CodigoAbreviacion__in:" + strings.Join(codigosParametroContratista, "|")
	var parametros []ParametrosCrudModels.Parametro
	step = "1"
	if outputError = ParametrosHelper.GetParametros(&parametros,
		queryParametros, -1, 0, fieldsParametros, empty, empty); outputError != nil {
		return
	}
	// ... y así se podría eliminar/reducir lo siguiente:
	for _, parametro := range parametros {
		for _, codigoContratista := range codigosParametroContratista {
			if parametro.CodigoAbreviacion == codigoContratista {
				parametroContratistaID[codigoContratista] = parametro.Id
				break
			}
		}
	}
	// logs.Debug("parametroContratistaID:", parametroContratistaID)

	// PARTE 2. Traer los terceros que tengan los ID anteriores en la tabla vinculacion

	// vinculacionesMap := make(map[int]models.Vinculacion)
	var vinculos = []string{}
	step = "2"
	for _, id := range parametroContratistaID {
		vinculos = append(vinculos, fmt.Sprint(id))
	}

	documentosMap := make(map[int]TercerosCrudModels.DatosIdentificacion)
	consultar := func(queryTercero, queryDocumento string) {
		var vinculacionesTerceros []TercerosCrudModels.Vinculacion
		fullQueryVinculaciones := "Activo:true"
		if idTercero > 0 {
			fullQueryVinculaciones += ",TerceroPrincipalId:" + fmt.Sprint(idTercero)
		}
		if queryTercero != "" {
			fullQueryVinculaciones += ",TerceroPrincipalId__NombreCompleto__icontains:" + queryTercero
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
	}
	consultar("", "")

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
