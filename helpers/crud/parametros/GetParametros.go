package parametros

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/parametros_crud/models"
	"github.com/udistrital/terceros_mid/helpers/utils"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func GetParametros(parametros *[]models.Parametro, query string,
	limit, offset int, fields, sortby, order []string) (
	outputError map[string]interface{}) {
	const funcion = "GetParametros - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro?"
	params, err := utils.PrepareBeegoQuery(query, fields, sortby, order, limit, offset)
	if err != nil {
		return err
	}
	urlParametros += params.Encode()
	// logs.Debug("urlParametros:", urlParametros)
	var data interface{}
	if resp, err := request.GetJsonTest(urlParametros, &data); err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = e.Error(funcion+"request.GetJsonTest(urlParametros, &data)",
			err, fmt.Sprint(http.StatusBadGateway))
		return
	}
	if err := utils.DesencapsularRespuesta(data, parametros); err != nil {
		logs.Error(err)
		outputError = e.Error(funcion+"utils.DesencapsularRespuesta(data, parametros)",
			err, fmt.Sprint(http.StatusInternalServerError))
		return
	}
	return
}
