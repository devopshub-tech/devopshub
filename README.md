# # DevOpsHub

## Organización de las carpetas:

```
devopshub/
|-- cmd/
|   |-- master/
|   |   |-- main.go           # Punto de entrada para el servicio master
|   |
|   |-- worker/
|       |-- main.go           # Punto de entrada para el servicio worker
|
|-- pkg/
|   |-- common/               # Paquetes compartidos entre los servicios
|   |   |-- logging/
|   |   |-- utils/
|   |
|   |-- master/               # Paquetes específicos para el servicio master
|   |   |-- ...
|   |
|   |-- worker/               # Paquetes específicos para el servicio worker
|       |-- ...
|
|-- internal/                 # Código que no será importado por otros repositorios
|   |-- config/               # Configuraciones específicas
|   |-- db/                   # Acceso a la base de datos
|
|-- api/                      # Definiciones de API compartidas o gRPC si aplica
|
|-- deployments/              # Archivos de despliegue (Docker, Kubernetes, etc.)
|
|-- scripts/                  # Scripts útiles para el desarrollo y despliegue
|
|-- README.md                 # Documentación principal del monorepo
|-- go.mod                    # Archivo de definición de módulos Go
```
