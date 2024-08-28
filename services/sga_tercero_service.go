package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/log"

	// "github.com/udistrital/terceros_mid/helpers"
	"github.com/udistrital/terceros_mid/models"
	"github.com/udistrital/utils_oas/formatdata"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func ActualizarPersona(data []byte) (interface{}, error) {
	var body map[string]interface{}
	response := make(map[string]interface{})
	if err := json.Unmarshal(data, &body); err == nil {

		if idTercero, ok := body["Tercero"].(map[string]interface{})["hasId"].(float64); ok {
			var updateTercero map[string]interface{}
			if body["Tercero"].(map[string]interface{})["hasId"] != nil {
				errtercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idTercero), &updateTercero)
				if errtercero == nil && updateTercero["Status"] != 404 {
					dataToUpdate := body["Tercero"].(map[string]interface{})["data"].(map[string]interface{})
					if PrimerNombre, ok := dataToUpdate["PrimerNombre"]; ok {
						updateTercero["PrimerNombre"] = PrimerNombre
					}
					if SegundoNombre, ok := dataToUpdate["SegundoNombre"]; ok {
						updateTercero["SegundoNombre"] = SegundoNombre
					}
					if PrimerApellido, ok := dataToUpdate["PrimerApellido"]; ok {
						updateTercero["PrimerApellido"] = PrimerApellido
					}
					if SegundoApellido, ok := dataToUpdate["SegundoApellido"]; ok {
						updateTercero["SegundoApellido"] = SegundoApellido
					}
					updateTercero["NombreCompleto"] = updateTercero["PrimerNombre"].(string) + " " + updateTercero["SegundoNombre"].(string) + " " + updateTercero["PrimerApellido"].(string) + " " + updateTercero["SegundoApellido"].(string)
					if FechaNacimiento, ok := dataToUpdate["FechaNacimiento"]; ok {
						updateTercero["FechaNacimiento"] = time_bogota.TiempoCorreccionFormato(FechaNacimiento.(string))
					}
					if UsuarioWSO2, ok := dataToUpdate["UsuarioWSO2"]; ok {
						updateTercero["UsuarioWSO2"] = UsuarioWSO2
					}

					var updateTerceroAns map[string]interface{}
					errUpdateTercero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", idTercero), "PUT", &updateTerceroAns, updateTercero)
					if errUpdateTercero == nil {
						response["tercero"] = updateTerceroAns
					} else {
						logs.Error("Error --> ", errUpdateTercero)
						return nil, errors.New(errUpdateTercero.Error())
					}
				} else {
					logs.Error("Error --> ", errtercero)
					return nil, errors.New(errtercero.Error())
				}
			}

			if body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["hasId"] != nil {
				idInfComp := body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["hasId"].(float64)
				var updateInfoComp map[string]interface{}
				errUpdtInfoComp := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%v", idInfComp), &updateInfoComp)
				if errUpdtInfoComp == nil && updateInfoComp["Status"] != 404 {
					updateInfoComp["Dato"] = body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["data"]

					formatdata.JsonPrint(updateInfoComp)

					var updateAnswer map[string]interface{}
					errupdateAnswer := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idInfComp), "PUT", &updateAnswer, updateInfoComp)
					if errupdateAnswer == nil {
						response["telefono"] = updateAnswer
					}
				}
			} else {
				IdTelefono, _ := models.IdInfoCompTercero("10", "TELEFONO")
				ItTel, _ := strconv.ParseFloat(IdTelefono, 64)
				newInfo := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": idTercero},
					"InfoComplementariaId": map[string]interface{}{"Id": ItTel},
					"Dato":                 body["Complementarios"].(map[string]interface{})["Telefono"].(map[string]interface{})["data"],
					"Activo":               true,
				}

				formatdata.JsonPrint(newInfo)
				var createinfo map[string]interface{}
				errCreateInfo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &createinfo, newInfo)
				if errCreateInfo == nil && fmt.Sprintf("%v", createinfo) != "map[]" && createinfo["Id"] != nil {
					response["telefono"] = createinfo
				}
			}
			return response, nil

		} else {
			return nil, errors.New("error del servicio ActualizarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio ActualizarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido" + err.Error())
	}
}

