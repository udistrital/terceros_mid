package tipos

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/udistrital/terceros_mid/models"
	e "github.com/udistrital/utils_oas/errorctrl"
	"github.com/udistrital/utils_oas/request"
)

func GetOrdenadores(idTercero int, query string) (terceros []map[string]interface{}, outputError map[string]interface{}) {
	const funcion = "GetOrdenadores - "
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "/GetOrdenadores - Uncaught Error!",
				"err":     err,
				"status":  "500", // Error no manejado!
			}
			panic(outputError)
		}
	}()

	if query != "" {
		err := errors.New("query no implementado")
		return nil, e.Error(funcion+`query != ""`, err, fmt.Sprint(http.StatusNotImplemented))
	}

	// PARTE 1. Traer los ID de los parámetros asociados a funcionarios de planta

	// Los siguientes son los códigos de los registros de la tabla "parametro" de la API
	// de parámetros, cuyo tipo_parámetro_id sea el asociado a Cargos.
	// Específicamente los códigos de parámetros que se asignen a ordenadores de gastos
	codigosParametroCargos := []string{"R", "VA", "VAF", "D", "DI"}
	CodigoTipoParamCargo := "C"
	CodigoAreaTipo := "C"
	parametroCargoID := make(map[string]int)

	var respBody models.RespuestaAPI1Arr
	urlParametros := "http://" + beego.AppConfig.String("parametrosService") + "parametro?limit=-1"
	urlParametros += "&fields=Id,CodigoAbreviacion"
	urlParametros += "&query=Activo:true,TipoParametroId__CodigoAbreviacion:" + CodigoTipoParamCargo
	urlParametros += ",TipoParametroId__AreaTipoId__CodigoAbreviacion:" + CodigoAreaTipo
	// fmt.Println(urlParametros)
	if resp, err := request.GetJsonTest(urlParametros, &respBody); err == nil && resp.StatusCode == 200 {
		// fmt.Printf("Data: %v\n", respBody.Data)

		for _, paramCargos := range respBody.Data {
			// fmt.Printf("Param #%d: %#v\n", k, paramVinculacion)
			codParam := paramCargos["CodigoAbreviacion"]
			// fmt.Printf("codParam (%T): %v\n", codParam, codParam)
			for _, codigoJefe := range codigosParametroCargos {
				if codigoJefe == codParam {
					// fmt.Printf("P=P %v - T(id):%T - v:%f\n", paramVinculacion, paramVinculacion["Id"], paramVinculacion["Id"])
					parametroCargoID[codigoJefe] = int(paramCargos["Id"].(float64))
				}
			}
		}
		// fmt.Printf("ids: %#v\n", parametroCargoID)
	} else {
		if err == nil {
			err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
		}
		logs.Error(err)
		outputError = map[string]interface{}{
			"funcion": "/GetOrdenadores - request.GetJsonTest(urlParametros, &respBody)",
			"err":     err,
			"status":  "502",
		}
		return nil, outputError
	}
	// PARTE 2. Traer los terceros que tengan estos IDs en la tabla vinculacion

	// NOTA: Esta parte se podría mejorar aplicando concurrencia. Vease:
	// https://gobyexample.com/goroutines
	// https://gobyexample.com/waitgroups
	// https://mayurwadekar2.medium.com/concurrency-and-parallelism-in-golang-c8327701fd94
	for _, paramID := range parametroCargoID {

		var vinculaciones []models.Vinculacion
		urlTerceros := "http://" + beego.AppConfig.String("tercerosService") + "vinculacion?limit=-1"
		urlTerceros += "&fields=Id,TerceroPrincipalId,DependenciaId,CargoId"
		urlTerceros += "&query=Activo:true,CargoId:" + fmt.Sprint(paramID)
		if idTercero > 0 {
			urlTerceros += ",TerceroPrincipalId__Id:" + fmt.Sprint(idTercero)
		}
		// fmt.Println(urlTerceros)
		if resp, err := request.GetJsonTest(urlTerceros, &vinculaciones); err == nil && resp.StatusCode == 200 {

			if len(vinculaciones) == 0 || vinculaciones[0].Id == 0 {
				continue
			}
			// fmt.Println("paramId:", paramID, "#vinculaciones: ", len(vinculaciones))

			// Lo siguiente es para que no se vuelva a agregar un tercero
			// cuando el tercero tenga más de una vinculación
			for _, vincul := range vinculaciones {
				add := true
				for _, tercero := range terceros {
					if mTercero := tercero["TerceroPrincipal"].(*models.Tercero); vincul.TerceroPrincipalId.Id == mTercero.Id {
						add = false
						break
					}
				}
				if add {
					terceros = append(terceros, map[string]interface{}{
						"TerceroPrincipal": vincul.TerceroPrincipalId,
						"Cargo":            vincul.CargoId,
						"DependenciaId":    vincul.DependenciaId,
					})
				}
			}
		} else {
			if err == nil {
				err = fmt.Errorf("Undesired status code - Got:%d", resp.StatusCode)
			}
			logs.Error(err)
			outputError = map[string]interface{}{
				"funcion": "/GetOrdenadores - request.GetJsonTest(urlTerceros, &vinculaciones)",
				"err":     err,
				"status":  "502",
			}
			return nil, outputError
		}
		// fmt.Printf("ids: %#v\n", terceros)
	}

	return terceros, nil

}
