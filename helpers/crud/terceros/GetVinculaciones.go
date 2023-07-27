package terceros

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/terceros_crud/models"
	"github.com/udistrital/terceros_mid/helpers/utils"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
)

func GetVinculaciones(vinculaciones *[]models.Vinculacion, query string,
	limit, offset int, fields, sortby, order []string) (
	outputError map[string]interface{}) {
	const funcion = "GetVinculaciones - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	urlVinculaciones := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?"
	params, err := utils.PrepareBeegoQuery(query, fields, sortby, order, limit, offset)
	if err != nil {
		return err
	}
	urlVinculaciones += params.Encode()
	var data interface{}
	// logs.Debug("urlVinculaciones:", urlVinculaciones, "- data:", data)
	if resp, err := request.GetJsonTest(urlVinculaciones, &data); err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = e.Error(funcion+"request.GetJsonTest(urlVinculaciones, &tercerosMap)",
			err, fmt.Sprint(http.StatusBadGateway))
		return
	}
	if err := formatdata.FillStruct(data, &vinculaciones); err != nil {
		logs.Error(err)
		outputError = e.Error(funcion+"formatdata.FillStruct(data, &terceros)",
			err, fmt.Sprint(http.StatusInternalServerError))
		return
	}
	if len(*vinculaciones) == 0 || (*vinculaciones)[0].TerceroPrincipalId == nil {
		*vinculaciones = []models.Vinculacion{}
	}
	return
}

func GetTrVinculacionIdentificacion(compuesto, vinculaciones, cargos, dependencias string) (
	vinculaciones_ []map[string]interface{}, outputError map[string]interface{}) {

	const funcion = "GetVinculaciones - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	urlTrVinculaciones := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion/identificacion?"

	parametros := []string{}
	if compuesto != "" {
		parametros = append(parametros, "query="+compuesto)
	}

	if vinculaciones != "" {
		parametros = append(parametros, "vinculaciones="+vinculaciones)
	}

	if cargos != "" {
		parametros = append(parametros, "cargos="+cargos)
	}

	if dependencias != "" {
		parametros = append(parametros, "dependencias="+dependencias)
	}

	payload := strings.Join(parametros, "&")

	var data interface{}

	if resp, err := request.GetJsonTest(urlTrVinculaciones+payload, &data); err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("undesired Status Code: %d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = e.Error(funcion+"request.GetJsonTest(urlTrVinculaciones, &tercerosMap)",
			err, fmt.Sprint(http.StatusBadGateway))
		return
	}
	err := formatdata.FillStruct(data, &vinculaciones_)
	if err != nil {
		logs.Error(err)
		outputError = e.Error(funcion+"formatdata.FillStruct(data, &terceros)",
			err, fmt.Sprint(http.StatusInternalServerError))
	}

	return
}
