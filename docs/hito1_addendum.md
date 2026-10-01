# Hito 1 · Addendum técnico

**Pareja:** Dhylan Parrales · Luis Anzules
**Paralelo:** Aplicaciones Web II B

## A. Estructura del proyecto
.
├── main.go               # Punto de entrada de la aplicación y servidor
├── go.mod                # Dependencias del proyecto en Go
├── internal/

│   └── bitfamily/        # Lógica central del negocio BitFamily
│       ├── modelos.go    # Definición de structs
│       ├── manejadores.go # Funciones controladoras (Handlers)
│       └── rutas.go      # Definición de rutas
└── docs/                 # Documentación del hito 1
## B. Configuración y secretos

| Variable | Para qué sirve | Ejemplo (sin datos reales) |
|----------|----------------|----------------------------|
| PORT | Puerto donde corre el servidor Go | 8080 |
| DB_HOST | Host de la base de datos PostgreSQL | localhost |
| DB_PASS | Contraseña de acceso a la base de datos | miclavesecreta |

**Archivo de ejemplo:** `.env.example`

## C. Pruebas

| Prueba | Qué caso cubre |
|--------|----------------|
| TestCrearTarea | Valida que un padre pueda insertar una tarea con saldo positivo. |
| TestTransicionProhibida | Valida que retorne error 409 al intentar pasar de `aprobada` a `disponible`. |

**Captura de `go test ./...`:** ![Pruebas](hito1_pruebas.png) 

## D. Boceto de la pantalla principal

![Boceto](hito1_boceto.png)

## E. Diagrama de secuencia del caso de uso principal

```mermaid
sequenceDiagram
    Hijo->>Pantalla Tareas: toca "Completar Tarea"
    Pantalla Tareas->>API: PATCH /api/tareas/1042/estado {"estado":"completada"}
    API->>BD: valida estado actual y roles
    BD-->>API: actualización ok
    API-->>Pantalla Tareas: 200 OK con tarea actualizada
    Padre->>Pantalla Tareas: revisa y toca "Aprobar Tarea"
    Pantalla Tareas->>API: PATCH /api/tareas/1042/estado {"estado":"aprobada"}
    API->>BD: valida la transición prohibida y suma dinero a saldo_usd del hijo
    BD-->>API: actualización ok
    API-->>Pantalla Tareas: 200 OK