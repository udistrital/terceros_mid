package tipos

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	TercerosCrudModels "github.com/udistrital/terceros_crud/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

// GetProveedor trae la lista de proveedores registrados en Terceros, con opcion de filtrar por ID
func GetProveedor(idProveedor int) (terceros []map[string]interface{}, outputError map[string]interface{}) {
	const funcion = "GetProveedor - "
	defer e.ErrorControlFunction(funcion+"Uncaught Error!", fmt.Sprint(http.StatusInternalServerError))

	// PARTE 1. Traer los ID de los tipo_tercero asociados a proveedores
	// Eliminada: Desde que esté registrado en Agora/terceros, debería poder
	// seleccionarse como proveedor. Si se requiere a futuro el reestablecer
	// algun criterio de filtrado, se podrían agregar opciones __in al query

	// PARTE 2 - Traer los terceros
	var tercerosMap []TercerosCrudModels.Tercero
	urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "tercero?limit=-1"
	urlTerceros += "&fields=Id,NombreCompleto"
	urlTerceros += "&query=Activo:true"
	if idProveedor > 0 {
		urlTerceros += ",Id:" + fmt.Sprint(idProveedor)
	}
	// logs.Debug("urlTerceros:", urlTerceros)
	if resp, err := request.GetJsonTest(urlTerceros, &tercerosMap); err == nil && resp.StatusCode == 200 {
		if len(tercerosMap) == 0 || tercerosMap[0].Id == 0 {
			return
		}
	} else {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = e.Error(funcion+"request.GetJsonTest(urlTerceros, &tercerosMap)",
			err, fmt.Sprint(http.StatusBadGateway))
		return
	}
	// formatdata.JsonPrint(tercerosSinMas)

	// PARTE 3: Completar información de identificación, de estar disponible
	for _, dataTercero := range tercerosMap {

		dataFinal := map[string]interface{}{
			"Tercero": map[string]interface{}{
				"Id":             dataTercero.Id,
				"NombreCompleto": dataTercero.NombreCompleto,
			},
		}

		const limitDocs = 10
		var dataTerceros []TercerosCrudModels.DatosIdentificacion
		urlDocTercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion"
		urlDocTercero += "?fields=TipoDocumentoId,Numero"
		urlDocTercero += "&query=Activo:true,TerceroId__Id:" + fmt.Sprint(dataTercero.Id)
		// logs.Debug("urlDocTercero: ", urlDocTercero)
		if resp, err := request.GetJsonTest(urlDocTercero, &dataTerceros); err == nil && resp.StatusCode == 200 {
			// TODO: Retornar los documentos únicos activos que tiene el tercero, en un arreglo o mapeo
			if len(dataTerceros) == 1 && dataTerceros[0].TipoDocumentoId.Id > 0 {
				dataRecortada := map[string]interface{}{
					"TipoDocumentoId": map[string]interface{}{
						"Id":     dataTerceros[0].TipoDocumentoId.Id,
						"Nombre": dataTerceros[0].TipoDocumentoId.Nombre,
					},
					"Numero": dataTerceros[0].Numero,
				}
				dataFinal["Identificacion"] = dataRecortada
			} else {
				s := ""
				if len(dataTerceros) >= limitDocs {
					s += " (o mas)"
				}
				err := fmt.Errorf("se esperaba UN (único) documento activo registrado para el Tercero con ID: %d, hay: %d%s", dataTercero.Id, len(dataTerceros), s)
				logs.Notice(err)
			}
		} else {
			if err == nil {
				err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = e.Error(funcion+"request.GetJsonTest(urlDocTercero, &dataTerceros)",
				err, fmt.Sprint(http.StatusBadGateway))
			return
		}

		terceros = append(terceros, dataFinal)
	}
	// formatdata.JsonPrint(terceros)

	return
}
