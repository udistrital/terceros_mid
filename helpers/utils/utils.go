package utils

import (
	"errors"

	"github.com/udistrital/utils_oas/formatdata"
)

func DesencapsularRespuesta(in, out interface{}) (err error) {
	payloadNames := []string{
		"Data",
		"Body",
	}
	var (
		ok               bool
		apiResponse      interface{}
		capsuledResponse map[string]interface{}
	)
	capsuledResponse, ok = in.(map[string]interface{})
	if !ok {
		return errors.New("not encapsulated body")
	}
	for _, name := range payloadNames {

		if apiResponse, ok = capsuledResponse[name]; ok {
			// Se pudo desencapsular con la opción actual

			if apiResponse == nil || apiResponse == "null" {
				// Si de una vez se concluye que viene vacío, ni molestarse...
				return
			}

			var arrayAttempt []map[string]interface{}
			if err2 := formatdata.FillStruct(apiResponse, &arrayAttempt); err2 == nil {
				// logs.Debug("se espera retornar un arreglo")
				if len(arrayAttempt) == 1 && len(arrayAttempt[0]) == 0 {
					// Si es un arreglo, pero solo viene
					// un objeto vacío [{}], ni molestarse...
					// logs.Debug("falsa alarma, no hay resultados")
					return
				}
			} else {
				// logs.Debug("se espera retornar un objeto")
			}
			return formatdata.FillStruct(apiResponse, &out)
		}
	}
	// Como ultimo intento, intentar llenar el objeto externo
	return formatdata.FillStruct(in, &out)
}