func GuardarPersona(data []byte) (interface{}, error) {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroPost map[string]interface{}

	var paramReq = []string{"PrimerNombre", "SegundoNombre", "PrimerApellido", "SegundoApellido", "FechaNacimiento", "Usuario",
		"TipoIdentificacion", "NumeroIdentificacion", "FechaExpedicion", "Genero", "Telefono"}
	var jsonOk bool = true

	if err := json.Unmarshal(data, &tercero); err == nil && fmt.Sprintf("%v", tercero) != "map[]" {
		for _, key := range paramReq {
			if _, ok := tercero[key]; !ok {
				jsonOk = false
				break
			}
		}
		if jsonOk {
			TipoContribuyenteId := map[string]interface{}{
				"Id": 1,
			}
			guardarpersona := map[string]interface{}{
				"NombreCompleto":      tercero["PrimerNombre"].(string) + " " + tercero["SegundoNombre"].(string) + " " + tercero["PrimerApellido"].(string) + " " + tercero["SegundoApellido"].(string),
				"PrimerNombre":        tercero["PrimerNombre"],
				"SegundoNombre":       tercero["SegundoNombre"],
				"PrimerApellido":      tercero["PrimerApellido"],
				"SegundoApellido":     tercero["SegundoApellido"],
				"FechaNacimiento":     time_bogota.TiempoCorreccionFormato(tercero["FechaNacimiento"].(string)),
				"Activo":              true,
				"TipoContribuyenteId": TipoContribuyenteId, // Persona natural actualmente tiene ese id en el api
				"UsuarioWSO2":         tercero["Usuario"],
			}
			errPersona := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &terceroPost, guardarpersona)

			if errPersona == nil && fmt.Sprintf("%v", terceroPost) != "map[]" && terceroPost["Id"] != nil {
				if terceroPost["Status"] != 400 {
					idTerceroCreado := terceroPost["Id"]
					var identificacion map[string]interface{}

					TipoDocumentoId := map[string]interface{}{
						"Id": tercero["TipoIdentificacion"].(map[string]interface{})["Id"],
					}
					TerceroId := map[string]interface{}{
						"Id": idTerceroCreado,
					}
					identificaciontercero := map[string]interface{}{
						"Numero":          tercero["NumeroIdentificacion"],
						"TipoDocumentoId": TipoDocumentoId,
						"TerceroId":       TerceroId,
						"Activo":          true,
						"FechaExpedicion": time_bogota.TiempoCorreccionFormato(tercero["FechaExpedicion"].(string)),
					}
					errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion", "POST", &identificacion, identificaciontercero)
					if errIdentificacion == nil && fmt.Sprintf("%v", identificacion) != "map[]" && identificacion["Id"] != nil {
						if identificacion["Status"] != 400 {
							IdTelefono, _ := models.IdInfoCompTercero("10", "TELEFONO")
							ItTel, _ := strconv.ParseFloat(IdTelefono, 64)
							newInfo := map[string]interface{}{
								"TerceroId":            TerceroId,
								"InfoComplementariaId": map[string]interface{}{"Id": ItTel},
								"Dato":                 fmt.Sprintf("{\"principal\":%.0f,\"alterno\":null}", tercero["Telefono"]),
								"Activo":               true,
							}
							formatdata.JsonPrint(newInfo)
							var createinfo map[string]interface{}
							errCreateInfo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &createinfo, newInfo)
							if errCreateInfo == nil && fmt.Sprintf("%v", createinfo) != "map[]" && createinfo["Id"] != nil {
								resultado = terceroPost
								resultado["NumeroIdentificacion"] = identificacion["Numero"]
								resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
								resultado["FechaExpedicion"] = identificacion["FechaExpedicion"]
								resultado["TelefonoId"] = createinfo["Id"]

								return resultado, nil
							}
						} else {
							//Si pasa un error borra todo lo creado al momento del registro del documento de identidad
							//var resultado2 map[string]interface{}
							models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]))
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion/%.f", identificacion["Id"]), "DELETE", &resultado2, nil)
							models.SetInactivo(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]))
							//request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
							logs.Error("Error --> ", errIdentificacion)
							return nil, errors.New("error del servicio GuardarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					} else {
						logs.Error("Error --> ", errIdentificacion)
						return nil, errors.New("error del servicio GuardarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errPersona)
					return nil, errors.New("error del servicio GuardarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			} else {
				logs.Error("Error --> ", errPersona)
				return nil, errors.New("error del servicio GuardarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			logs.Error("Error --> ", "Body contains an incorrect data type or an invalid parameter")
			return nil, errors.New("error del servicio GuardarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio GuardarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
	return nil, errors.New("error del servicio GuardarPersona: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func GuardarDatosComplementarios(data []byte) (interface{}, error) {
	var tercero map[string]interface{}     // Body POST
	var HayError bool = false              // Handle Errors
	var resultado []map[string]interface{} // Response POST

	var terceroget map[string]interface{} // Info Tercero for PUT lugar
	var terceroOrg map[string]interface{} // tercero info orig if error
	var LugarPut map[string]interface{}   // resp Put lugar

	if err := json.Unmarshal(data, &tercero); err == nil {

		errtercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", tercero["Tercero"].(float64)), &terceroget)
		if errtercero == nil && terceroget["Status"] != 400 {
			terceroOrg = terceroget
		} else {
			HayError = true
		}

		// Registro de Grupo Sanguineo y Rh
		if !HayError {
			var grupoSanguineoPost map[string]interface{}
			InfoComplementariaId2 := map[string]interface{}{
				"Id": tercero["GrupoSanguineo"],
			}
			grupoSanguineo := map[string]interface{}{
				"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
				"InfoComplementariaId": InfoComplementariaId2,
				"Activo":               true,
			}
			errGrupoSanguineoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &grupoSanguineoPost, grupoSanguineo)
			if errGrupoSanguineoPost == nil && fmt.Sprintf("%v", grupoSanguineoPost) != "map[]" && grupoSanguineoPost["Id"] != nil {
				if grupoSanguineoPost["Status"] != 400 {
					//Ok POST Gr sang
					resultado = append(resultado, grupoSanguineoPost)

					var FactorRhPost map[string]interface{}
					InfoComplementariaId3 := map[string]interface{}{
						"Id": tercero["Rh"],
					}
					factorRh := map[string]interface{}{
						"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
						"InfoComplementariaId": InfoComplementariaId3,
						"Activo":               true,
					}
					errFactorRhPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &FactorRhPost, factorRh)
					if errFactorRhPost == nil && fmt.Sprintf("%v", FactorRhPost) != "map[]" && FactorRhPost["Id"] != nil {
						if FactorRhPost["Status"] != 400 {
							// Ok POST Rh
							resultado = append(resultado, FactorRhPost)

						} else {
							HayError = true
							return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					} else {
						HayError = true
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}

				} else {
					HayError = true
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			} else {
				HayError = true
				return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		}

		// Registro de Poblaciones
		if !HayError {
			poblaciones := tercero["TipoPoblacion"].([]interface{})
			for i := 0; i < len(poblaciones); i++ {
				var poblacionPost1 map[string]interface{}
				TipoPoblacion := poblaciones[i].(map[string]interface{})
				nuevaPoblacion := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": TipoPoblacion["Id"].(float64)},
					"Activo":               true,
				}
				errPoblacionPost1 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &poblacionPost1, nuevaPoblacion)
				if errPoblacionPost1 == nil && fmt.Sprintf("%v", poblacionPost1) != "map[]" && poblacionPost1["Id"] != nil {
					if poblacionPost1["Status"] != 400 {
						//Ok POST select Poblacion
						resultado = append(resultado, poblacionPost1)
					} else {
						HayError = true
						logs.Error("Error --> ", errPoblacionPost1)
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					HayError = true
					logs.Error("Error --> ", errPoblacionPost1)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
			if fmt.Sprintf("%v", reflect.TypeOf(tercero["ComprobantePoblacion"])) == "map[string]interface {}" {
				var poblacionPost2 map[string]interface{}
				comprobantePoblacion := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 315},
					"Activo":               true,
					"Dato":                 `{"value":` + fmt.Sprintf("%v", tercero["ComprobantePoblacion"].(map[string]interface{})["Id"]) + `}`,
				}
				errPoblacionPost2 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &poblacionPost2, comprobantePoblacion)
				if errPoblacionPost2 == nil && fmt.Sprintf("%v", poblacionPost2) != "map[]" && poblacionPost2["Id"] != nil {
					if poblacionPost2["Status"] != 400 {
						//Ok POST comp pobl
						resultado = append(resultado, poblacionPost2)

					} else {
						HayError = true
						logs.Error("Error --> ", errPoblacionPost2)
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					HayError = true
					logs.Error("Error --> ", errPoblacionPost2)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
		}

		// Registro de Lugar de origen
		if !HayError {
			terceroget["LugarOrigen"] = tercero["Lugar"].(map[string]interface{})["Id"].(float64)
			errLugarPut := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", tercero["Tercero"].(float64)), "PUT", &LugarPut, terceroget)
			if errLugarPut == nil && fmt.Sprintf("%v", LugarPut) != "map[]" && LugarPut["Id"] != nil {
				if LugarPut["Status"] != 400 {
					//Ok PUT lugarId tercero
					//resultado = append(resultado, LugarPut)

				} else {
					HayError = true
					logs.Error("Error --> ", errLugarPut)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			} else {
				HayError = true
				logs.Error("Error --> ", errLugarPut)
				return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		}

		// Registro de Discapacidades
		if !HayError {
			discapacidades := tercero["TipoDiscapacidad"].([]interface{})
			for i := 0; i < len(discapacidades); i++ {
				var discapacidadPost1 map[string]interface{}
				discapacidad := discapacidades[i].(map[string]interface{})
				nuevadiscapacidad := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": discapacidad["Id"].(float64)},
					"Activo":               true,
				}
				errDiscapacidadPost1 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &discapacidadPost1, nuevadiscapacidad)
				if errDiscapacidadPost1 == nil && fmt.Sprintf("%v", discapacidadPost1) != "map[]" && discapacidadPost1["Id"] != nil {
					if discapacidadPost1["Status"] != 400 {
						//Ok POST select discapacidad
						resultado = append(resultado, discapacidadPost1)
					} else {
						HayError = true
						logs.Error("Error --> ", errDiscapacidadPost1)
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					HayError = true
					logs.Error("Error --> ", errDiscapacidadPost1)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
			if fmt.Sprintf("%v", reflect.TypeOf(tercero["ComprobanteDiscapacidad"])) == "map[string]interface {}" {
				var discapacidadPost2 map[string]interface{}
				comprobanteDiscapacidad := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 310},
					"Activo":               true,
					"Dato":                 `{"value":` + fmt.Sprintf("%v", tercero["ComprobanteDiscapacidad"].(map[string]interface{})["Id"]) + `}`,
				}
				errDiscapacidadPost2 := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &discapacidadPost2, comprobanteDiscapacidad)
				if errDiscapacidadPost2 == nil && fmt.Sprintf("%v", discapacidadPost2) != "map[]" && discapacidadPost2["Id"] != nil {
					if discapacidadPost2["Status"] != 400 {
						//Ok POST comp disca
						resultado = append(resultado, discapacidadPost2)

					} else {
						HayError = true
						logs.Error("Error --> ", errDiscapacidadPost2)
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					HayError = true
					logs.Error("Error --> ", errDiscapacidadPost2)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
		}

		// Registro de EPS
		if !HayError {
			// Registro de EPS
			if (tercero["EPS"] != nil) && (tercero["FechaVinculacionEPS"] != nil) {

				var postEPS map[string]interface{}
				nuevaEPS := map[string]interface{}{
					"TerceroId":              map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"TerceroEntidadId":       map[string]interface{}{"Id": tercero["EPS"].(map[string]interface{})["Id"].(float64)},
					"FechaInicioVinculacion": tercero["FechaVinculacionEPS"].(string),
					"Activo":                 true,
				}
				errNuevaEPS := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero", "POST", &postEPS, nuevaEPS)
				if errNuevaEPS == nil && fmt.Sprintf("%v", postEPS) != "map[]" && postEPS["Id"] != nil {
					if postEPS["Status"] == 400 {
						HayError = true
						logs.Error("Error --> ", errNuevaEPS)
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					} else {
						resultado = append(resultado, postEPS)
					}
				} else {
					HayError = true
					logs.Error("Error --> ", errNuevaEPS)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
		}

		// Registro de ARL
		if !HayError {
			// Registro de Grupo de sisben
			if tercero["GrupoSisben"] != nil {
				var postGrupoSisben map[string]interface{}
				grSisben := map[string]interface{}{
					"value": fmt.Sprintf("%v", tercero["GrupoSisben"]),
				}
				jsonGrupoSisben, _ := json.Marshal(grSisben)
				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 42},
					"Activo":               true,
					"Dato":                 string(jsonGrupoSisben),
				}
				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postGrupoSisben, nuevoGrupoSisben)
				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postGrupoSisben) != "map[]" && postGrupoSisben["Id"] != nil {
					if postGrupoSisben["Status"] == 400 {
						HayError = true
						logs.Error("Error --> ", errGrupoSisbenPost)
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					} else {
						resultado = append(resultado, postGrupoSisben)
					}
				} else {
					HayError = true
					logs.Error("Error --> ", errGrupoSisbenPost)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
		}

		// Registro de Número de hermanos
		if !HayError {
			// Registro de Número de hermanos
			if tercero["NumeroHermanos"] != nil {
				var postNumeroHermanos map[string]interface{}
				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 319},
					"Activo":               true,
					"Dato":                 fmt.Sprintf("%v", tercero["NumeroHermanos"]),
				}
				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postNumeroHermanos, nuevoGrupoSisben)
				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postNumeroHermanos) != "map[]" && postNumeroHermanos["Id"] != nil {
					if postNumeroHermanos["Status"] == 400 {
						HayError = true
						logs.Error("Error --> ", errGrupoSisbenPost)
						return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					} else {
						resultado = append(resultado, postNumeroHermanos)
					}
				} else {
					HayError = true
					logs.Error("Error --> ", errGrupoSisbenPost)
					return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
		}

		//Registro de estado civil
		if !HayError {
			// DATOS QUE VAN EL BODY DE LA PETICION A INFO_COMPLEMENTARIA_TERCERO
			bodyInfoComplementariaTercero := map[string]interface{}{
				"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
				"InfoComplementariaId": map[string]interface{}{"Id": tercero["EstadoCivil"].(map[string]interface{})["Id"].(float64)},
				"Activo":               true,
			}
			// Esta variable guardara la respuesta del api en caso que sea correcta.
			var responsePostEstadoCivil map[string]interface{}
			errResponsePostEstadoCivil := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &responsePostEstadoCivil, bodyInfoComplementariaTercero)
			// Si la respuesta es correcta, se procede a validar si el estado civil fue guardado correctamente.
			if errResponsePostEstadoCivil == nil && fmt.Sprintf("%v", errResponsePostEstadoCivil) != "map[]" && responsePostEstadoCivil["Id"] != nil {
				if responsePostEstadoCivil["Status"] != 400 {
					//Ok PUT lugarId tercero
					//resultado = append(resultado, LugarPut)

				} else {
					HayError = true
					logs.Error("Error --> ", errResponsePostEstadoCivil)
					return nil, errors.New("service sga_tercero_service: (GuardarDatosComplementarios) error al guardar el estado civil")
				}
			} else {
				HayError = true
				logs.Error("Error --> ", errResponsePostEstadoCivil)
				return nil, errors.New("service sga_tercero_service: (GuardarDatosComplementarios) error al guardar el estado civil")
			}
		}

		//Registro de estado orientacion sexual
		if !HayError {
			// DATOS QUE VAN EL BODY DE LA PETICION A INFO_COMPLEMENTARIA_TERCERO
			bodyInfoComplementariaTercero := map[string]interface{}{
				"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
				"InfoComplementariaId": map[string]interface{}{"Id": tercero["OrientacionSexual"].(map[string]interface{})["Id"].(float64)},
				"Activo":               true,
			}
			// Esta variable guardara la respuesta del api en caso que sea correcta.
			var responsePostOrientacionSexual map[string]interface{}
			errResponsePostOrientacionSexual := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &responsePostOrientacionSexual, bodyInfoComplementariaTercero)
			// Si la respuesta es correcta, se procede a validar si el estado civil fue guardado correctamente.
			if errResponsePostOrientacionSexual == nil && fmt.Sprintf("%v", errResponsePostOrientacionSexual) != "map[]" && responsePostOrientacionSexual["Id"] != nil {
				if responsePostOrientacionSexual["Status"] != 400 {
					resultado = append(resultado, responsePostOrientacionSexual)
				} else {
					HayError = true
					logs.Error("Error --> ", errResponsePostOrientacionSexual)
					return nil, errors.New("service sga_tercero_service: (GuardarDatosComplementarios) error al guardar la orientación sexual")
				}
			} else {
				HayError = true
				logs.Error("Error --> ", errResponsePostOrientacionSexual)
				return nil, errors.New("service sga_tercero_service: (GuardarDatosComplementarios) error al guardar la orientación sexual")
			}
		}

		//Registro de identidad de genero
		if !HayError {
			// DATOS QUE VAN EL BODY DE LA PETICION A INFO_COMPLEMENTARIA_TERCERO
			bodyInfoComplementariaTercero := map[string]interface{}{
				"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
				"InfoComplementariaId": map[string]interface{}{"Id": tercero["IdentidadGenero"].(map[string]interface{})["Id"].(float64)},
				"Activo":               true,
			}
			// Esta variable guardara la respuesta del api en caso que sea correcta.
			var responsePostIdentidadGenero map[string]interface{}
			errResponsePostIdentidadGenero := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &responsePostIdentidadGenero, bodyInfoComplementariaTercero)
			// Si la respuesta es correcta, se procede a validar si el estado civil fue guardado correctamente.
			if errResponsePostIdentidadGenero == nil && fmt.Sprintf("%v", errResponsePostIdentidadGenero) != "map[]" && responsePostIdentidadGenero["Id"] != nil {
				if responsePostIdentidadGenero["Status"] != 400 {
					resultado = append(resultado, responsePostIdentidadGenero)
				} else {
					HayError = true
					logs.Error("Error --> ", errResponsePostIdentidadGenero)
					return nil, errors.New("service sga_tercero_service: (GuardarDatosComplementarios) error al guardar la identidad de genero")
				}
			} else {
				HayError = true
				logs.Error("Error --> ", errResponsePostIdentidadGenero)
				return nil, errors.New("service sga_tercero_service: (GuardarDatosComplementarios) error al guardar la identidad de genero")
			}
		}

	} else {
		HayError = true
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}

	if !HayError { // if all ok, pass response
		resultado = append(resultado, LugarPut)
		return resultado, nil
	} else { // Delete POSTed if error
		for _, infoComp := range resultado {
			var respDel map[string]interface{}
			request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", infoComp["Id"]), "DELETE", &respDel, nil)
		}
		var respPut map[string]interface{} // restore Put data tercero
		request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", terceroOrg["Id"].(float64)), "PUT", &respPut, terceroOrg)
	}
	return nil, errors.New("error del servicio GuardarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func GuardarDatosComplementariosParAcademico(data []byte) (interface{}, error) {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroget map[string]interface{}
	var tercerooriginal map[string]interface{}
	var Area_Conocimiento map[string]interface{}
	var Nivel_Formacion map[string]interface{}
	var Institucionr map[string]interface{}

	if err := json.Unmarshal(data, &tercero); err == nil {

		errtercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/"+fmt.Sprintf("%v", tercero["Tercero"].(map[string]interface{})["Id"]), &terceroget)
		if errtercero == nil && terceroget["Status"] != 400 {
			tercerooriginal = terceroget
		} else {
			logs.Error(errtercero.Error())
			return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [errtercero] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

		Area_ConocimientoTemp := tercero["AreaConocimiento"].(map[string]interface{})["AREA_CONOCIMIENTO"].([]interface{})
		for _, areatemp := range Area_ConocimientoTemp {
			Area_Conocimiento = areatemp.(map[string]interface{})
		}

		var AreaConocimientoPost map[string]interface{}

		//Codifica en un map separado la informacion del area Conocimiento
		AreaConocimiento := map[string]interface{}{
			"AreaConocimiento": tercero["AreaConocimiento"].(map[string]interface{})["AreaConocimiento"],
		}
		//la convierte en json
		jsonAreaConocimientoString, _ := json.Marshal(AreaConocimiento)

		informacionParAcademico := map[string]interface{}{
			"TerceroId":            tercerooriginal,
			"InfoComplementariaId": Area_Conocimiento,
			"Activo":               true,
			"Dato":                 string(jsonAreaConocimientoString),
		}
		formatdata.JsonPrint(informacionParAcademico)
		errAreaConocimientoPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &AreaConocimientoPost, informacionParAcademico)
		if errAreaConocimientoPost == nil && fmt.Sprintf("%v", AreaConocimientoPost) != "map[]" && AreaConocimientoPost["Id"] != nil {
			if AreaConocimientoPost["Status"] != 400 {
				Nivel_FormacionTemp := tercero["FormacionAcademica"].(map[string]interface{})["NIVEL_FORMACION"].([]interface{})
				for _, areatemp := range Nivel_FormacionTemp {
					Nivel_Formacion = areatemp.(map[string]interface{})
				}

				var NivelformacionPost map[string]interface{}

				NivelFormacion := map[string]interface{}{
					"NivelFormacion": tercero["FormacionAcademica"].(map[string]interface{})["FormacionAcademica"],
				}
				jsonNivelFomracion, _ := json.Marshal(NivelFormacion)

				informacionParAcademico2 := map[string]interface{}{
					"TerceroId":            tercerooriginal,
					"InfoComplementariaId": Nivel_Formacion,
					"Activo":               true,
					"Dato":                 string(jsonNivelFomracion),
				}
				errNivelFormacionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &NivelformacionPost, informacionParAcademico2)
				if errNivelFormacionPost == nil && fmt.Sprintf("%v", NivelformacionPost) != "map[]" && NivelformacionPost["Id"] != nil {
					if NivelformacionPost["Status"] != 400 {

						InstucionTemp := tercero["Institucion"].(map[string]interface{})["INSTITUCION"].([]interface{})
						for _, areatemp := range InstucionTemp {
							Institucionr = areatemp.(map[string]interface{})
						}
						var InstitucionPost map[string]interface{}

						Institucion := map[string]interface{}{
							"Institucion": tercero["Institucion"].(map[string]interface{})["Institucion"],
						}
						jsonInstitucion, _ := json.Marshal(Institucion)

						informacionParAcademico3 := map[string]interface{}{
							"TerceroId":            tercerooriginal,
							"InfoComplementariaId": Institucionr,
							"Activo":               true,
							"Dato":                 string(jsonInstitucion),
						}
						errInstitucionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &InstitucionPost, informacionParAcademico3)
						if errInstitucionPost == nil && fmt.Sprintf("%v", InstitucionPost) != "map[]" && InstitucionPost["Id"] != nil {
							if InstitucionPost["Status"] != 400 {

								resultado = tercero
								return resultado, nil
							} else {
								var resultado2 map[string]interface{}
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", NivelformacionPost["Id"]), "DELETE", &resultado2, nil)
								request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", AreaConocimientoPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error("Error --> ", errInstitucionPost)
								return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [errInstitucionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							logs.Error("Error --> ", errInstitucionPost)
							return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [errInstitucionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					} else {
						var resultado2 map[string]interface{}
						request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/"+fmt.Sprintf("%v", AreaConocimientoPost["Id"]), "DELETE", &resultado2, nil)

						logs.Error("Error --> ", errNivelFormacionPost)
						return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [errNivelFormacionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errNivelFormacionPost)
					return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [errNivelFormacionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}

			} else {

				logs.Error("Error --> ", errAreaConocimientoPost)
				return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [errAreaConocimientoPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			logs.Error("Error --> ", errAreaConocimientoPost)
			return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [errAreaConocimientoPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

	} else {
		return nil, errors.New("error del servicio GuardarDatosComplementariosParAcademico: [error] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func ActualizarDatosComplementarios(data []byte) (interface{}, error) {
	//Persona a la cual se van a agregar los datos complementarios
	var persona map[string]interface{}
	//Grupo etnico al que pertenece la persona
	// var GrupoEtnico map[string]interface{}
	var Poblacion map[string]interface{}
	var PoblacionAux []map[string]interface{}
	//Discapacidades que tiene la persona
	var Discapacidad map[string]interface{}
	//Discapacidad = make(map[string]interface{})
	var DiscapacidadAux []map[string]interface{}
	//Grupo sanguineo de la persona
	// var GrupoSanguineo map[string]interface{}
	GrupoSanguineo := make(map[string]interface{})
	// var GrupoRh map[string]interface{}
	GrupoRh := make(map[string]interface{})
	var GrupoSanguineoAux []map[string]interface{}
	var GrupoSAux []map[string]interface{}
	//resultado de la consulta por ente
	var resultado []map[string]interface{}
	var idpersona_rh []map[string]interface{}
	var idpersona_grupo_sanguineo []map[string]interface{}
	var postEPS map[string]interface{}
	var postGrupoSisben map[string]interface{}
	var postNumeroHermanos map[string]interface{}
	//Resultado de agregar grupo sanguineo y discapacidades
	var resultado2 map[string]interface{}
	//Resultado de agregar grupo sanguineo y discapacidades
	var resultado3 map[string]interface{}
	var resultado4 map[string]interface{}
	var resultado5 map[string]interface{}
	var resultado6 map[string]interface{}
	var resultado7 map[string]interface{}
	var resultado8 map[string]interface{}
	//acumulado de errores
	errores := []interface{}{"acumulado de alertas"}

	//comprobar que el JSON de entrada sea correcto
	if err := json.Unmarshal(data, &persona); err == nil {
		errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/?query=Id:"+fmt.Sprintf("%.f", persona["Ente"]), &resultado)
		if errPersona == nil && resultado != nil {

			//GET para traer las poblaciones registradas del tercero
			poblacion := persona["TipoPoblacion"].([]interface{})
			idPersona := persona["Ente"]
			//var auxDeleteP map[string]interface{}
			//var errDeleteP error
			var OkInactive1 bool
			errPoblacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:3&sortby=Id&order=desc&limit=0", &PoblacionAux)
			if errPoblacion == nil {
				if len(PoblacionAux) > 0 {
					for _, registro := range PoblacionAux {
						idPoblacionAux := fmt.Sprintf("%.f", registro["Id"].(float64))
						//errDeleteP = request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idPoblacionAux, "DELETE", &auxDeleteP, nil)
						OkInactive1 = models.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + idPoblacionAux)
					}
				}
				if OkInactive1 {
					for _, poblaciones := range poblacion {
						nuevaPoblacion := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": poblaciones.(map[string]interface{})["Id"].(float64)},
							"Activo":               true,
						}

						errPoblacionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Poblacion, nuevaPoblacion)
						if errPoblacionPost == nil && fmt.Sprintf("%v", Poblacion) != "map[]" && Poblacion["Id"] != nil {
							if Poblacion["Status"] != 400 {

							} else {
								logs.Error("Error --> ", errPoblacionPost)
								return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errPoblacionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							logs.Error("Error --> ", errPoblacionPost)
							return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errPoblacionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					}

					if fmt.Sprintf("%v", reflect.TypeOf(persona["ComprobantePoblacion"])) == "map[string]interface {}" {
						comprobantePoblacion := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 315},
							"Activo":               true,
							"Dato":                 `{"value":` + fmt.Sprintf("%v", persona["ComprobantePoblacion"].(map[string]interface{})["Id"]) + `}`,
						}
						errPoblacionPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Poblacion, comprobantePoblacion)

						if errPoblacionPost == nil && fmt.Sprintf("%v", Poblacion) != "map[]" && Poblacion["Id"] != nil {
							if Poblacion["Status"] != 400 {

							} else {
								logs.Error("Error --> ", errPoblacionPost)
								return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errPoblacionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							logs.Error("Error --> ", errPoblacionPost)
							return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errPoblacionPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					}
				}
			}

			if (persona["GrupoSanguineo"] != nil || persona["GrupoSanguineo"] != 0) && (persona["Rh"] != nil || persona["Rh"] != 0) {
				//GET para obtener toda la informacion del rh
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=Id:"+fmt.Sprintf("%.f", persona["Rh"]), &GrupoSanguineoAux)
				GrupoRh["InfoComplementariaId"] = GrupoSanguineoAux[0]
				GrupoRh["TerceroId"] = resultado[0]
				GrupoRh["Activo"] = true
				idRh := GrupoRh["InfoComplementariaId"].(map[string]interface{})["GrupoInfoComplementariaId"].(map[string]interface{})["Id"]
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:"+fmt.Sprintf("%.f", idRh)+"&sortby=Id&order=desc&limit=1", &idpersona_rh)
				//PUT RH
				errGrupoRh := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idpersona_rh[0]["Id"]), "PUT", &resultado3, GrupoRh)
				if errGrupoRh == nil {
					errores = append(errores, []interface{}{"OK grupo_sanquineo_persona"})
				} else {
					errores = append(errores, []interface{}{"err grupo_sanquineo_persona", errGrupoRh.Error()})
				}
				//GET grupo sanguineo
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?query=Id:"+fmt.Sprintf("%.f", persona["GrupoSanguineo"]), &GrupoSAux)
				GrupoSanguineo["TerceroId"] = resultado[0]
				GrupoSanguineo["InfoComplementariaId"] = GrupoSAux[0]
				GrupoSanguineo["Activo"] = true
				idGrupoSan := GrupoSanguineo["InfoComplementariaId"].(map[string]interface{})["GrupoInfoComplementariaId"].(map[string]interface{})["Id"]
				request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:"+fmt.Sprintf("%.f", idGrupoSan)+"&sortby=Id&order=desc&limit=1", &idpersona_grupo_sanguineo)
				errGrupoSanguineo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idpersona_grupo_sanguineo[0]["Id"]), "PUT", &resultado4, GrupoSanguineo)
				if errGrupoSanguineo == nil {
					errores = append(errores, []interface{}{"OK grupo_sanquineo_persona"})
				} else {
					errores = append(errores, []interface{}{"err grupo_sanquineo_persona", errGrupoSanguineo.Error()})
				}
			} else {
				errores = append(errores, []interface{}{"el grupo sanguineo es incorrecto:", persona["GrupoSanguineo"], persona["Rh"]})
			}

			//GET para traer las discapacidades registradas del tercero
			discapacidad := persona["TipoDiscapacidad"].([]interface{})
			//var auxDelete map[string]interface{}
			//var errDelete error
			var OkInactive2 bool
			errDiscapacidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId__GrupoInfoComplementariaId__Id:1&sortby=Id&order=desc&limit=0", &DiscapacidadAux)
			if errDiscapacidad == nil {
				if len(DiscapacidadAux) > 0 {
					for _, registro := range DiscapacidadAux {
						idDiscapacidadAux := fmt.Sprintf("%.f", registro["Id"].(float64))
						//errDelete = request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+idDiscapacidadAux, "DELETE", &auxDelete, nil)
						OkInactive2 = models.SetInactivo("http://" + beego.AppConfig.String("TercerosService") + "info_complementaria_tercero/" + idDiscapacidadAux)
					}
				}
				if OkInactive2 {
					for _, discapacidades := range discapacidad {
						nuevadiscapacidad := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": discapacidades.(map[string]interface{})["Id"].(float64)},
							"Activo":               true,
						}

						errDiscapacidadPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Discapacidad, nuevadiscapacidad)
						if errDiscapacidadPost == nil && fmt.Sprintf("%v", Discapacidad) != "map[]" && Discapacidad["Id"] != nil {
							if Discapacidad["Status"] != 400 {

							} else {
								logs.Error("Error --> ", errDiscapacidadPost)
								return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errDiscapacidadPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							logs.Error("Error --> ", errDiscapacidadPost)
							return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errDiscapacidadPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					}

					if fmt.Sprintf("%v", reflect.TypeOf(persona["ComprobanteDiscapacidad"])) == "map[string]interface {}" {
						comprobanteDiscapacidad := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 310},
							"Activo":               true,
							"Dato":                 `{"value":` + fmt.Sprintf("%v", persona["ComprobanteDiscapacidad"].(map[string]interface{})["Id"]) + `}`,
						}
						errDiscapacidadPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &Discapacidad, comprobanteDiscapacidad)

						if errDiscapacidadPost == nil && fmt.Sprintf("%v", Discapacidad) != "map[]" && Discapacidad["Id"] != nil {
							if Discapacidad["Status"] != 400 {

							} else {
								logs.Error("Error --> ", errDiscapacidadPost)
								return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errDiscapacidadPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							logs.Error("Error --> ", errDiscapacidadPost)
							return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errDiscapacidadPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					}
				}
			}

			ubicacion := resultado[0]
			ubicacion["LugarOrigen"] = persona["Lugar"].(map[string]interface{})["Id"]
			if errUbicacionEnte := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%.f", ubicacion["Id"]), "PUT", &resultado6, ubicacion); errUbicacionEnte == nil {
				if resultado6["Type"] == "error" {
					errores = append(errores, resultado2["Body"])
				} else {
					errores = append(errores, []interface{}{"OK update ubicacion_ente"})
				}
			}

			if (persona["EPS"] != nil) && (persona["FechaVinculacionEPS"] != nil) {
				var EPS []map[string]interface{}

				nuevaEPS := map[string]interface{}{
					"TerceroId":              map[string]interface{}{"Id": idPersona.(float64)},
					"TerceroEntidadId":       map[string]interface{}{"Id": persona["EPS"].(map[string]interface{})["Id"].(float64)},
					"FechaInicioVinculacion": persona["FechaVinculacionEPS"].(string),
					"Activo":                 true,
				}

				nuevo := true

				errEPS := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero?query=Activo:true,TerceroId.Id:"+fmt.Sprintf("%.f", idPersona), &EPS)
				if errEPS == nil && fmt.Sprintf("%v", EPS) != "[map[]]" {
					if fmt.Sprintf("%v", EPS) != "[map[]]" {
						if EPS[0]["TerceroEntidadId"].(map[string]interface{})["Id"] == nuevaEPS["TerceroEntidadId"].(map[string]interface{})["Id"] && EPS[0]["FechaInicioVinculacion"] == nuevaEPS["FechaInicioVinculacion"] {
							nuevo = false
						}
					}
					if nuevo {
						EPS[0]["Activo"] = false
						if errEPS := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero/"+fmt.Sprintf("%.f", EPS[0]["Id"]), "PUT", &resultado5, EPS[0]); errEPS == nil {
							if resultado6["Type"] == "error" {
								errores = append(errores, resultado5["Body"])
							}
						}
					}
				}

				if nuevo {
					errNuevaEPS := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero", "POST", &postEPS, nuevaEPS)
					if errNuevaEPS == nil && fmt.Sprintf("%v", postEPS) != "map[]" && postEPS["Id"] != nil {
						if postEPS["Status"] == 400 {
							logs.Error("Error --> ", errNuevaEPS)
							return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errNuevaEPS] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					} else {
						logs.Error("Error --> ", errNuevaEPS)
						return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errNuevaEPS] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				}
			}

			if persona["GrupoSisben"] != nil {
				var GrupoSisben []map[string]interface{}

				errGrupoSisben := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId.Id:42,Activo:True&sortby=Id&order=desc&limit=1", &GrupoSisben)
				if errGrupoSisben == nil && fmt.Sprintf("%v", GrupoSisben) != "[map[]]" {
					GrupoSisben[0]["Activo"] = false
					if errGrupoSisben := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", GrupoSisben[0]["Id"]), "PUT", &resultado7, GrupoSisben[0]); errGrupoSisben == nil {
						if resultado7["Type"] == "error" {
							errores = append(errores, resultado7["Body"])
						}
					}
				}

				grSisben := map[string]interface{}{
					"value": fmt.Sprintf("%v", persona["GrupoSisben"]),
				}
				jsonGrupoSisben, _ := json.Marshal(grSisben)

				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 42},
					"Activo":               true,
					"Dato":                 string(jsonGrupoSisben),
				}

				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postGrupoSisben, nuevoGrupoSisben)
				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postGrupoSisben) != "map[]" && postGrupoSisben["Id"] != nil {
					if postGrupoSisben["Status"] == 400 {
						logs.Error("Error --> ", errGrupoSisbenPost)
						return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errGrupoSisbenPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errGrupoSisbenPost)
					return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errGrupoSisbenPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}

			if persona["NumeroHermanos"] != nil {
				var numeroHermanos []map[string]interface{}

				errNumeroHermanos := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%.f", idPersona)+",InfoComplementariaId.Id:319,Activo:True&sortby=Id&order=desc&limit=1", &numeroHermanos)
				if errNumeroHermanos == nil && fmt.Sprintf("%v", numeroHermanos) != "[map[]]" {
					numeroHermanos[0]["Activo"] = false
					if errNumeroHermanos := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", numeroHermanos[0]["Id"]), "PUT", &resultado8, numeroHermanos[0]); errNumeroHermanos == nil {
						if resultado8["Type"] == "error" {
							errores = append(errores, resultado8["Body"])
						}
					}
				}

				nuevoGrupoSisben := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": idPersona.(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 319},
					"Activo":               true,
					"Dato":                 fmt.Sprintf("%v", persona["NumeroHermanos"]),
				}
				errGrupoSisbenPost := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &postNumeroHermanos, nuevoGrupoSisben)

				if errGrupoSisbenPost == nil && fmt.Sprintf("%v", postNumeroHermanos) != "map[]" && postNumeroHermanos["Id"] != nil {
					if postNumeroHermanos["Status"] == 400 {
						logs.Error("Error --> ", errGrupoSisbenPost)
						return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errGrupoSisbenPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errGrupoSisbenPost)
					return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errGrupoSisbenPost] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}

			if persona["EstadoCivil"] != nil {
				var dataEstado []map[string]interface{}
				var dataEstadoPut map[string]interface{}

				//Consulta la existencia de el campo para actualizarlo
				errEstadoData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId:"+fmt.Sprintf("%.f", idPersona)+",Activo:True,InfoComplementariaId__GrupoInfoComplementariaId__Id:2", &dataEstado)
				if errEstadoData == nil && dataEstado[0]["Id"] != nil {
					//Interface para reasignar el Id del estado
					dataPut := map[string]interface{}{
						"Id": persona["EstadoCivil"].(map[string]interface{})["Id"].(float64),
					}
					//Asignacion a la data de envio
					dataEstado[0]["InfoComplementariaId"] = dataPut
					//Actualizacion de la informacion
					errUpdateEstado := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", dataEstado[0]["Id"]), "PUT", &dataEstadoPut, dataEstado[0])
					if errUpdateEstado == nil {
						errores = append(errores, "Estado civil actualizado")
					} else {
						logs.Error("Error --> ", errUpdateEstado)
						return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errUpdateEstado] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errEstadoData)
					return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errEstadoData] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}

			}

			if persona["IdentidadGenero"] != nil {
				var dataIdentidad []map[string]interface{}
				var dataIdentidadPut map[string]interface{}

				//Consulta la existencia de el campo para actualizarlo
				errEstadoData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId:"+fmt.Sprintf("%.f", idPersona)+",Activo:True,InfoComplementariaId__GrupoInfoComplementariaId__Id:1637", &dataIdentidad)
				if errEstadoData == nil && dataIdentidad[0]["Id"] != nil {
					//Interface para reasignar el Id de la identidad
					dataPut := map[string]interface{}{
						"Id": persona["IdentidadGenero"].(map[string]interface{})["Id"].(float64),
					}
					//Asignacion a la data de envio
					dataIdentidad[0]["InfoComplementariaId"] = dataPut
					//Actualizacion de la informacion
					errUpdateEstado := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", dataIdentidad[0]["Id"]), "PUT", &dataIdentidadPut, dataIdentidad[0])
					if errUpdateEstado == nil {
						errores = append(errores, "Identidad actualizada")
					} else {
						logs.Error("Error --> ", errUpdateEstado)
						return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errUpdateEstado] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errEstadoData)
					return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errEstadoData] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}

			if persona["OrientacionSexual"] != nil {
				var dataOrientacion []map[string]interface{}
				var dataOrientacionPut map[string]interface{}

				//Consulta la existencia de el campo para actualizarlo
				errEstadoData := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId:"+fmt.Sprintf("%.f", idPersona)+",Activo:True,InfoComplementariaId__GrupoInfoComplementariaId__Id:1636", &dataOrientacion)
				if errEstadoData == nil && dataOrientacion[0]["Id"] != nil {
					//Interface para reasignar el Id de la orientacion
					dataPut := map[string]interface{}{
						"Id": persona["OrientacionSexual"].(map[string]interface{})["Id"].(float64),
					}
					//Asignacion a la data de envio
					dataOrientacion[0]["InfoComplementariaId"] = dataPut
					//Actualizacion de la informacion
					errUpdateEstado := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", dataOrientacion[0]["Id"]), "PUT", &dataOrientacionPut, dataOrientacion[0])
					if errUpdateEstado == nil {
						errores = append(errores, "Orientacion actualizada")
					} else {
						logs.Error("Error --> ", errUpdateEstado)
						return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errUpdateEstado] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errEstadoData)
					return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errEstadoData] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			}
			return errores, nil
		} else {
			if errPersona != nil {
				return nil, errors.New("error del servicio ActualizarDatosComplementarios: [errPersona] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
			if len(resultado) == 0 {
				return nil, errors.New("error del servicio ActualizarDatosComplementarios: NO existe ninguna persona con este ente")
			}
			return nil, errors.New("error del servicio ActualizarDatosComplementarios: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		logs.Error(err.Error())
		return nil, errors.New("error del servicio ActualizarDatosComplementarios: " + err.Error())
	}
}

func ConsultarExistenciaPersona(numeroDocumento string) (interface{}, error) {
	var resultados []map[string]interface{}

	var documentos []map[string]interface{}
	errDocumentos := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,Numero:"+numeroDocumento+"&sortby=FechaCreacion&order=desc&limit=0", &documentos)
	if errDocumentos == nil && fmt.Sprintf("%v", documentos) != "[map[]]" {
		for _, doc := range documentos {
			preparedoc := doc["TerceroId"].(map[string]interface{})

			IdTercero := fmt.Sprintf("%v", doc["TerceroId"].(map[string]interface{})["Id"])

			preparedoc["NumeroIdentificacion"] = doc["Numero"]
			preparedoc["TipoIdentificacion"] = doc["TipoDocumentoId"]
			preparedoc["FechaExpedicion"] = doc["FechaExpedicion"]
			preparedoc["SoporteDocumento"] = doc["DocumentoSoporte"]
			preparedoc["IdentificacionId"] = doc["Id"]

			var estado []map[string]interface{}
			var genero []map[string]interface{}
			var orientacionSexual []map[string]interface{}
			var identidadGenero []map[string]interface{}
			var telefono []map[string]interface{}

			errEstado := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:2", &estado)
			if errEstado == nil && fmt.Sprintf("%v", estado[0]) != "map[]" {
				if estado[0]["Status"] != 404 {
					preparedoc["EstadoCivil"] = estado[0]["InfoComplementariaId"]
					preparedoc["EstadoCivilId"] = estado[0]["Id"]
				} else {
					if estado[0]["Message"] == "Not found resource" {

					} else {
						logs.Error("Error --> ", estado)
					}
				}
			} else {
				logs.Error("Error --> ", estado)
			}

			errGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:6", &genero)
			if errGenero == nil && fmt.Sprintf("%v", genero[0]) != "map[]" {
				if genero[0]["Status"] != 404 {
					preparedoc["Genero"] = genero[0]["InfoComplementariaId"]
					preparedoc["GeneroId"] = genero[0]["Id"]
				} else {
					if genero[0]["Message"] == "Not found resource" {
					} else {
						logs.Error("Error --> ", genero)
					}
				}
			} else {
				logs.Error("Error --> ", genero)
			}

			errOrientacionSexual := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1636", &orientacionSexual)
			if errOrientacionSexual == nil && fmt.Sprintf("%v", orientacionSexual[0]) != "map[]" {
				if orientacionSexual[0]["Status"] != 404 {
					preparedoc["OrientacionSexual"] = orientacionSexual[0]["InfoComplementariaId"]
					preparedoc["OrientacionSexualId"] = orientacionSexual[0]["Id"]
				} else {
					if orientacionSexual[0]["Message"] == "Not found resource" {
					} else {
						logs.Error("Error --> ", orientacionSexual)
					}
				}
			} else {
				logs.Error("Error --> ", orientacionSexual)
			}

			errIdentidadGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				IdTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1637", &identidadGenero)
			if errIdentidadGenero == nil && fmt.Sprintf("%v", identidadGenero[0]) != "map[]" {
				if identidadGenero[0]["Status"] != 404 {
					preparedoc["IdentidadGenero"] = identidadGenero[0]["InfoComplementariaId"]
					preparedoc["IdentidadGeneroId"] = identidadGenero[0]["Id"]
				} else {
					if identidadGenero[0]["Message"] == "Not found resource" {
					} else {
						logs.Error("Error --> ", identidadGenero)
					}
				}
			} else {
				logs.Error("Error --> ", identidadGenero)
			}

			IdTelefono, _ := models.IdInfoCompTercero("10", "TELEFONO")
			errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Activo:true,TerceroId.Id:"+
				IdTercero+",InfoComplementariaId__Id:"+IdTelefono+"&sortby=Id&order=desc&limit=1", &telefono)
			if errTelefono == nil && fmt.Sprintf("%v", telefono) != "[map[]]" {
				var dataJson map[string]interface{}
				if err := json.Unmarshal([]byte(telefono[0]["Dato"].(string)), &dataJson); err == nil {
					preparedoc["Telefono"] = dataJson["principal"]
					preparedoc["TelefonoAlterno"] = dataJson["alterno"]
					preparedoc["TelefonoId"] = telefono[0]["Id"]
				}
			} else {
				logs.Error("Error --> ", telefono)
			}

			resultados = append(resultados, preparedoc)
		}
		return resultados, nil
	} else {
		logs.Error("Error --> ", documentos)
		return nil, errors.New("error del servicio ConsultarExistenciaPersona: [errDocumentos] La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func ConsultarPersona(idTercero string) (interface{}, error) {
	var resultado map[string]interface{}
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idTercero, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			//formatdata.JsonPrint(persona)

			var identificacion []map[string]interface{}

			errIdentificacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion?query=Activo:true,TerceroId.Id:"+idTercero+",TipoDocumentoId__Id__lt:14&sortby=Id&order=desc&limit=0", &identificacion)
			if errIdentificacion == nil && fmt.Sprintf("%v", identificacion[0]) != "map[]" {
				if identificacion[0]["Status"] != 404 {
					var genero []map[string]interface{}
					var telefono []map[string]interface{}

					resultado = persona[0]
					resultado["NumeroIdentificacion"] = identificacion[0]["Numero"]
					resultado["TipoIdentificacion"] = identificacion[0]["TipoDocumentoId"]
					resultado["FechaExpedicion"] = identificacion[0]["FechaExpedicion"]
					resultado["SoporteDocumento"] = identificacion[0]["DocumentoSoporte"]
					//fmt.Println("Resultado identificacion")
					//formatdata.JsonPrint(resultado)

					errGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:6", &genero)
					if errGenero == nil && fmt.Sprintf("%v", genero[0]) != "map[]" {
						if genero[0]["Status"] != 404 {
							resultado["Genero"] = genero[0]["InfoComplementariaId"]
							resultado["GeneroId"] = genero[0]["Id"]
						} else {
							if genero[0]["Message"] == "Not found resource" {
								logs.Error("Not found resource")
							} else {
								logs.Error("Error --> ", genero)
							}
						}
					} else {
						logs.Error("Error --> ", genero)
					}

					IdTelefono, _ := models.IdInfoCompTercero("10", "TELEFONO")
					errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Activo:true,TerceroId.Id:"+
						fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:"+IdTelefono+"&sortby=Id&order=desc&limit=1", &telefono)
					if errTelefono == nil && fmt.Sprintf("%v", telefono) != "[map[]]" {
						var dataJson map[string]interface{}
						if err := json.Unmarshal([]byte(telefono[0]["Dato"].(string)), &dataJson); err == nil {
							resultado["Telefono"] = dataJson["principal"]
							resultado["TelefonoAlterno"] = dataJson["alterno"]
							resultado["TelefonoId"] = telefono[0]["Id"]
						}
					} else {
						logs.Error("Error --> ", telefono)
					}

					return resultado, nil

				} else {
					if identificacion[0]["Message"] == "Not found resource" {
						return nil, errors.New("identificacion no encontrada")
					} else {
						logs.Error("Error --> ", identificacion)
						return nil, errors.New("fallo al encontrar la identificacion")
					}
				}
			} else {
				logs.Error("Error --> ", identificacion)
				return nil, errors.New("fallo al encontrar la identificacion")
			}
		} else {
			if persona[0]["Message"] == "Not found resource" {
				return nil, errors.New("recurso no encontrado [persona]")
			} else {
				logs.Error("Error --> ", persona)
				return nil, errors.New("error al consultar persona")
			}
		}
	} else {
		logs.Error("Error --> ", persona)
		return nil, errors.New(errPersona.Error())
	}
}

func GuardarDatosContacto(data []byte) (interface{}, error) {
	var resultado map[string]interface{}
	var tercero map[string]interface{}
	var EstratoPost map[string]interface{}

	if err := json.Unmarshal(data, &tercero); err == nil {

		// estrato tercero
		estrato := map[string]interface{}{

			"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
			"InfoComplementariaId": map[string]interface{}{"Id": 41}, // Id para estrato
			"Dato":                 tercero["EstratoTercero"],
			"Activo":               true,
		}
		// formatdata.JsonPrint(estrato)
		errEstrato := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &EstratoPost, estrato)
		if errEstrato == nil && fmt.Sprintf("%v", EstratoPost) != "map[]" && EstratoPost["Id"] != nil {

			if EstratoPost["Status"] != 400 {
				//codigo Postal
				var codigopostalPost map[string]interface{}

				codigo := fmt.Sprintf("%v", tercero["Contactotercero"].(map[string]interface{})["CodigoPostal"])
				requestBod := "{\n    \"Data\": \"" + codigo + "\"\n  }"

				codigopostaltercero := map[string]interface{}{
					"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
					"InfoComplementariaId": map[string]interface{}{"Id": 55}, // Id para codigo postal
					"Dato":                 requestBod,
					"Activo":               true,
				}
				//formatdata.JsonPrint(codigopostaltercero)
				errCodigoPostal := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &codigopostalPost, codigopostaltercero)
				if errCodigoPostal == nil && fmt.Sprintf("%v", codigopostalPost) != "map[]" && codigopostalPost["Id"] != nil {
					if codigopostalPost["Status"] != 400 {
						// Telefono
						var telefonoPost map[string]interface{}

						telefonotercero := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
							"InfoComplementariaId": map[string]interface{}{"Id": 51}, // Id para telefono
							"Dato":                 tercero["Contactotercero"].(map[string]interface{})["Telefono"],
							"Activo":               true,
						}

						errTelefono := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &telefonoPost, telefonotercero)
						if errTelefono == nil && fmt.Sprintf("%v", telefonoPost) != "map[]" && telefonoPost["Id"] != nil {
							if telefonoPost["Status"] != 400 {
								// Telefono alternativo
								var telefonoalternativoPost map[string]interface{}

								telefonoalternativotercero := map[string]interface{}{
									"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
									"InfoComplementariaId": map[string]interface{}{"Id": 52}, // Id para telefono alternativo
									"Dato":                 tercero["Contactotercero"].(map[string]interface{})["TelefonoAlterno"],
									"Activo":               true,
								}

								errTelefonoAlterno := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &telefonoalternativoPost, telefonoalternativotercero)
								if errTelefonoAlterno == nil && fmt.Sprintf("%v", telefonoalternativoPost) != "map[]" && telefonoalternativoPost["Id"] != nil {

									if telefonoalternativotercero["Status"] != 400 {
										// Lugar residencia
										var lugarresidenciaPost map[string]interface{}

										lugarresidenciatercero := map[string]interface{}{
											"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
											"InfoComplementariaId": map[string]interface{}{"Id": 58}, // Id para lugar de residencia
											"Dato":                 fmt.Sprintf("%g", tercero["UbicacionTercero"].(map[string]interface{})["Lugar"].(map[string]interface{})["Id"]),
											"Activo":               true,
										}

										errLugarResidencia := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &lugarresidenciaPost, lugarresidenciatercero)
										if errLugarResidencia == nil && fmt.Sprintf("%v", lugarresidenciaPost) != "map[]" && lugarresidenciaPost["Id"] != nil {
											if lugarresidenciatercero["Status"] != 400 {
												// Direccion de residencia
												var direccionPost map[string]interface{}
												direcion := fmt.Sprintf("%v", tercero["UbicacionTercero"].(map[string]interface{})["Direccion"])
												requestBody := "{\n    \"Data\": \"" + direcion + "\"\n  }"

												direcciontercero := map[string]interface{}{
													"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
													"InfoComplementariaId": map[string]interface{}{"Id": 54}, // Id para direccion de residencia
													"Dato":                 requestBody,
													"Activo":               true,
												}

												errDireccion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &direccionPost, direcciontercero)
												if errDireccion == nil && fmt.Sprintf("%v", direccionPost) != "map[]" && direccionPost["Id"] != nil {
													if direcciontercero["Status"] != 400 {
														// Estrato de quien costea
														var estratoquiencosteaPost map[string]interface{}

														estratoquiencosteatercero := map[string]interface{}{
															"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
															"InfoComplementariaId": map[string]interface{}{"Id": 57}, // Id para estrato de responsable
															"Dato":                 fmt.Sprintf("%v", tercero["EstratoQuienCostea"]),
															"Activo":               true,
														}

														errEstratoResponsable := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &estratoquiencosteaPost, estratoquiencosteatercero)
														if errEstratoResponsable == nil && fmt.Sprintf("%v", estratoquiencosteaPost) != "map[]" && estratoquiencosteaPost["Id"] != nil {
															if estratoquiencosteatercero["Status"] != 400 {
																// Correo electronico tercero
																var correoelectronicoPost map[string]interface{}

																direcion := fmt.Sprintf("%v", tercero["Contactotercero"].(map[string]interface{})["Correo"])
																requestBody1 := "{\n    \"Data\": \"" + direcion + "\"\n  }"

																correoelectronicotercero := map[string]interface{}{
																	"TerceroId":            map[string]interface{}{"Id": tercero["Tercero"].(float64)},
																	"InfoComplementariaId": map[string]interface{}{"Id": 53}, // Id para correo electronico
																	"Dato":                 requestBody1,
																	"Activo":               true,
																}

																errCorreo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero", "POST", &correoelectronicoPost, correoelectronicotercero)
																if errCorreo == nil && fmt.Sprintf("%v", correoelectronicoPost) != "map[]" && correoelectronicoPost["Id"] != nil {
																	if correoelectronicotercero["Status"] != 400 {
																		// Resultado final
																		resultado = tercero
																		return resultado, nil
																	} else {
																		//Si pasa un error borra todo lo creado al momento del registro del correo electronico
																		var resultado2 map[string]interface{}
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", lugarresidenciaPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", direccionPost["Id"]), "DELETE", &resultado2, nil)
																		request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", estratoquiencosteaPost["Id"]), "DELETE", &resultado2, nil)
																		logs.Error("Error --> ", errCorreo)
																		return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																	}
																} else {
																	logs.Error("Error --> ", errCorreo)
																	return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																}

															} else {
																//Si pasa un error borra todo lo creado al momento del registro del estrato de quien costea
																var resultado2 map[string]interface{}
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", lugarresidenciaPost["Id"]), "DELETE", &resultado2, nil)
																request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", direccionPost["Id"]), "DELETE", &resultado2, nil)
																logs.Error("Error --> ", errEstratoResponsable)
																return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
															}
														} else {
															logs.Error("Error --> ", errEstratoResponsable)
															return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
														}

													} else {
														//Si pasa un error borra todo lo creado al momento del registro de la direccion
														var resultado2 map[string]interface{}
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
														request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", lugarresidenciaPost["Id"]), "DELETE", &resultado2, nil)
														logs.Error("Error --> ", errDireccion)
														return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													}
												} else {
													logs.Error("Error --> ", errDireccion)
													return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
												}
											} else {
												//Si pasa un error borra todo lo creado al momento del registro del lugar de residencia
												var resultado2 map[string]interface{}
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)
												request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoalternativoPost["Id"]), "DELETE", &resultado2, nil)
												logs.Error("Error --> ", errLugarResidencia)
												return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											}
										} else {
											logs.Error("Error --> ", errLugarResidencia)
											return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
										}
									} else {
										//Si pasa un error borra todo lo creado al momento del registro del telefono alterno
										var resultado2 map[string]interface{}
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
										request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", telefonoPost["Id"]), "DELETE", &resultado2, nil)

										logs.Error("Error --> ", errTelefonoAlterno)
										return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
									}
								} else {
									logs.Error("Error --> ", errTelefonoAlterno)
									return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
								}
							} else {
								//Si pasa un error borra todo lo creado al momento del registro del telefono
								var resultado2 map[string]interface{}
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
								request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", codigopostalPost["Id"]), "DELETE", &resultado2, nil)
								logs.Error("Error --> ", errTelefono)
								return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							logs.Error("Error --> ", errTelefono)
							return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
						}
					} else {
						//Si pasa un error borra todo lo creado al momento del registro del codigo postal
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/%.f", EstratoPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error("Error --> ", errCodigoPostal)
						return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					logs.Error("Error --> ", errCodigoPostal)
					return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			} else {
				logs.Error("Error --> ", errEstrato)
				return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			logs.Error("Error --> ", errEstrato)
			return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio GuardarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")

	}
}

