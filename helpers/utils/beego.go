// Utils asociados a queries de APIs generadas con Beego

package utils

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	e "github.com/udistrital/utils_oas/errorctrl"
)

// Acepta los parámetros típicos y retorna un objeto url.Values
// al que se le podrían modificar los parámetros antes de
// codificarlos para ser usados en URLs, por ejemplo:
//
//   limit := 10
//   offset := 20
//   params, err := PrepareBeegoQuery("",[],[],[],limit,offset)
//   url := "http://host:puerto/v1/endpoint?" + params.Encode()
//
// Vease también: https://golang.cafe/blog/how-to-url-encode-string-in-golang-example.html
func PrepareBeegoQuery(query string,
	fields, sortby, order []string, limit, offset string) (
	params url.Values, outputError map[string]interface{}) {
	const funcion = "PrepareBeegoQuery - "
	defer e.ErrorControlFunction(funcion+"unhandled error!", fmt.Sprint(http.StatusInternalServerError))

	if len(query) > 0 {
		params.Add("query", query)
	}

	if len(fields) > 0 {
		params.Add("fields", strings.Join(fields, ","))
	}

	if (len(sortby) > 0 || len(order) > 0) && len(sortby) == len(order) {
		params.Add("sortby", strings.Join(sortby, ","))
		params.Add("order", strings.Join(order, ","))
	} else if len(sortby) != len(order) {
		err := errors.New("sortby shall have same length as order")
		return nil, e.Error(funcion+`len(sortby) != len(order)`, err, fmt.Sprint(http.StatusBadRequest))
	}

	if len(limit) > 0 {
		params.Add("limit", limit)
	}
	if len(offset) > 0 {
		params.Add("offset", offset)
	}
	return
}
