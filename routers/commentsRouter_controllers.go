package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:PropiedadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:PropiedadesController"],
		beego.ControllerComments{
			Method:           "GetPropiedades",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:PropiedadesController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:PropiedadesController"],
		beego.ControllerComments{
			Method:           "GetPropiedadesDeUnTerceroId",
			Router:           "/:propiedad/:idTercero",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "GuardarPersona",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ActualizarPersona",
			Router:           "/",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarPersona",
			Router:           "/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosComplementarios",
			Router:           "/:tercero_id/complementarios",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosContacto",
			Router:           "/:tercero_id/contacto",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosFamiliar",
			Router:           "/:tercero_id/familiar",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosFormacionPregrado",
			Router:           "/:tercero_id/formacion-pregrado",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarInfoEstudiante",
			Router:           "/:tercero_id/info-solicitante",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "GuardarAutor",
			Router:           "/autores",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ActualizarDatosComplementarios",
			Router:           "/complementarios",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "GuardarDatosComplementarios",
			Router:           "/complementarios",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "GuardarDatosComplementariosParAcademico",
			Router:           "/complementarios-par-academico",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "GuardarDatosContacto",
			Router:           "/contacto",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarExistenciaPersona",
			Router:           "/existencia/:numeroDocumento",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ActualizarInfoFamiliar",
			Router:           "/info-familiar",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ObtenerTercerosConNIT",
			Router:           "/nit",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarDatosAcudiente",
			Router:           "/datos-acudiente/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "GuardarDatosAcudiente",
			Router:           "/datos-acudiente/:tercero_id",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ActualizarDatosAcudiente",
			Router:           "/datos-acudiente/:tercero_id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarLocalidades",
			Router:           "/localidades",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ConsultarInfoAcademicaAspirante",
			Router:           "/localidades/:tercero_id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "CrearLocalidades",
			Router:           "/localidades/:tercero_id",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_tercero_mid/controllers:SgaTercerosController"],
		beego.ControllerComments{
			Method:           "ActualizarInfoAcademicaAspirante",
			Router:           "/localidades/:tercero_id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"],
		beego.ControllerComments{
			Method:           "GetTipos",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"],
		beego.ControllerComments{
			Method:           "GetByTipo",
			Router:           "/:tipo",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"] = append(beego.GlobalControllerRouter["github.com/udistrital/terceros_mid/controllers:TiposController"],
		beego.ControllerComments{
			Method:           "GetByTipoAndID",
			Router:           "/:tipo/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
