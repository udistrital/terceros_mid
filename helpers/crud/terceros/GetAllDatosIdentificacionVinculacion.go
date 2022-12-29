package terceros

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func GetAllDatosIdentificacionVinculacion(documentos *[]map[string]interface{}, query, vinculaciones string) (
	outputError map[string]interface{}) {

	const funcion = "GetAllDatosIdentificacionVinculacion - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	urlcrud := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion/identificacion?" +
		"query=" + query + "&vinculaciones=" + vinculaciones

	var data interface{}
	if resp, err := request.GetJsonTest(urlcrud, &data); err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}

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

	return
}
