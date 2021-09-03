# terceros_mid

MID de Terceros

- [x] Tipos de Terceros
- [ ] Propiedades de Terceros

Api intermediaria entre el cliente de plan de necesidades y las apis necesarios para la gestión de la información para estos mismos con respecto a las dependencias.
Api mid para el subsistema de terceros que hace parte del sistema kronos

## Especificaciones Técnicas

### Tecnologías Implementadas y Versiones

- [Golang](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/golang.md)
- [BeeGo](https://github.com/udistrital/introduccion_oas/blob/master/instalacion_de_herramientas/beego.md)

### Variables de Entorno

- Listadas/Documentadas en la plantilla [template.env](template.env)
- Para usarlas:

  1. copiar la plantilla, por ejemplo como ".env": `cp template.env .env`
  2. Ajustar los valores en la copia creada, por ejemplo, con el editor nano: `nano .env`
  3. Traer sus valores al entorno con `source .env`

### Ejecución del Proyecto

```shell
#1. Obtener el repositorio con Go
go get github.com/udistrital/terceros_mid

#2. Moverse a la carpeta del repositorio
cd $GOPATH/src/github.com/udistrital/terceros_mid

# 3. Moverse a la rama **develop**
git pull origin develop && git checkout develop

# 4. alimentar todas las variables de entorno que utiliza el proyecto.
cp template.env .env
source .env
PORT=8080
```

### Ejecución Dockerfile

```shell
# Generar imagen de despliegue (release/master)
docker build .
```

### APIs Requeridas

1. [terceros_crud](https://github.com/udistrital/terceros_crud)
2. [oikos2_crud](https://github.com/udistrital/oikos_api)
3. [parametrod_crud](https://github.com/udistrital/parametros_crud)

### Ejecución Pruebas

Pruebas unitarias

```shell
go test
```

## Estado CI

| Develop | Release 0.1.0 | Master |
| -- | -- | -- |
| [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/terceros_mid/status.svg?ref=refs/heads/develop)](https://hubci.portaloas.udistrital.edu.co/udistrital/terceros_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/terceros_mid/status.svg?ref=refs/heads/release/0.1.0)](https://hubci.portaloas.udistrital.edu.co/udistrital/terceros_mid) | [![Build Status](https://hubci.portaloas.udistrital.edu.co/api/badges/udistrital/terceros_mid/status.svg?ref=refs/heads/master)](https://hubci.portaloas.udistrital.edu.co/udistrital/terceros_mid) |

## Licencia

This file is part of terceros_mid

terceros_mid is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

terceros_mid is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with terceros_mid. If not, see https://www.gnu.org/licenses/.
