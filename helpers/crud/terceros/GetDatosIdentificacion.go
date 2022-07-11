package terceros

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/terceros_crud/models"
	"github.com/udistrital/terceros_mid/helpers/utils"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func GetDatosIdentificacion(documentos *[]models.DatosIdentificacion, query string,
	limit, offset int, fields, sortby, order []string) (
	outputError map[string]interface{}) {
	const funcion = "GetDatosIdentificacion - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	urlDatosIdentificacion := "http://" + beego.AppConfig.String("tercerosService") + "datos_identificacion?"
	params, err := utils.PrepareBeegoQuery(query, fields, sortby, order, limit, offset)
	if err != nil {
		return err
	}
	urlDatosIdentificacion += params.Encode()
	logs.Debug("urlDatosIdentificacion:", urlDatosIdentificacion)
	var data interface{}
	if resp, err := request.GetJsonTest(urlDatosIdentificacion, &data); err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error("carajoDI:", err)
		outputError = e.Error(funcion+"request.GetJsonTest(urlDatosIdentificacion, &tercerosMap)",
			err, fmt.Sprint(http.StatusBadGateway))
		return
	}
	if err := formatdata.FillStruct(data, &documentos); err != nil {
		logs.Error(err)
		outputError = e.Error(funcion+"formatdata.FillStruct(data, &terceros)",
			err, fmt.Sprint(http.StatusInternalServerError))
		return
	}
	if len(*documentos) == 0 || (*documentos)[0].Id == 0 {
		*documentos = []models.DatosIdentificacion{}
	}
	return
}
