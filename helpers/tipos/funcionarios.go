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
