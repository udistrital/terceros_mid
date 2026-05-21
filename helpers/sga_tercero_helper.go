package helpers

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/utils_oas/request"
)

func UpdateOrCreateInfoComplementaria(tipoInfo string, infoComp map[string]interface{}, idTercero float64) (map[string]interface{}, bool) {
	resp := map[string]interface{}{}
	ok := false

	if infoComp[tipoInfo].(map[string]interface{})["hasId"] != nil {
		idInfComp := infoComp[tipoInfo].(map[string]interface{})["hasId"].(float64)
		var updateInfoComp map[string]interface{}
		errUpdtInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%v", idInfComp), &updateInfoComp)
		if errUpdtInfoComp == nil && updateInfoComp["Status"] != 404 {
			dataToUpdate := infoComp[tipoInfo].(map[string]interface{})["data"].(map[string]interface{})
			updateInfoComp["InfoComplementariaId"] = dataToUpdate

			var updateAnswer map[string]interface{}
			errupdateAnswer := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idInfComp), "PUT", &updateAnswer, updateInfoComp)
			if errupdateAnswer == nil {
				resp = updateAnswer
				ok = true
			}
		}
	} else {
		newInfo := map[string]interface{}{
			"TerceroId":            map[string]interface{}{"Id": idTercero},
			"InfoComplementariaId": infoComp[tipoInfo].(map[string]interface{})["data"].(map[string]interface{}),
			"Activo":               true,
		}
		var createinfo map[string]interface{}
		errCreateInfo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &createinfo, newInfo)
		if errCreateInfo == nil && fmt.Sprintf("%v", createinfo) != "map[]" && createinfo["Id"] != nil {
			resp = createinfo
			ok = true
		}
	}

	return resp, ok
}

func ObtenerInfoComplementariaId(codigo string) (float64, error) {
	var infoComp []map[string]interface{}
	var okTipo bool
	var errTipo error

	// varios intentos para obtener el dato
	for intento := 1; intento <= 2; intento++ {

		errTipo = request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=Activo:true,CodigoAbreviacion:"+codigo, &infoComp)

		if errTipo != nil {
			logs.Error("Error consultando ", codigo, ": ", errTipo)
			continue
		}

		if len(infoComp) == 0 {
			logs.Error("Respuesta vacía para ", codigo)
			continue
		}

		okTipo = true
		break
	}

	if !okTipo {
		return 0, errors.New("no fue posible obtener info_complementaria para " + codigo)
	}

	return infoComp[0]["Id"].(float64), nil
}
