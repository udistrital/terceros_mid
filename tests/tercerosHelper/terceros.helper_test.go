package tercerosHelper_test

import (
	"flag"
	"os"
	"testing"

	"github.com/astaxie/beego"
	"github.com/udistrital/terceros_mid/helpers/tipos"
)

var parameters struct {
	PARAMETROS_CRUD  string
	TERCEROS_SERVICE string
	OIKOS2_CRUD      string
}

func TestMain(m *testing.M) {
	parameters.PARAMETROS_CRUD = os.Getenv("PARAMETROS_CRUD")
	beego.AppConfig.Set("ParametroService", parameters.PARAMETROS_CRUD)
	parameters.TERCEROS_SERVICE = os.Getenv("TERCEROS_SERVICE")
	beego.AppConfig.Set("TercerosService", parameters.TERCEROS_SERVICE)
	parameters.OIKOS2_CRUD = os.Getenv("OIKOS2_CRUD")
	beego.AppConfig.Set("OikosService", parameters.OIKOS2_CRUD)
	flag.Parse()
	os.Exit(m.Run())
}

// TestGetContratista ...
func TestGetContratista(t *testing.T) {

	testStr := "perez"
	if valor, err := tipos.GetContratista(9825, testStr); err != nil {
		t.Error("No se pudo consultar el contratista", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetContratista Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetFuncionariosPlanta(t *testing.T) {

	testStr := "perez"
	if valor, err := tipos.GetFuncionariosPlanta(9801, testStr); err != nil {
		t.Error("No se pudo consultar el funcionario", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetFuncionariosPlanta Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetOrdenadores(t *testing.T) {

	testStr := "perez"
	if valor, err := tipos.GetOrdenadores(9804, testStr); err != nil {
		t.Error("No se pudo consultar el ordenador", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetOrdenadores Finalizado Correctamente")
	}
}

// TestGetFuncionariosPlanta ...
func TestGetProveedor(t *testing.T) {

	testStr := "perez"
	if valor, err := tipos.GetProveedor(9769, testStr); err != nil {
		t.Error("No se pudo consultar el proveedor", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetProveedor Finalizado Correctamente")
	}
}

// TestGetTipos ...
func TestGetTipos(t *testing.T) {

	if valor, err := tipos.GetTipos(); err != nil {
		t.Error("No se pudo consultar los tipos", err)
		t.Fail()
	} else {
		t.Log(valor)
		t.Log("TestGetTipos Finalizado Correctamente")
	}
}

func TestEndPointParametrosService(t *testing.T) {
	t.Log("Testing EndPoint ParametroService")
	t.Log(parameters.PARAMETROS_CRUD)
	t.Log("Testing EndPoint TercerosService")
	t.Log(parameters.TERCEROS_SERVICE)
	t.Log("Testing EndPoint OikosService")
	t.Log(parameters.OIKOS2_CRUD)
}