func ConsultarDatosComplementarios(idTercero string) (interface{}, error) {
	//Id de la persona

	//resultado datos complementarios persona
	respuesta := make(map[string]interface{})
	var resultado map[string]interface{}
	var tercero []map[string]interface{}

	errTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero/?query=Id:"+idTercero, &tercero)

	if errTercero == nil && fmt.Sprintf("%v", tercero[0]) != "map[]" {
		if tercero[0]["Status"] != 404 {

			var poblaciones []map[string]interface{}
			resultado = map[string]interface{}{"Tercero": tercero[0]["Id"]}

			// Bloque para obtener el estado civil
			var estado []map[string]interface{}
			errEstado := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:2", &estado)
			if errEstado == nil && len(estado) > 0 && estado[0]["Status"] != 404 {
				resultado["EstadoCivil"] = estado[0]["InfoComplementariaId"]
				resultado["EstadoCivilId"] = estado[0]["Id"]
			} else if errEstado != nil {
				log.Error(errEstado)
				return nil, errors.New("error obteniendo el estado civil")
			}

			// Bloque para obtener la orientación sexual
			var orientacionSexual []map[string]interface{}
			errOrientacionSexual := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1636", &orientacionSexual)
			if errOrientacionSexual == nil && len(orientacionSexual) > 0 && orientacionSexual[0]["Status"] != 404 {
				resultado["OrientacionSexual"] = orientacionSexual[0]["InfoComplementariaId"]
				resultado["OrientacionSexualId"] = orientacionSexual[0]["Id"]
			} else if errOrientacionSexual != nil {
				log.Error(errOrientacionSexual)
				return nil, errors.New("error obteniendo la orientación sexual")
			}

			// Bloque para obtener la identidad de género
			var identidadGenero []map[string]interface{}
			errIdentidadGenero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+
				fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1637", &identidadGenero)
			if errIdentidadGenero == nil && len(identidadGenero) > 0 && identidadGenero[0]["Status"] != 404 {
				resultado["IdentidadGenero"] = identidadGenero[0]["InfoComplementariaId"]
				resultado["IdentidadGeneroId"] = identidadGenero[0]["Id"]
			} else if errIdentidadGenero != nil {
				log.Error(errIdentidadGenero)
				return nil, errors.New("error obteniendo la identidad de género")
			}

			errPoblacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:3&sortby=Id&order=desc&limit=0", &poblaciones)
			if errPoblacion == nil && fmt.Sprintf("%v", poblaciones[0]) != "map[]" {
				if poblaciones[0]["Status"] != 404 {

					var tipoPoblacion []map[string]interface{}
					for i := 0; i < len(poblaciones); i++ {
						if len(poblaciones) > 0 {
							poblacion := poblaciones[i]["InfoComplementariaId"].(map[string]interface{})
							if poblacion["Nombre"] == "DOCUMENTO_SOPORTE_POBLACION" { // Documento soporte
								var documento map[string]interface{}

								if err := json.Unmarshal([]byte(poblaciones[i]["Dato"].(string)), &documento); err == nil {
									resultado["IdDocumentoPoblacion"] = documento["value"]
								}
							} else {
								tipoPoblacion = append(tipoPoblacion, poblacion)
							}
						}
					}
					resultado["TipoPoblacion"] = tipoPoblacion
					var grupoSanguineo []map[string]interface{}

					errGrupoSanguineo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:7&sortby=Id&order=desc&limit=1", &grupoSanguineo)

					if errGrupoSanguineo == nil && fmt.Sprintf("%v", grupoSanguineo[0]) != "map[]" {
						if grupoSanguineo[0]["Status"] != 404 {

							resultado["GrupoSanguineo"] = grupoSanguineo[0]["InfoComplementariaId"]
							var fatorRHGet []map[string]interface{}
							errFactorRh := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:8&sortby=Id&order=desc&limit=1", &fatorRHGet)
							if errFactorRh == nil && fmt.Sprintf("%v", fatorRHGet[0]) != "map[]" {
								if fatorRHGet[0]["Status"] != 404 {

									resultado["Rh"] = fatorRHGet[0]["InfoComplementariaId"]

									var discapacidades []map[string]interface{}
									errDiscapacidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=Activo:true,terceroId.Id:"+fmt.Sprintf("%v", tercero[0]["Id"])+",InfoComplementariaId.GrupoInfoComplementariaId.Id:1&limit=0", &discapacidades)
									if errDiscapacidad == nil && fmt.Sprintf("%v", discapacidades[0]) != "map[]" {
										if discapacidades[0]["Status"] != 404 {

											var tipoDiscapacidad []map[string]interface{}
											for i := 0; i < len(discapacidades); i++ {
												if len(discapacidades) > 0 {
													discapacidad := discapacidades[i]["InfoComplementariaId"].(map[string]interface{})
													if discapacidad["Nombre"] == "DOCUMENTO_SOPORTE_DISCAPACIDAD" { // Documento soporte
														var documento map[string]interface{}

														if err := json.Unmarshal([]byte(discapacidades[i]["Dato"].(string)), &documento); err == nil {
															resultado["IdDocumentoDiscapacidad"] = documento["value"]
														}
													} else {
														tipoDiscapacidad = append(tipoDiscapacidad, discapacidad)
													}
												}
											}
											resultado["TipoDiscapacidad"] = tipoDiscapacidad

											var ubicacionEnte map[string]interface{}
											errUbicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+idTercero, &ubicacionEnte)

											if errUbicacion == nil && fmt.Sprintf("%v", ubicacionEnte) != "map[]" {
												if ubicacionEnte["Status"] != 404 {
													//Consulta ciudad, departamento y pais
													var lugar map[string]interface{}
													errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"relacion_lugares/jerarquia_lugar/"+fmt.Sprintf("%v", ubicacionEnte["LugarOrigen"]), &lugar)

													if errLugar == nil {
														if lugar["Status"] != 404 {
															ubicacionEnte["Lugar"] = lugar
															resultado["Lugar"] = ubicacionEnte

															var grupoSisben []map[string]interface{}

															errGrupoSisben := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Activo:true,TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:42&sortby=Id&order=desc&limit=1", &grupoSisben)
															if errGrupoSisben == nil && fmt.Sprintf("%v", grupoSisben) != "[map[]]" {
																var grSisben map[string]interface{}

																if err := json.Unmarshal([]byte(grupoSisben[0]["Dato"].(string)), &grSisben); err == nil {
																	resultado["GrupoSisben"] = grSisben["value"]
																}

															}

															var EPS []map[string]interface{}

															errEPS := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"seguridad_social_tercero?query=Activo:true,TerceroId.Id:"+idTercero, &EPS)
															if errEPS == nil && fmt.Sprintf("%v", EPS) != "[map[]]" {
																resultado["EPS"] = EPS[0]["TerceroEntidadId"]
																resultado["FechaVinculacionEPS"] = EPS[0]["FechaInicioVinculacion"]
															}

															var hermanosUnivesidad []map[string]interface{}

															errHermanosUni := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=Activo:true,TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:319&sortby=Id&order=desc&limit=1", &hermanosUnivesidad)
															if errHermanosUni == nil && fmt.Sprintf("%v", hermanosUnivesidad) != "[map[]]" {
																resultado["hermanosUnivesidad"] = hermanosUnivesidad[0]["Dato"]
															}

															respuesta["Data"] = resultado
															return resultado, nil
														} else {
															if lugar["Message"] == "Not found resource" {
																return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
															} else {
																return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
															}
														}
													} else {
														return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													}
												} else {
													if ubicacionEnte["Message"] == "Not found resource" {
														return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													} else {
														log.Error(errUbicacion)
														return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													}
												}
											} else {
												log.Error(errUbicacion)
												return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											}

										} else {
											if discapacidades[0]["Message"] == "Not found resource" {
												return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											} else {
												log.Error(errDiscapacidad)
												return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											}
										}
									} else {
										log.Error(errDiscapacidad)
										return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
									}
								} else {
									if fatorRHGet[0]["Message"] == "Not found resource" {
										return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
									} else {
										log.Error(errFactorRh)
										return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
									}
								}
							} else {
								log.Error(errFactorRh)
								return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							if grupoSanguineo[0]["Message"] == "Not found resource" {
								return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							} else {
								log.Error(errGrupoSanguineo)
								return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						}
					} else {
						log.Error(errGrupoSanguineo)
						return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					if poblaciones[0]["Message"] == "Not found resource" {
						return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					} else {
						log.Error(errPoblacion)
						return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				}
			} else {
				log.Error(errPoblacion)
				return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			if tercero[0]["Message"] == "Not found resource" {
				return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			} else {
				log.Error(errTercero)
				return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		}
	} else {
		log.Error(errTercero)
		return nil, errors.New("error del servicio ConsultarDatosComplementarios:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func ConsultarDatosContacto(idTercero string) (interface{}, error) {
	//resultado datos complementarios persona
	var resultado map[string]interface{}
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero?query=Id:"+idTercero, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {
			var estratotercero []map[string]interface{}
			resultado = map[string]interface{}{"Ente": persona[0]["Ente"], "Persona": persona[0]["Id"]}

			errEstrato := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:41", &estratotercero)
			if errEstrato == nil && fmt.Sprintf("%v", estratotercero[0]) != "map[]" {

				if estratotercero[0]["Status"] != 404 {

					resultado["EstratoTercero"] = estratotercero[0]["Dato"]

					var estratoacudiente []map[string]interface{}

					errEstratoAcudiente := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:57", &estratoacudiente)
					if errEstratoAcudiente == nil && fmt.Sprintf("%v", estratoacudiente[0]) != "map[]" {
						if estratoacudiente[0]["Status"] != 404 {
							var CodigoPostal []map[string]interface{}
							resultado["EstratoAcudiente"] = estratoacudiente[0]["Dato"]

							errCodigoPostal := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:55", &CodigoPostal)
							if errCodigoPostal == nil && fmt.Sprintf("%v", CodigoPostal[0]) != "map[]" {
								if CodigoPostal[0]["Status"] != 404 {
									var lugar map[string]interface{}
									resultado["CodigoPostal"] = CodigoPostal[0]["Dato"]

									var Telefono []map[string]interface{}
									errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:51", &Telefono)
									if errTelefono == nil && fmt.Sprintf("%v", Telefono[0]) != "map[]" {
										if Telefono[0]["Status"] != 404 {
											resultado["Telefono"] = Telefono[0]["Dato"]

											var TelefonoAlterno []map[string]interface{}
											errTelefonoAlterno := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:52", &TelefonoAlterno)
											if errTelefonoAlterno == nil && fmt.Sprintf("%v", TelefonoAlterno[0]) != "map[]" {
												if TelefonoAlterno[0]["Status"] != 404 {
													resultado["TelefonoAlterno"] = TelefonoAlterno[0]["Dato"]

													var Direccion []map[string]interface{}
													errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:54", &Direccion)
													if errDireccion == nil && fmt.Sprintf("%v", Direccion[0]) != "map[]" {
														if Direccion[0]["Status"] != 404 {
															resultado["Direccion"] = Direccion[0]["Dato"]

															var Correo []map[string]interface{}
															errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:53", &Correo)
															if errCorreo == nil && fmt.Sprintf("%v", Correo[0]) != "map[]" {
																if Correo[0]["Status"] != 404 {
																	resultado["Correo"] = Correo[0]["Dato"]

																	var ubicacionEnte []map[string]interface{}
																	errUbicacion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.Id:58", &ubicacionEnte)
																	if errUbicacion == nil && fmt.Sprintf("%v", ubicacionEnte[0]) != "map[]" {
																		if ubicacionEnte[0]["Status"] != 404 {

																			errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+
																				fmt.Sprintf("%v", ubicacionEnte[0]["Dato"]), &lugar)
																			if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
																				if lugar["Status"] != 404 {
																					ubicacionEnte[0]["Lugar"] = lugar
																					resultado["UbicacionEnte"] = ubicacionEnte[0]
																					return resultado, nil
																				} else {
																					if lugar["Message"] == "Not found resource" {
																						return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																					} else {
																						logs.Error("Error --> ", lugar)
																						return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																					}
																				}
																			} else {
																				logs.Error("Error --> ", lugar)
																				return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																			}

																		} else {
																			if ubicacionEnte[0]["Message"] == "Not found resource" {
																				return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																			} else {
																				logs.Error("Error --> ", ubicacionEnte)
																				return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																			}
																		}
																	} else {
																		logs.Error("Error --> ", ubicacionEnte)
																		return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																	}
																} else {
																	if Correo[0]["Message"] == "Not found resource" {
																		return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																	} else {
																		logs.Error("Error --> ", Correo)
																		return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																	}
																}
															} else {
																logs.Error("Error --> ", Correo)
																return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
															}
														} else {
															if Direccion[0]["Message"] == "Not found resource" {
																return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
															} else {
																logs.Error("Error --> ", Direccion)
																return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
															}
														}
													} else {
														logs.Error("Error --> ", Direccion)
														return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													}

												} else {
													if TelefonoAlterno[0]["Message"] == "Not found resource" {
														return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													} else {
														logs.Error("Error --> ", TelefonoAlterno)
														return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													}
												}
											} else {
												logs.Error("Error --> ", TelefonoAlterno)
												return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											}

										} else {
											if Telefono[0]["Message"] == "Not found resource" {
												return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											} else {
												logs.Error("Error --> ", Telefono)
												return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											}
										}
									} else {
										logs.Error("Error --> ", Telefono)
										return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
									}
								} else {
									if CodigoPostal[0]["Message"] == "Not found resource" {
										return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
									} else {
										logs.Error("Error --> ", CodigoPostal)
										return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
									}
								}
							} else {
								logs.Error("Error --> ", errCodigoPostal.Error())
								return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						} else {
							if estratoacudiente[0]["Message"] == "Not found resource" {
								return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							} else {
								logs.Error("Error --> ", estratoacudiente)
								return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						}
					} else {
						logs.Error("Error --> ", errEstratoAcudiente)
						return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					if estratotercero[0]["Message"] == "Not found resource" {
						return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					} else {
						logs.Error("Error --> ", estratotercero)
						return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				}
			} else {
				logs.Error("Error --> ", errEstrato)
				return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			if persona[0]["Message"] == "Not found resource" {
				return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			} else {
				logs.Error("Error --> ", persona)
				return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		}
	} else {
		logs.Error("Error --> ", errPersona)
		return nil, errors.New("error del servicio ConsultarDatosContacto:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func ConsultarDatosFamiliar(idTercero string) (interface{}, error) {
	resultado := make(map[string]interface{})
	var terceros []map[string]interface{}
	var correos []map[string]interface{}
	var telefonos []map[string]interface{}
	var direcciones []map[string]interface{}
	var errorGetAll bool

	errTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/?query=TerceroId__Id:"+idTercero+"&sortby=Id&order=asc&limit=0", &terceros)
	if errTercero == nil {
		if terceros != nil {
			if fmt.Sprintf("%v", terceros[0]) != "map[]" {
				resultado["NombreFamiliarPrincipal"] = terceros[0]["TerceroFamiliarId"].(map[string]interface{})["NombreCompleto"]
				resultado["NombreFamiliarAlterno"] = terceros[1]["TerceroFamiliarId"].(map[string]interface{})["NombreCompleto"]

				idPrincipal := fmt.Sprintf("%.f", terceros[0]["TerceroFamiliarId"].(map[string]interface{})["Id"])
				idAlterno := fmt.Sprintf("%.f", terceros[1]["TerceroFamiliarId"].(map[string]interface{})["Id"])

				// GET de correos
				//Correo principal
				errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idPrincipal+",InfoComplementariaId__Id:53", &correos)
				if errCorreo == nil {
					if correos != nil {
						var CorreoJson map[string]interface{}
						if err := json.Unmarshal([]byte(correos[0]["Dato"].(string)), &CorreoJson); err != nil {
							resultado["CorreoElectronico"] = nil
						} else {
							resultado["CorreoElectronico"] = CorreoJson["value"]
							//Correo alterno
							errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idAlterno+",InfoComplementariaId__Id:53", &correos)
							if errCorreo == nil {
								if correos != nil {
									if err := json.Unmarshal([]byte(correos[0]["Dato"].(string)), &CorreoJson); err != nil {
										resultado["CorreoElectronicoAlterno"] = nil
									} else {
										resultado["CorreoElectronicoAlterno"] = CorreoJson["value"]

										//GET Telefono
										//Telefono principal
										errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idPrincipal+",InfoComplementariaId__Id:51", &telefonos)
										if errTelefono == nil {
											if telefonos != nil {
												var TelefonoJson map[string]interface{}
												if err := json.Unmarshal([]byte(telefonos[0]["Dato"].(string)), &TelefonoJson); err != nil {
													resultado["Telefono"] = nil
												} else {
													resultado["Telefono"] = TelefonoJson["value"]
													//Telefono alterno
													errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idAlterno+",InfoComplementariaId__Id:51", &telefonos)
													if errTelefono == nil {
														if telefonos != nil {
															if err := json.Unmarshal([]byte(telefonos[0]["Dato"].(string)), &TelefonoJson); err != nil {
																resultado["TelefonoAlterno"] = nil
															} else {
																resultado["TelefonoAlterno"] = TelefonoJson["value"]

																//GET Direcciones
																//Direccion principal
																errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idPrincipal+",InfoComplementariaId__Id:54", &direcciones)
																if errDireccion == nil {
																	if direcciones != nil {
																		var DireccionJson map[string]interface{}
																		if err := json.Unmarshal([]byte(direcciones[0]["Dato"].(string)), &DireccionJson); err != nil {
																			resultado["DireccionResidencia"] = nil
																		} else {
																			resultado["DireccionResidencia"] = DireccionJson["value"]
																			errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idAlterno+",InfoComplementariaId__Id:54", &direcciones)
																			if errDireccion == nil {
																				if direcciones != nil {
																					if err := json.Unmarshal([]byte(direcciones[0]["Dato"].(string)), &DireccionJson); err != nil {
																						resultado["DireccionResidenciaAlterno"] = nil
																					} else {
																						resultado["DireccionResidenciaAlterno"] = DireccionJson["value"]
																						resultado["Parentesco"] = map[string]interface{}{
																							"Id":     terceros[0]["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
																							"Nombre": terceros[0]["TipoParentescoId"].(map[string]interface{})["Nombre"],
																						}
																						resultado["ParentescoAlterno"] = map[string]interface{}{
																							"Id":     terceros[1]["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
																							"Nombre": terceros[1]["TipoParentescoId"].(map[string]interface{})["Nombre"],
																						}
																					}
																				} else {
																					errorGetAll = true
																					return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																				}
																			} else {
																				errorGetAll = true
																				log.Error(errDireccion.Error())
																				return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																			}
																		}
																	} else {
																		errorGetAll = true
																		return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																	}
																} else {
																	errorGetAll = true
																	log.Error(errDireccion.Error())
																	return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
																}
															}
														} else {
															errorGetAll = true
															return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
														}
													} else {
														errorGetAll = true
														log.Error(errTelefono.Error())
														return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
													}
												}
											} else {
												errorGetAll = true
												return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
											}
										} else {
											errorGetAll = true
											log.Error(errTelefono.Error())
											return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
										}
									}
								} else {
									errorGetAll = true
									return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
								}
							} else {
								errorGetAll = true
								log.Error(errCorreo.Error())
								return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						}
					} else {
						errorGetAll = true
						return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					errorGetAll = true
					log.Error(errCorreo.Error())
					return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
				}
			} else {
				errorGetAll = true
				return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		} else {
			errorGetAll = true
			return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		switch {
		case errTercero != nil:
			return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		case len(terceros) == 0:
			return nil, errors.New("error del servicio ConsultarDatosFamiliar:   No existen familiares asociados a esta persona")
		default:
			return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	}

	if !errorGetAll {
		return resultado, nil
	}

	return nil, errors.New("error del servicio ConsultarDatosFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func ConsultarDatosFormacionPregrado(idTercero string) (interface{}, error) {
	var resultado map[string]interface{}
	var personaInscrita []map[string]interface{}
	var IdColegioGet float64
	resultado = make(map[string]interface{})
	var errorGetAll bool

	errPersona := request.GetJson("http://"+beego.AppConfig.String("InscripcionService")+"/inscripcion_pregrado?query=Activo:true,InscripcionId.PersonaId:"+idTercero, &personaInscrita)
	if errPersona == nil {
		if fmt.Sprintf("%v", personaInscrita[0]) != "map[]" {
			resultado = map[string]interface{}{"Persona Inscrita": personaInscrita[0]}
			resultado["TipoIcfes"] = personaInscrita[0]["TipoIcfesId"]
			resultado["NúmeroRegistroIcfes"] = personaInscrita[0]["CodigoIcfes"]
			resultado["Valido"] = personaInscrita[0]["Valido"]

			var NumeroSemestre []map[string]interface{}
			errNumeroSemestre := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero/?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.Id:14&sortby=FechaCreacion&order=desc&limit=1", &NumeroSemestre)
			if errNumeroSemestre == nil && fmt.Sprintf("%v", NumeroSemestre[0]) != "map[]" {
				if NumeroSemestre[0]["Status"] != 404 {
					resultado["numeroSemestres"] = NumeroSemestre[0]
					//cargar id colegio relacionado
					var IdColegio []map[string]interface{}

					errIdColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+idTercero+",InfoComplementariaId__Id:313,Activo:true&limit=0&sortby=FechaCreacion&order=desc", &IdColegio)
					if errIdColegio == nil {
						if fmt.Sprintf("%v", IdColegio[0]) != "map[]" {
							var formacion map[string]interface{}

							for i := 0; i < len(IdColegio); i++ {
								if err := json.Unmarshal([]byte(IdColegio[i]["Dato"].(string)), &formacion); err == nil {
									if formacion["ProgramaAcademico"] == "colegio" {
										IdColegioGet = (formacion["NitUniversidad"].(map[string]interface{})["Id"]).(float64)

										// Cargar Direccion
										var direccionColegio []map[string]interface{}
										errLugarColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId:"+fmt.Sprintf("%v", IdColegioGet)+",InfoComplementariaId:54", &direccionColegio)
										if errLugarColegio == nil && fmt.Sprintf("%v", direccionColegio[0]) != "map[]" {
											if direccionColegio[0]["Status"] != 404 {
												var direccion map[string]interface{}

												if err := json.Unmarshal([]byte(direccionColegio[0]["Dato"].(string)), &direccion); err == nil {
													resultado["DireccionColegio"] = direccion["DIRECCIÓN"]

												}
											} else {
												errorGetAll = true
											}
										} else {
											errorGetAll = true
											log.Error(errLugarColegio)
										}

										//cargar id Lugar colegio
										var IdLugarColegio []map[string]interface{}

										var jsondata map[string]interface{}
										errIdLugarColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?query=TerceroId:"+fmt.Sprintf("%v", IdColegioGet)+",InfoComplementariaId:89", &IdLugarColegio)
										if errIdLugarColegio == nil && fmt.Sprintf("%v", IdLugarColegio[0]) != "map[]" {
											if IdLugarColegio[0]["Status"] != 404 {

												IdString := IdLugarColegio[0]["Dato"]
												if _, err := strconv.ParseInt(IdString.(string), 10, 64); err == nil {
													jsondata = map[string]interface{}{"dato": IdString}

												} else {

													if err := json.Unmarshal([]byte(IdString.(string)), &jsondata); err != nil {
														panic(err)
													}
												}

												var lugar map[string]interface{}

												errLugar := request.GetJson("http://"+beego.AppConfig.String("UbicacionesService")+"/relacion_lugares/jerarquia_lugar/"+
													fmt.Sprintf("%v", jsondata["dato"]), &lugar)
												if errLugar == nil && fmt.Sprintf("%v", lugar) != "map[]" {
													if lugar["Status"] != 404 {

														resultado["Lugar"] = lugar

														var colegio []map[string]interface{}

														errcolegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/tercero_tipo_tercero?query=TerceroId:"+
															fmt.Sprintf("%v", IdColegioGet), &colegio)
														if errcolegio == nil && fmt.Sprintf("%v", colegio[0]) != "map[]" {
															if colegio[0]["Status"] != 404 {
																resultado["TipoColegio"] = colegio[0]["TipoTerceroId"].(map[string]interface{})["Id"]
																resultado["Colegio"] = colegio[0]["TerceroId"]

															} else {
																if colegio[0]["Message"] == "Not found resource" {

																} else {
																	errorGetAll = true
																	log.Error(colegio)

																}
															}
														} else {
															errorGetAll = true
															log.Error(errcolegio)

														}
													} else {
														if lugar["Message"] == "Not found resource" {

														} else {
															errorGetAll = true
															log.Error(lugar)

														}
													}
												} else {
													errorGetAll = true
													log.Error(errLugar)

												}

											} else {
												if IdLugarColegio[0]["Message"] == "Not found resource" {

												} else {
													errorGetAll = true

												}
											}
										} else {
											errorGetAll = true

										}

										break
									}
								}
							}

						} else {
							if IdColegio[0]["Message"] == "Not found resource" {
								return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							} else {
								errorGetAll = true
								return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
							}
						}
					} else {
						errorGetAll = true
						return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				} else {
					if NumeroSemestre[0]["Message"] == "Not found resource" {
						return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					} else {
						errorGetAll = true
						return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
					}
				}
			} else {
				errorGetAll = true
				return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}

		} else {
			if personaInscrita[0]["Message"] == "Not found resource" {
				return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			} else {
				errorGetAll = true
				return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
			}
		}
	} else {
		errorGetAll = true
		return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}

	if !errorGetAll {
		return resultado, nil
	}
	return nil, errors.New("error del servicio ConsultarDatosFormacionPregrado:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func ActualizarInfoFamiliar(data []byte) (interface{}, error) {
	var InfoFamiliar map[string]interface{}
	var Familiares []map[string]interface{}
	var ParentescoPut map[string]interface{}
	var Telefono []map[string]interface{}
	var TelefonoPut map[string]interface{}
	var Correo []map[string]interface{}
	var CorreoPut map[string]interface{}
	var Direccion []map[string]interface{}
	var DireccionPut map[string]interface{}
	resultado := make(map[string]interface{})
	var errorGetAll bool

	if err := json.Unmarshal(data, &InfoFamiliar); err == nil {
		Familiar := InfoFamiliar["Familiares"].([]interface{})
		IdTercero := fmt.Sprintf("%.f", InfoFamiliar["Tercero_Familiar"].(map[string]interface{})["Id"])

		//GET para traer el id de los familiares asociados al tercero
		errFamiliares := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar?query=TerceroId__Id:"+IdTercero, &Familiares)
		if errFamiliares == nil {
			if Familiares != nil {
				idPrincipal := Familiares[0]["TerceroFamiliarId"].(map[string]interface{})["Id"]
				idAlterno := Familiares[1]["TerceroFamiliarId"].(map[string]interface{})["Id"]

				//Almacena la informacion de contacto del familiar
				ContactoPrincipal := Familiar[0].(map[string]interface{})["InformacionContacto"].([]interface{})
				ContactoAlterno := Familiar[1].(map[string]interface{})["InformacionContacto"].([]interface{})

				//PUT Parentesco
				// Familiar principal
				ParentescoPrincipal := Familiar[0].(map[string]interface{})["Familiar"].(map[string]interface{})["TipoParentescoId"]
				Familiares[0]["TipoParentescoId"] = ParentescoPrincipal
				errParentesco := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/"+fmt.Sprintf("%.f", Familiares[0]["Id"]), "PUT", &ParentescoPut, Familiares[0])
				if errParentesco == nil {
					if ParentescoPut != nil {
						resultado["Parentesco"] = map[string]interface{}{
							"Id":     ParentescoPut["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
							"Nombre": ParentescoPut["TipoParentescoId"].(map[string]interface{})["Nombre"],
						}
						//Familiar alterno
						ParentescoAlterno := Familiar[1].(map[string]interface{})["Familiar"].(map[string]interface{})["TipoParentescoId"]
						Familiares[1]["TipoParentescoId"] = ParentescoAlterno
						errParentesco := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero_familiar/"+fmt.Sprintf("%.f", Familiares[1]["Id"]), "PUT", &ParentescoPut, Familiares[1])
						if errParentesco == nil {
							if ParentescoPut != nil {
								resultado["ParentescoAlterno"] = map[string]interface{}{
									"Id":     ParentescoPut["TipoParentescoId"].(map[string]interface{})["Id"].(float64),
									"Nombre": ParentescoPut["TipoParentescoId"].(map[string]interface{})["Nombre"],
								}
							} else {
								errorGetAll = true
							}
						} else {
							errorGetAll = true
						}
					} else {
						errorGetAll = true
					}
				} else {
					errorGetAll = true
				}

				//PUT Telefono (Info complementaria 51)
				// Familiar Principal
				errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPrincipal)+",InfoComplementariaId__Id:51", &Telefono)
				if errTelefono == nil {
					if Telefono != nil {
						Telefono[0]["Dato"] = ContactoPrincipal[0].(map[string]interface{})["Dato"]
						errTelefono := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Telefono[0]["Id"]), "PUT", &TelefonoPut, Telefono[0])
						if errTelefono == nil {
							if TelefonoPut != nil {
								resultado["Telefono"] = TelefonoPut["Dato"]
								// Familiar alterno
								errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idAlterno)+",InfoComplementariaId__Id:51", &Telefono)
								if errTelefono == nil {
									if Telefono != nil {
										Telefono[0]["Dato"] = ContactoAlterno[0].(map[string]interface{})["Dato"]
										errTelefono := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Telefono[0]["Id"]), "PUT", &TelefonoPut, Telefono[0])
										if errTelefono == nil {
											if TelefonoPut != nil {
												resultado["TelefonoAlterno"] = TelefonoPut["Dato"]
											} else {
												errorGetAll = true
											}
										} else {
											errorGetAll = true
										}
									} else {
										errorGetAll = true
									}
								} else {
									errorGetAll = true
								}
							} else {
								errorGetAll = true
							}
						} else {
							errorGetAll = true
						}
					} else {
						errorGetAll = true
					}
				} else {
					errorGetAll = true
				}

				//PUT Correo (Info complementaria 53)
				// Correo principal
				errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPrincipal)+",InfoComplementariaId__Id:53", &Correo)
				if errCorreo == nil {
					if Correo != nil {
						Correo[0]["Dato"] = ContactoPrincipal[1].(map[string]interface{})["Dato"]
						errCorreo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Correo[0]["Id"]), "PUT", &CorreoPut, Correo[0])
						if errCorreo == nil {
							if Correo != nil {
								resultado["Correo"] = CorreoPut["Dato"]
								// Correo alterno
								errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idAlterno)+",InfoComplementariaId__Id:53", &Correo)
								if errCorreo == nil {
									if Correo != nil {
										Correo[0]["Dato"] = ContactoAlterno[1].(map[string]interface{})["Dato"]
										errCorreo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Correo[0]["Id"]), "PUT", &CorreoPut, Correo[0])
										if errCorreo == nil {
											if Correo != nil {
												resultado["CorreoAlterno"] = CorreoPut["Dato"]
											} else {
												errorGetAll = true
											}
										} else {
											errorGetAll = true
										}
									} else {
										errorGetAll = true
									}
								} else {
									errorGetAll = true
								}
							} else {
								errorGetAll = true
							}
						} else {
							errorGetAll = true
						}
					} else {
						errorGetAll = true
					}
				} else {
					errorGetAll = true
				}

				// PUT Direccion (Info complementaria 54)
				//Direccion principal
				errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idPrincipal)+",InfoComplementariaId__Id:54", &Direccion)
				if errDireccion == nil {
					if Direccion != nil {
						Direccion[0]["Dato"] = ContactoPrincipal[2].(map[string]interface{})["Dato"]
						errDireccion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Direccion[0]["Id"]), "PUT", &DireccionPut, Direccion[0])
						if errDireccion == nil {
							if DireccionPut != nil {
								resultado["Direccion"] = DireccionPut["Dato"]
								//Direccion alterna
								errDireccion := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId__Id:"+fmt.Sprintf("%.f", idAlterno)+",InfoComplementariaId__Id:54", &Direccion)
								if errDireccion == nil {
									if Direccion != nil {
										Direccion[0]["Dato"] = ContactoAlterno[2].(map[string]interface{})["Dato"]
										errDireccion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", Direccion[0]["Id"]), "PUT", &DireccionPut, Direccion[0])
										if errDireccion == nil {
											if DireccionPut != nil {
												resultado["DireccionAlterno"] = DireccionPut["Dato"]
											} else {
												errorGetAll = true
											}
										} else {
											errorGetAll = true
										}
									} else {
										errorGetAll = true
									}
								} else {
									errorGetAll = true
								}
							} else {
								errorGetAll = true
							}
						} else {
							errorGetAll = true
						}
					} else {
						errorGetAll = true
					}
				} else {
					errorGetAll = true
				}

			} else {
				errorGetAll = true
			}
		} else {
			errorGetAll = true
		}
	} else {
		errorGetAll = true
	}

	if !errorGetAll {
		return resultado, nil
	}

	return nil, errors.New("error del servicio ActualizarInfoFamiliar:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
}

