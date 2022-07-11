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

func GetTerceros(terceros *[]models.Tercero, query string,
	limit, offset int, fields, sortby, order []string) (
	outputError map[string]interface{}) {
	const funcion = "GetTerceros - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "tercero?"
	params, err := utils.PrepareBeegoQuery(query, fields, sortby, order, limit, offset)
	if err != nil {
		return err
	}
	urlTerceros += params.Encode()
	logs.Debug("urlTerceros:", urlTerceros)
	var data interface{}
	if resp, err := request.GetJsonTest(urlTerceros, &data); err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error("carajoTE:", err)
		outputError = e.Error(funcion+"request.GetJsonTest(urlTerceros, &tercerosMap)",
			err, fmt.Sprint(http.StatusBadGateway))
		return
	}
	if err := formatdata.FillStruct(data, &terceros); err != nil {
		logs.Error(err)
		outputError = e.Error(funcion+"formatdata.FillStruct(data, &terceros)",
			err, fmt.Sprint(http.StatusInternalServerError))
		return
	}
	if len(*terceros) == 0 || (*terceros)[0].Id == 0 {
		*terceros = []models.Tercero{}
	}
	return
}
