package tipos

import (
	"fmt"
	"net/http"

	TercerosCrudModels "github.com/udistrital/terceros_crud/models"
	TercerosHelper "github.com/udistrital/terceros_mid/helpers/crud/terceros"
	e "github.com/udistrital/utils_oas/errorctrl"
)

// GetProveedor trae la lista de proveedores registrados en Terceros, con opcion de filtrar por ID
func GetProveedor(idProveedor int, query string) (terceros []map[string]interface{}, outputError map[string]interface{}) {
	const funcion = "GetProveedor - "
	defer e.ErrorControlFunction(funcion+"Uncaught Error!", fmt.Sprint(http.StatusInternalServerError))

	// PREPARAR
	const (
		// TODO: Dar soporte a paginación, traer limit y offset como argumentos de esta función
		limit  = -1
		offset = 0

		QueryBaseTerceros   = "Activo:true"
		QueryBaseDocumentos = "Activo:true,TerceroId__Activo:true"
	)
	var (
		fieldsDocumentos    = []string{"Id", "TipoDocumentoId", "Numero"}
		fullQueryDocumentos string
		respuestaDocumentos []TercerosCrudModels.DatosIdentificacion
		empty               = []string{}
		documentosMap       = make(map[int]TercerosCrudModels.DatosIdentificacion)

		// TODO: consultar a la tabla Terceros para traer también terceros sin documentos:
		// fieldsTerceros      = []string{"Id", "NombreCompleto"}
		// fullQueryTerceros   string
		// respuestaTerceros   []TercerosCrudModels.Tercero
		// tercerosMap         = make(map[int]TercerosCrudModels.Tercero)
	)
	if query != "" { // Si se especificó un parámetro de busqueda
		// 1.1 Terceros que coincidan por nombre
		fullQueryDocumentos = QueryBaseDocumentos + ",TerceroId__NombreCompleto__icontains:" + query
		if err := TercerosHelper.GetDatosIdentificacion(&respuestaDocumentos,
			fullQueryDocumentos, limit, offset, fieldsDocumentos, empty, empty); err != nil {
			outputError = err
			return
		}
		for _, v := range respuestaDocumentos {
			documentosMap[v.Id] = v
		}
		// 1.2 Terceros que coincidan por documento
		fullQueryDocumentos = QueryBaseDocumentos + ",Numero__icontains:" + query
		if err := TercerosHelper.GetDatosIdentificacion(&respuestaDocumentos,
			fullQueryDocumentos, limit, offset, fieldsDocumentos, empty, empty); err != nil {
			outputError = err
			return
		}
		for _, v := range respuestaDocumentos {
			documentosMap[v.Id] = v
		}
	} else { // Todos los terceros a no ser que se pase un id
		fullQueryDocumentos = QueryBaseDocumentos
		if idProveedor > 0 {
			fullQueryDocumentos += ",TerceroId__Id:" + fmt.Sprint(idProveedor)
		}
		if err := TercerosHelper.GetDatosIdentificacion(&respuestaDocumentos,
			fullQueryDocumentos, limit, offset, fieldsDocumentos, empty, empty); err != nil {
			outputError = err
			return
		}
		for _, v := range respuestaDocumentos {
			documentosMap[v.Id] = v
		}
	}
	// 2. Procesar y retornar respuesta
	terceros = make([]map[string]interface{}, 0)
	for _, v := range documentosMap {
		terceros = append(terceros, map[string]interface{}{
			"Tercero": map[string]interface{}{
				"Id":             v.TerceroId.Id,
				"NombreCompleto": v.TerceroId.NombreCompleto,
			},
			"Identificacion": map[string]interface{}{
				"TipoDocumentoId": map[string]interface{}{
					"Id":     v.TipoDocumentoId.Id,
					"Nombre": v.TipoDocumentoId.Nombre,
				},
				"Numero": v.Numero,
			},
		})
	}
	return
}