func ConsultarInfoEstudiante(idTercero string) (interface{}, error) {
	resultado := make(map[string]interface{})
	var persona []map[string]interface{}

	errPersona := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idTercero, &persona)
	if errPersona == nil && fmt.Sprintf("%v", persona[0]) != "map[]" {
		if persona[0]["Status"] != 404 {

			var correoPersonal []map[string]interface{}
			var programa []map[string]interface{}
			var telefono []map[string]interface{}
			var codigo []map[string]interface{}
			var programaNombre []map[string][]map[string]map[string]interface{}
			var correoInstitucional []map[string]interface{}
			var jsondata map[string]interface{}

			// resultado = persona[0]

			resultado["Nombre"] = persona[0]["NombreCompleto"]
			resultado["Id"] = persona[0]["Id"]

			errCorreo := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:53", &correoPersonal)
			if errCorreo == nil && fmt.Sprintf("%v", correoPersonal[0]) != "map[]" {
				if correoPersonal[0]["Status"] != 404 {
					correoaux := correoPersonal[0]["Dato"]
					if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
						panic(err)
					}

					if valor, existe := jsondata["Data"]; existe {
						resultado["CorreoPersonal"] = valor
					} else if valor, existe := jsondata["value"]; existe {
						resultado["CorreoPersonal"] = valor
					} else {
						resultado["CorreoPersonal"] = nil
					}
				}
			}

			errPrograma := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:95", &programa)
			if errPrograma == nil && fmt.Sprintf("%v", programa[0]) != "map[]" {
				if programa[0]["Status"] != 404 {
					programa := programa[0]["Dato"]
					if err := json.Unmarshal([]byte(programa.(string)), &jsondata); err != nil {
						panic(err)
					}

					errProgramaNombre := request.GetJson("http://"+beego.AppConfig.String("ProyectoAcademicoService")+"tr_proyecto_academico/"+fmt.Sprintf("%v", jsondata["value"]), &programaNombre)

					fmt.Println(errProgramaNombre)

					if fmt.Sprintf("%v", programaNombre[0]) != "map[]" {
						resultado["Carrera"] = programaNombre[0]["Enfasis"][0]["ProyectoAcademicoInstitucionId"]["Nombre"]
					}
				}
			}

			errTelefono := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:51", &telefono)
			if errTelefono == nil && fmt.Sprintf("%v", telefono[0]) != "map[]" {
				if telefono[0]["Status"] != 404 {
					telefonoaux := telefono[0]["Dato"]

					if err := json.Unmarshal([]byte(telefonoaux.(string)), &jsondata); err != nil {
						resultado["Telefono"] = telefono[0]["Dato"]
					} else {
						resultado["Telefono"] = jsondata["principal"]
					}
				}
			}

			errCodigoEst := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:93", &codigo)
			if errCodigoEst == nil && fmt.Sprintf("%v", codigo[0]) != "map[]" {
				if codigo[0]["Status"] != 404 {
					resultado["Codigo"] = codigo[0]["Dato"]
				}
			}

			errCorreoIns := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+fmt.Sprintf("%v", persona[0]["Id"])+",InfoComplementariaId__Id:94", &correoInstitucional)
			if errCorreoIns == nil && fmt.Sprintf("%v", correoInstitucional[0]) != "map[]" {
				if correoInstitucional[0]["Status"] != 404 {
					correoaux := correoInstitucional[0]["Dato"]
					if err := json.Unmarshal([]byte(correoaux.(string)), &jsondata); err != nil {
						panic(err)
					}

					resultado["CorreoInstitucional"] = jsondata["value"]
				}
			}
			return resultado, nil
		} else {
			return nil, errors.New("error del servicio ConsultarInfoEstudiante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}
	} else {
		logs.Error("Error --> ", errPersona)
		return nil, errors.New("error del servicio ConsultarInfoEstudiante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func GuardarAutor(data []byte) (interface{}, error) {
	//resultado solicitud de descuento
	var resultado map[string]interface{}
	//solicitud de descuento
	var tercero map[string]interface{}
	var terceroPost map[string]interface{}

	if err := json.Unmarshal(data, &tercero); err == nil {
		guardarpersona := map[string]interface{}{
			"NombreCompleto":      tercero["NombreCompleto"].(string),
			"Activo":              false,
			"TipoContribuyenteId": tercero["TipoContribuyenteId"],
		}
		errPersona := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero", "POST", &terceroPost, guardarpersona)

		if errPersona == nil && fmt.Sprintf("%v", terceroPost) != "map[]" && terceroPost["Id"] != nil {
			if terceroPost["Status"] != 400 {
				idTerceroCreado := terceroPost["Id"]
				var identificacion map[string]interface{}

				TerceroId := map[string]interface{}{
					"Id": idTerceroCreado,
				}
				identificaciontercero := map[string]interface{}{
					"Numero":          tercero["NumeroIdentificacion"],
					"TipoDocumentoId": tercero["TipoDocumentoId"],
					"TerceroId":       TerceroId,
					"Activo":          true,
				}
				errIdentificacion := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"datos_identificacion", "POST", &identificacion, identificaciontercero)
				if errIdentificacion == nil && fmt.Sprintf("%v", identificacion) != "map[]" && identificacion["Id"] != nil {
					if identificacion["Status"] != 400 {
						resultado = terceroPost

						resultado["NumeroIdentificacion"] = identificacion["Numero"]
						resultado["TipoIdentificacionId"] = identificacion["TipoDocumentoId"].(map[string]interface{})["Id"]
						return resultado, nil

					} else {
						//Si pasa un error borra todo lo creado al momento del registro del documento de identidad
						var resultado2 map[string]interface{}
						request.SendJson(fmt.Sprintf("http://"+beego.AppConfig.String("TercerosService")+"tercero/%.f", terceroPost["Id"]), "DELETE", &resultado2, nil)
						logs.Error("Error --> ", errIdentificacion)
						return nil, errors.New("error del servicio GuardarAutor")
					}
				} else {
					logs.Error("Error --> ", errIdentificacion)
					return nil, errors.New("error del servicio GuardarAutor")
				}
			} else {
				logs.Error("Error --> ", errPersona)
				return nil, errors.New("error del servicio GuardarAutor")
			}
		} else {
			logs.Error("Error --> ", errPersona)
			return nil, errors.New("error del servicio GuardarAutor")
		}
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio GuardarAutor")
	}
}

func ConsultarDatosAcudiente(idTercero string) (interface{}, error) {

	type InfoComplementariaTercero []struct {
		ID                   int    `json:"Id"`
		Dato                 string `json:"Dato"`
		InfoComplementariaId struct {
			CodigoAbreviacion string `json:"CodigoAbreviacion"`
		} `json:"InfoComplementariaId"`
	}

	// Obtener Id del grupo mediante el codigo de abreviacion

	var idGrupo []interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/grupo_info_complementaria?limit=0&query=CodigoAbreviacion:DATOS_CONT_ACU&query=Activo:true", &idGrupo)
	if err == nil && len(idGrupo) > 0 {
		if firstItem, ok := idGrupo[0].(map[string]interface{}); ok {
			if idGrup, ok := firstItem["Id"].(float64); ok {

				// Obtener la info complementaria mediante el id del grupo
				id, _ := strconv.Atoi(idTercero)
				id2, _ := strconv.Atoi(strconv.Itoa(int(idGrup)))
				var idInfoComplementaria []interface{}
				err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=0&&query=TerceroId.Id:"+strconv.Itoa(id)+",InfoComplementariaId.GrupoInfoComplementariaId.Id:"+strconv.Itoa(id2), &idInfoComplementaria)
				if err == nil && len(idInfoComplementaria) > 0 {
					idInfoComplementariaBytes, _ := json.Marshal(idInfoComplementaria)

					var informacionTablaRompiminetoResp InfoComplementariaTercero
					if err := json.Unmarshal(idInfoComplementariaBytes, &informacionTablaRompiminetoResp); err != nil {
						return nil, errors.New("error del servicio GuardarDatosAcudiente")
					} else {
						return informacionTablaRompiminetoResp, nil
					}

				} else {
					return nil, errors.New("error del servicio GuardarDatosAcudiente")
				}

			} else {
				return nil, errors.New("error del servicio GuardarDatosAcudiente")
			}
		} else {
			return nil, errors.New("error del servicio GuardarDatosAcudiente")
		}
	} else {
		return nil, errors.New("error del servicio GuardarDatosAcudiente")
	}
}

func GuardarDatosAcudiente(idTercero string, data []byte) (interface{}, error) {

	type Data struct {
		Nombre            string `json:"nombre"`
		Parentezco        string `json:"parentezco"`
		Correo            string `json:"correo"`
		Direccion         string `json:"direccion"`
		Telefono          string `json:"telefono"`
		NombreSegundo     string `json:"nombreSegundo"`
		ParentezcoSegundo string `json:"parentezcoSegundo"`
		CorreoSegundo     string `json:"correoSegundo"`
		DireccionSegundo  string `json:"direccionSegundo"`
		TelefonoSegundo   string `json:"telefonoSegundo"`
	}

	type InfoComplementaria []struct {
		ID                int    `json:"Id"`
		CodigoAbreviacion string `json:"CodigoAbreviacion"`
	}

	// Obtener Id del grupo mediante el codigo de abreviacion
	var idGrupo []interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/grupo_info_complementaria?limit=0&query=CodigoAbreviacion:DATOS_CONT_ACU&query=Activo:true", &idGrupo)
	if err == nil && len(idGrupo) > 0 {
		if firstItem, ok := idGrupo[0].(map[string]interface{}); ok {
			if idGrup, ok := firstItem["Id"].(float64); ok {

				// Obtener Id de la info complementaria mediante el id del grupo
				var idInfoComplementaria []interface{}
				err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria?limit=0&query=GrupoInfoComplementariaId.Id:"+strconv.Itoa(int(idGrup))+"&query=Activo:true", &idInfoComplementaria)
				if err == nil && len(idInfoComplementaria) > 0 {
					idInfoComplementariaBytes, _ := json.Marshal(idInfoComplementaria)
					var idInfoComplementariaResp InfoComplementaria
					if err := json.Unmarshal(idInfoComplementariaBytes, &idInfoComplementariaResp); err != nil {
						return nil, errors.New("error del servicio GuardarDatosAcudiente")
					} else {

						var datos Data
						if err := json.Unmarshal(data, &datos); err != nil {
							log.Fatal(err)
						}
						// Almacena los datos en una slice
						datosSlice := []string{
							datos.Nombre,
							datos.Parentezco,
							datos.Correo,
							datos.Direccion,
							datos.Telefono,
							datos.NombreSegundo,
							datos.ParentezcoSegundo,
							datos.CorreoSegundo,
							datos.DireccionSegundo,
							datos.TelefonoSegundo,
						}

						// Asegúrate de que hay suficientes elementos en idInfoComplementariaResp
						if len(idInfoComplementariaResp) < len(datosSlice) {
							log.Fatal("No hay suficientes elementos en idInfoComplementariaResp")
						}

						for i, element := range idInfoComplementariaResp {

							id, _ := strconv.Atoi(idTercero)

							newInfo := map[string]interface{}{
								"TerceroId":            map[string]interface{}{"Id": id},
								"InfoComplementariaId": map[string]interface{}{"Id": element.ID},
								"Dato":                 fmt.Sprintf("{\"Dato\": \"%s\"}", datosSlice[i]),
								"Activo":               true,
							}

							// Guardar en la base de datos
							var apiResp map[string]interface{}
							err := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", newInfo, &apiResp)
							if err != nil {
								logs.Error("Error al enviar la solicitud a la API CRUD", err)
								return nil, errors.New("error del servicio GuardarDatosAcudiente")
							}
						}
						return nil, nil
					}
				} else {
					return nil, errors.New("error del servicio GuardarDatosAcudiente")
				}
			} else {
				return nil, errors.New("error del servicio GuardarDatosAcudiente")
			}
		} else {
			return nil, errors.New("error del servicio GuardarDatosAcudiente")
		}
	} else {
		return nil, errors.New("error del servicio GuardarDatosAcudiente")
	}
}

func ActualizarDatosAcudiente(idTercero string, data []byte) (interface{}, error) {

	type Data struct {
		Nombre            string `json:"nombre"`
		Parentezco        string `json:"parentezco"`
		Correo            string `json:"correo"`
		Direccion         string `json:"direccion"`
		Telefono          string `json:"telefono"`
		NombreSegundo     string `json:"nombreSegundo"`
		ParentezcoSegundo string `json:"parentezcoSegundo"`
		CorreoSegundo     string `json:"correoSegundo"`
		DireccionSegundo  string `json:"direccionSegundo"`
		TelefonoSegundo   string `json:"telefonoSegundo"`
	}

	var datos Data
	if err := json.Unmarshal(data, &datos); err != nil {
		log.Fatal(err)
	}

	type GrupoInfo []struct {
		ID int `json:"Id"`
	}

	type InfoComplementaria []struct {
		ID                int    `json:"Id"`
		CodigoAbreviacion string `json:"CodigoAbreviacion"`
	}

	type InfoComplementariaTercero []struct {
		ID                   int    `json:"Id"`
		Dato                 string `json:"Dato"`
		InfoComplementariaId struct {
			CodigoAbreviacion string `json:"CodigoAbreviacion"`
		} `json:"InfoComplementariaId"`
	}

	// Obtener Id del grupo mediante el codigo de abreviacion
	var idGrupo []interface{}
	err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/grupo_info_complementaria?limit=0&query=CodigoAbreviacion:DATOS_CONT_ACU&query=Activo:true", &idGrupo)
	if err == nil && len(idGrupo) > 0 {
		if firstItem, ok := idGrupo[0].(map[string]interface{}); ok {
			if idGrup, ok := firstItem["Id"].(float64); ok {

				// Obtener la info complementaria mediante el id del grupo
				id2, _ := strconv.Atoi(strconv.Itoa(int(idGrup)))
				var idInfoComplementaria []interface{}
				err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria?limit=0&query=GrupoInfoComplementariaId.Id:"+strconv.Itoa(id2)+"&query=Activo:true", &idInfoComplementaria)
				if err == nil && len(idInfoComplementaria) > 0 {
					idInfoComplementariaBytes, _ := json.Marshal(idInfoComplementaria)
					var idInfoComplementariaResp InfoComplementaria
					if err := json.Unmarshal(idInfoComplementariaBytes, &idInfoComplementariaResp); err != nil {
						return nil, errors.New("error del servicio ActualizarDatosAcudiente")
					} else {
						// Obtener Id de la info complementaria mediante el id del grupo
						id, _ := strconv.Atoi(idTercero)
						var idInfoComplementariaTercero []interface{}
						err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=0&query=TerceroId.Id:"+strconv.Itoa(id)+",InfoComplementariaId.GrupoInfoComplementariaId.Id:"+strconv.Itoa(id2), &idInfoComplementariaTercero)
						if err == nil && len(idInfoComplementariaTercero) > 0 {
							idInfoComplementariaTerceroBytes, _ := json.Marshal(idInfoComplementariaTercero)
							var idInfoComplementariaTerceroResp InfoComplementariaTercero
							if err := json.Unmarshal(idInfoComplementariaTerceroBytes, &idInfoComplementariaTerceroResp); err != nil {
								return nil, errors.New("error del servicio ActualizarDatosAcudiente")
							} else {
								// Almacena los datos en una slice
								datosSlice := []string{
									datos.Nombre,
									datos.Parentezco,
									datos.Correo,
									datos.Direccion,
									datos.Telefono,
									datos.NombreSegundo,
									datos.ParentezcoSegundo,
									datos.CorreoSegundo,
									datos.DireccionSegundo,
									datos.TelefonoSegundo,
								}

								// Asegúrate de que hay suficientes elementos en idInfoComplementariaResp
								if len(idInfoComplementariaResp) < len(datosSlice) {
									log.Fatal("No hay suficientes elementos en idInfoComplementariaResp")
								}

								id, _ := strconv.Atoi(idTercero)
								var informacionTablaRompimineto []interface{}
								err := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"/info_complementaria_tercero?limit=0&query=TerceroId.Id:"+strconv.Itoa(id)+",InfoComplementariaId.GrupoInfoComplementariaId.Id:"+strconv.Itoa(id2), &informacionTablaRompimineto)
								if err == nil && len(informacionTablaRompimineto) > 0 {
									informacionTablaRompiminetoBytes, _ := json.Marshal(informacionTablaRompimineto)
									var informacionTablaRompiminetoResp InfoComplementariaTercero
									if err := json.Unmarshal(informacionTablaRompiminetoBytes, &informacionTablaRompiminetoResp); err != nil {
										return nil, errors.New("error del servicio ActualizarDatosAcudiente")
									} else {

										for i, element := range idInfoComplementariaResp {
											idPut := 0
											switch informacionTablaRompiminetoResp[i].InfoComplementariaId.CodigoAbreviacion {
											case "NOM_PRI_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "PAREN_PRI_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "CORREO_PRI_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "DIREC_PRI_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "TEL_PRI_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "NOM_SEG_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "PAREN_SEG_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "CORREO_SEG_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "DIREC_SEG_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											case "TEL_SEG_ACU":
												idPut = informacionTablaRompiminetoResp[i].ID
											}

											id, _ := strconv.Atoi(idTercero)
											// Formatea los datos en un mapa para la solicitud de actualización
											updateInfo := map[string]interface{}{
												"TerceroId":            map[string]interface{}{"Id": id},
												"InfoComplementariaId": map[string]interface{}{"Id": element.ID},
												"Dato":                 fmt.Sprintf("{\"Dato\": \"%s\"}", datosSlice[i]),
												"Activo":               true,
											}

											// Realiza la solicitud de actualización a la API
											var apiResp map[string]interface{}
											err := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idPut), "PUT", updateInfo, &apiResp)
											if err != nil {
												logs.Error("Error al enviar la solicitud a la API CRUD", err)
												return nil, errors.New("error del servicio ActualizarDatosAcudiente")
											}
										}
									}
								} else {
									return nil, errors.New("error del servicio ActualizarDatosAcudiente")
								}

							}
						} else {
							return nil, errors.New("error del servicio ActualizarDatosAcudiente")
						}
					}
				} else {
					return nil, errors.New("error del servicio ActualizarDatosAcudiente")
				}
			} else {
				return nil, errors.New("error del servicio ActualizarDatosAcudiente")
			}
		} else {
			return nil, errors.New("error del servicio ActualizarDatosAcudiente")
		}
	} else {
		return nil, errors.New("error del servicio ActualizarDatosAcudiente")
	}
	return nil, nil
}

