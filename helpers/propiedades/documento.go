package propiedades

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/mitchellh/mapstructure"
	"github.com/udistrital/terceros_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func GetDocumento(terceroId string) (documentos []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetDocumento - Uncaught Error!",
				"err":     err,
				"status":  "500",
			}
			panic(outputError)
		}
	}()

	var dataTerceros []map[string]interface{}
	urlDocTercero := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion"
	urlDocTercero += "?sortby=Id&order=desc&fields=TipoDocumentoId,Numero"
	urlDocTercero += "&query=Activo:true,TerceroId__Id:" + fmt.Sprint(terceroId)
	// logs.Debug("urlDocTercero: ", urlDocTercero)
	if resp, err := request.GetJsonTest(urlDocTercero, &dataTerceros); err == nil && resp.StatusCode == 200 {

		if len(dataTerceros) == 1 && len(dataTerceros[0]) > 0 && dataTerceros[0]["Id"] != 0 {

			var dataModel models.DatosIdentificacion
			if err := mapstructure.Decode(dataTerceros[0], &dataModel); err != nil {
				logs.Error(err)
				outputError = map[string]interface{}{
					"funcion": "GetDocumento - mapstructure.Decode(dataTerceros[0], &dataModel)",
					"err":     err,
					"status":  "500",
				}
				return nil, outputError
			}

			dataRecortada := map[string]interface{}{
				// "TipoDocumentoId": dataModel.TipoDocumentoId,
				"TipoDocumentoId": map[string]interface{}{
					"Id":     dataModel.TipoDocumentoId.Id,
					"Nombre": dataModel.TipoDocumentoId.Nombre,
				},
				"Numero": dataModel.Numero,
			}
			documentos = append(documentos, dataRecortada)
			return documentos, nil

		} else {
			err := fmt.Errorf("Hay +/- un documento registrado como Activo para el Tercero con ID: %s", terceroId)
			logs.Warn(err)
		}
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "GetDocumento - request.GetJsonTest(urlDocTercero, &dataTerceros)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	return
}
