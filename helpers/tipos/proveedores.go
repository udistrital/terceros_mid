package tipos

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"

	TercerosCrudModels "github.com/udistrital/terceros_crud/models"
	TercerosHelper "github.com/udistrital/terceros_mid/helpers/crud/terceros"
	e "github.com/udistrital/utils_oas/errorctrl"
)

// GetProveedor trae la lista de proveedores registrados en Terceros, con opcion de filtrar por ID
func GetProveedor(idProveedor int, query string) (terceros []map[string]interface{}, outputError map[string]interface{}) {
	const funcion = "GetProveedor - "
	defer e.ErrorControlFunction(funcion+"Uncaught Error!", fmt.Sprint(http.StatusInternalServerError))

	if query != "" {
		err := errors.New("query no implementado")
		return nil, e.Error(funcion+`query != ""`, err, fmt.Sprint(http.StatusNotImplemented))
	}

	// PARTE 1. Traer los ID de los tipo_tercero asociados a proveedores
	// Eliminada: Desde que esté registrado en Agora/terceros, debería poder
	// seleccionarse como proveedor. Si se requiere a futuro el reestablecer
	// algun criterio de filtrado, se podrían agregar opciones __in al query

	// PARTE 2 - Traer los terceros

	// TODO: Dar soporte a paginación, traer limit y offset como argumentos de la función
	const limit = -1
	const offset = 0
	empty := []string{}
	fields := []string{"Id", "NombreCompleto"}
	query2 := "Activo:true"
	if idProveedor > 0 {
		query2 += ",Id:" + fmt.Sprint(idProveedor)
	}

	var tercerosMap []TercerosCrudModels.Tercero
	if err := TercerosHelper.GetTerceros(&tercerosMap, query2, limit, offset, fields, empty, empty); err != nil {
		outputError = err
		return
	}
	if len(tercerosMap) == 0 {
		return
	}
	// logs.Debug("tercerosMap:", tercerosMap)

	// PARTE 3: Completar información de identificación, de estar disponible
	fields = []string{"TipoDocumentoId", "Numero"}
	const limitDocs = 10
	for _, dataTercero := range tercerosMap {

		dataFinal := map[string]interface{}{
			"Tercero": map[string]interface{}{
				"Id":             dataTercero.Id,
				"NombreCompleto": dataTercero.NombreCompleto,
			},
		}

		var dataTerceros []TercerosCrudModels.DatosIdentificacion
		query2 = "Activo:true,TerceroId__Id:" + fmt.Sprint(dataTercero.Id)
		// TODO: Usar un buffer antes de consultar
		if err := TercerosHelper.GetDatosIdentificacion(&dataTerceros, query2, limitDocs, 0, fields, empty, empty); err != nil {
			outputError = err
			return
		}
		// TODO: Retornar los documentos únicos activos que tiene el tercero, en un arreglo o mapeo
		if len(dataTerceros) == 1 && dataTerceros[0].Numero != "" {
			dataRecortada := map[string]interface{}{
				"TipoDocumentoId": map[string]interface{}{
					"Id":     dataTerceros[0].TipoDocumentoId.Id,
					"Nombre": dataTerceros[0].TipoDocumentoId.Nombre,
				},
				"Numero": dataTerceros[0].Numero,
			}
			dataFinal["Identificacion"] = dataRecortada
		} else {

			found := len(dataTerceros)
			if found == 1 && dataTerceros[0].Numero == "" {
				found = 0
			}

			s := ""
			if found >= limitDocs {
				s += " (o mas)"
			}

			err := fmt.Errorf("se esperaba UN (único) documento activo registrado para el Tercero con ID: %d, hay: %d%s", dataTercero.Id, found, s)
			logs.Notice(err)
		}

		terceros = append(terceros, dataFinal)
	}
	return
}