func ConsultarLocalidades() (interface{}, error) {

	type Localidad struct {
		Id     int    `json:"Id"`
		Nombre string `json:"Nombre"`
	}

	type TipoColegio struct {
		Id     int    `json:"Id"`
		Nombre string `json:"Nombre"`
	}

	type ValidoBachillerato struct {
		Id     int    `json:"Id"`
		Nombre string `json:"Nombre"`
	}

	type SemestresTranscurridos struct {
		Id     int    `json:"Id"`
		Nombre string `json:"Nombre"`
	}

	type Medio struct {
		Id     int    `json:"Id"`
		Nombre string `json:"Nombre"`
	}

	type Presentaciones struct {
		Id     int    `json:"Id"`
		Nombre string `json:"Nombre"`
	}

	type Respuestas struct {
		Localidades  []Localidad              `json:"localidades"`
		TipoColegio  []TipoColegio            `json:"tipoColegio"`
		Valido       []ValidoBachillerato     `json:"valido"`
		Semestres    []SemestresTranscurridos `json:"semestres"`
		Medio        []Medio                  `json:"medio"`
		Presentacion []Presentaciones         `json:"presentacion"`
	}

	var localidades []Localidad
	errLoclidades := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?limit=0&query=GrupoInfoComplementariaId.CodigoAbreviacion:LOCALIDADES&query=Activo:true", &localidades)
	if errLoclidades == nil {
		var tipoColegio []TipoColegio
		errColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?limit=0&query=GrupoInfoComplementariaId.CodigoAbreviacion:TIPO_COLEGIO&query=Activo:true", &tipoColegio)
		if errColegio == nil {
			var validoBachillerato []ValidoBachillerato
			errValido := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?limit=0&query=GrupoInfoComplementariaId.CodigoAbreviacion:VALIDO&query=Activo:true", &validoBachillerato)
			if errValido == nil {
				var semestresTranscurridos []SemestresTranscurridos
				errSemestres := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?limit=0&query=GrupoInfoComplementariaId.CodigoAbreviacion:S_TRANSCURRIDOS&query=Activo:true", &semestresTranscurridos)
				if errSemestres == nil {
					var medio []Medio
					errMedio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?limit=0&query=GrupoInfoComplementariaId.CodigoAbreviacion:MEDIO&query=Activo:true", &medio)
					if errMedio == nil {
						var presentaciones []Presentaciones
						errPresentaciones := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria?limit=0&query=GrupoInfoComplementariaId.CodigoAbreviacion:PRESENTACIONES_U&query=Activo:true", &presentaciones)
						if errPresentaciones == nil {
							respuesta := Respuestas{
								Localidades:  localidades,
								TipoColegio:  tipoColegio,
								Valido:       validoBachillerato,
								Semestres:    semestresTranscurridos,
								Medio:        medio,
								Presentacion: presentaciones,
							}
							fmt.Println("GRUPO INFO", respuesta)
							return respuesta, nil
						} else {
							logs.Error("Error --> ", errPresentaciones)
							return nil, errors.New("error del servicio al consultar los datos de Presentaciones")
						}
					} else {
						logs.Error("Error --> ", errMedio)
						return nil, errors.New("error del servicio al consultar los datos de Medio")
					}
				} else {
					logs.Error("Error --> ", errSemestres)
					return nil, errors.New("error del servicio al consultar los datos de Semestres Transcurridos")
				}
			} else {
				logs.Error("Error --> ", errValido)
				return nil, errors.New("error del servicio al consultar los datos de Valido Bachillerato")
			}
		} else {
			logs.Error("Error --> ", errColegio)
			return nil, errors.New("error del servicio al consultar los colegios")
		}
	} else {
		logs.Error("Error --> ", errLoclidades)
		return nil, errors.New("error del servicio al consultar las Localidades")
	}
}

func ConsultarInfoAcademicaAspirante(idTercero string) (interface{}, error) {
	type Aspirante struct {
		Id int `json:"Id"`
	}

	type Localidad struct {
		Id                   int `json:"Id"`
		InfoComplementariaId struct {
			Id     int    `json:"Id"`
			Nombre string `json:"Nombre"`
		} `json:"InfoComplementariaId"`
	}

	type Colegio struct {
		Id                   int `json:"Id"`
		InfoComplementariaId struct {
			Id     int    `json:"Id"`
			Nombre string `json:"Nombre"`
		} `json:"InfoComplementariaId"`
	}

	type Valido struct {
		Id                   int `json:"Id"`
		InfoComplementariaId struct {
			Id     int    `json:"Id"`
			Nombre string `json:"Nombre"`
		} `json:"InfoComplementariaId"`
	}

	type Semestres struct {
		Id                   int `json:"Id"`
		InfoComplementariaId struct {
			Id     int    `json:"Id"`
			Nombre string `json:"Nombre"`
		} `json:"InfoComplementariaId"`
	}

	type Medio struct {
		Id                   int `json:"Id"`
		InfoComplementariaId struct {
			Id     int    `json:"Id"`
			Nombre string `json:"Nombre"`
		} `json:"InfoComplementariaId"`
	}

	type Presentaciones struct {
		Id                   int `json:"Id"`
		InfoComplementariaId struct {
			Id     int    `json:"Id"`
			Nombre string `json:"Nombre"`
		} `json:"InfoComplementariaId"`
	}

	type Respuestas struct {
		Localidades  []Localidad      `json:"localidades"`
		TipoColegio  []Colegio        `json:"colegio"`
		Valido       []Valido         `json:"valido"`
		Semestres    []Semestres      `json:"semestres"`
		Medio        []Medio          `json:"medio"`
		Presentacion []Presentaciones `json:"presentacion"`
	}

	var aspirante []Aspirante
	errAspirante := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero?query=Id:"+idTercero, &aspirante)
	if errAspirante == nil {

		var localidad []Localidad
		errLocalidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:LOCALIDADES", &localidad)
		if errLocalidad == nil {
			fmt.Println("ASPIRANTE RESPUESTA", localidad)
		} else {
			logs.Error("Error --> ", errLocalidad)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

		var colegio []Colegio
		errColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:TIPO_COLEGIO", &colegio)
		if errColegio == nil {
			fmt.Println("ASPIRANTE RESPUESTA", colegio)
		} else {
			logs.Error("Error --> ", errColegio)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

		var valido []Valido
		errValido := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:VALIDO", &valido)
		if errValido == nil {
			fmt.Println("ASPIRANTE RESPUESTA", valido)
		} else {
			logs.Error("Error --> ", errValido)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

		var semestres []Semestres
		errSemestres := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:S_TRANSCURRIDOS", &semestres)
		if errSemestres == nil {
			fmt.Println("ASPIRANTE RESPUESTA", semestres)
		} else {
			logs.Error("Error --> ", errSemestres)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

		var medio []Medio
		errMedio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:MEDIO", &medio)
		if errMedio == nil {
			fmt.Println("ASPIRANTE RESPUESTA", medio)
		} else {
			logs.Error("Error --> ", errMedio)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

		var presentaciones []Presentaciones
		errPresentaciones := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:PRESENTACIONES_U", &presentaciones)
		if errPresentaciones == nil {
			fmt.Println("ASPIRANTE RESPUESTA", presentaciones)
		} else {
			logs.Error("Error --> ", errPresentaciones)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
		}

		respuesta := Respuestas{
			Localidades:  localidad,
			TipoColegio:  colegio,
			Valido:       valido,
			Semestres:    semestres,
			Medio:        medio,
			Presentacion: presentaciones,
		}

		return respuesta, nil
	} else {
		logs.Error("Error --> ", errAspirante)
		return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante:   La solicitud contiene un tipo de dato incorrecto o un parámetro inválido")
	}
}

func CrearLocalidades(idTercero string, data []byte) (interface{}, error) {

	type Data struct {
		Localidad   int `json:"localidad"`
		TipoColegio int `json:"tipoColegio"`
		Valido      int `json:"valido"`
		Semestre    int `json:"semestre"`
		Medio       int `json:"medio"`
		Oportunidad int `json:"oportunidad"`
	}

	var datos Data
	if err := json.Unmarshal(data, &datos); err == nil {
		id, _ := strconv.Atoi(idTercero)
		datosSlice := []int{
			datos.Localidad,
			datos.TipoColegio,
			datos.Valido,
			datos.Semestre,
			datos.Medio,
			datos.Oportunidad,
		}
		for _, element := range datosSlice {
			updateInfo := map[string]interface{}{
				"TerceroId":            map[string]interface{}{"Id": id},
				"InfoComplementariaId": map[string]interface{}{"Id": element},
				"Activo":               true,
			}
			var apiResp map[string]interface{}
			err := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", updateInfo, &apiResp)
			if err != nil {
				logs.Error("Error al enviar la solicitud a la API CRUD", err)
				return nil, errors.New("error del servicio al crear las localidades")
			}
		}
		return nil, nil
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio al crear las localidades")
	}
}

func ActualizarInfoAcademicaAspirante(idTercero string, data []byte) (interface{}, error) {
	type Data struct {
		Localidad   int `json:"localidad"`
		TipoColegio int `json:"tipoColegio"`
		Valido      int `json:"valido"`
		Semestre    int `json:"semestre"`
		Medio       int `json:"medio"`
		Oportunidad int `json:"oportunidad"`
	}

	type Localidad struct {
		Id int `json:"Id"`
	}

	type Colegio struct {
		Id int `json:"Id"`
	}

	type Valido struct {
		Id int `json:"Id"`
	}

	type Semestres struct {
		Id int `json:"Id"`
	}

	type Medio struct {
		Id int `json:"Id"`
	}

	type Presentaciones struct {
		Id int `json:"Id"`
	}

	var datos Data
	if err := json.Unmarshal(data, &datos); err == nil {

		var localidad []Localidad
		errLocalidad := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:LOCALIDADES", &localidad)
		if errLocalidad == nil {
		} else {
			logs.Error("Error --> ", errLocalidad)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante")
		}

		var colegio []Colegio
		errColegio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:TIPO_COLEGIO", &colegio)
		if errColegio == nil {
		} else {
			logs.Error("Error --> ", errColegio)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante")
		}

		var valido []Valido
		errValido := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:VALIDO", &valido)
		if errValido == nil {
		} else {
			logs.Error("Error --> ", errValido)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante")
		}

		var semestres []Semestres
		errSemestres := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:S_TRANSCURRIDOS", &semestres)
		if errSemestres == nil {
		} else {
			logs.Error("Error --> ", errSemestres)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante")
		}

		var medio []Medio
		errMedio := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:MEDIO", &medio)
		if errMedio == nil {
		} else {
			logs.Error("Error --> ", errMedio)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante")
		}

		var presentaciones []Presentaciones
		errPresentaciones := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero?query=TerceroId.Id:"+idTercero+",InfoComplementariaId.GrupoInfoComplementariaId.CodigoAbreviacion:PRESENTACIONES_U", &presentaciones)
		if errPresentaciones == nil {
		} else {
			logs.Error("Error --> ", errPresentaciones)
			return nil, errors.New("error del servicio ConsultarInfoAcademicaAspirante")
		}

		//Datos de la peticion del cliente
		datosActualizarSlice := []int{
			datos.Localidad,
			datos.TipoColegio,
			datos.Valido,
			datos.Semestre,
			datos.Medio,
			datos.Oportunidad,
		}

		// Ids de la info complementaria que vamos a actualizar
		idSlice := []int{
			localidad[0].Id,
			colegio[0].Id,
			valido[0].Id,
			semestres[0].Id,
			medio[0].Id,
			presentaciones[0].Id,
		}

		for i := range datosActualizarSlice {
			id, _ := strconv.Atoi(idTercero)
			// Formatea los datos en un mapa para la solicitud de actualización
			updateInfo := map[string]interface{}{
				"TerceroId":            map[string]interface{}{"Id": id},
				"InfoComplementariaId": map[string]interface{}{"Id": datosActualizarSlice[i]},
				"Activo":               true,
			}

			// Realiza la solicitud de actualización a la API
			var apiResp map[string]interface{}
			err := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero/"+fmt.Sprintf("%.f", idSlice[i]), "PUT", updateInfo, &apiResp)
			if err != nil {
				logs.Error("Error al enviar la solicitud a la API CRUD", err)
				return nil, errors.New("error del servicio al actualizar las localidades")
			}
		}
		return nil, nil
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio al actualizar las localidades")
	}
}

func AsignarCorreoInstitucional(data []byte) (interface{}, error) {
	var body []map[string]interface{}
	var response []map[string]interface{}
	if err := json.Unmarshal(data, &body); err == nil {
		IdCorreo, _ := models.IdInfoCompTercero("10", "CI")
		ItCorreo, _ := strconv.ParseFloat(IdCorreo, 64)
		for _, info := range body {
			var updateTercero map[string]interface{}
			errUpdtTercero := request.GetJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", info["Id"]), &updateTercero)
			if errUpdtTercero == nil && updateTercero["Status"] != 404 {
				// Verificar si ya tiene un correo institucional
				if updateTercero["UsuarioWSO2"] == nil || updateTercero["UsuarioWSO2"] == "" {
					updateTercero["UsuarioWSO2"] = info["Correo"]

					var updateAnswer map[string]interface{}
					errupdateAnswer := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"tercero/"+fmt.Sprintf("%v", info["Id"]), "PUT", &updateAnswer, updateTercero)
					if errupdateAnswer == nil {
						response = append(response, updateAnswer)
						newInfo := map[string]interface{}{
							"TerceroId":            map[string]interface{}{"Id": info["Id"]},
							"InfoComplementariaId": map[string]interface{}{"Id": ItCorreo},
							"Dato":                 fmt.Sprintf("{\"value\": \"%v\"}", info["Correo"]),
							"Activo":               true,
						}

						var createinfo map[string]interface{}
						errCreateInfo := request.SendJson("http://"+beego.AppConfig.String("TercerosService")+"info_complementaria_tercero", "POST", &createinfo, newInfo)
						if errCreateInfo == nil && fmt.Sprintf("%v", createinfo) != "map[]" && createinfo["Id"] != nil {
							response = append(response, createinfo)
						} else {
							log.Error(errCreateInfo)
							response = append(response, map[string]interface{}{
								"Id":    info["Id"],
								"Error": "Error al crear registro",
							})
						}

						// Asignar rol "ESTUDIANTE"
						roleInfo := map[string]interface{}{
							"TerceroId": info["Id"],
							"role":      "ESTUDIANTE",
						}
						var roleResponse map[string]interface{}
						errRole := request.SendJson("http://"+beego.AppConfig.String("RolesService")+"roles", "POST", &roleResponse, roleInfo)
						if errRole == nil {
							response = append(response, roleResponse)
						} else {
							log.Error(errRole)
							response = append(response, map[string]interface{}{
								"Id":    info["Id"],
								"Error": "Error al asignar rol",
							})
						}
					} else {
						log.Error(errupdateAnswer)
						response = append(response, map[string]interface{}{
							"Id":    info["Id"],
							"Error": "Error al actualizar la información",
						})
					}
				} else {
					response = append(response, map[string]interface{}{
						"Id":    info["Id"],
						"Error": "Correo institucional ya existe",
					})
				}
			} else {
				log.Error(errUpdtTercero)
				response = append(response, map[string]interface{}{
					"Id":    info["Id"],
					"Error": "Hubo problemas al consultar el servicio",
				})
			}
		}
		return response, nil
	} else {
		logs.Error("Error --> ", err)
		return nil, errors.New("error del servicio AsignarCorreoInstitucional: La solicitud contiene un tipo de dato incorrecto o un parámetro inválido" + err.Error())
	}
}
