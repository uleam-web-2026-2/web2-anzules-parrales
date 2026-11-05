# Decisiones de Arquitectura – BitFamily API

## Semana 8: Estructura en Capas

Para cumplir con los principios de responsabilidad única y separación de conceptos en nuestra API, hemos organizado el código de la entidad principal (`tareas`) en las siguientes capas dentro del directorio `internal/`:

1. **Capa de Transporte (`internal/transporte/`)**:
   - **Qué contiene**: Los manejadores HTTP (handlers), las rutas configuradas con Chi y la decodificación de payloads JSON.
   - **Responsabilidad**: Atiende las peticiones HTTP y devuelve las respuestas correspondientes. **No contiene lógica de negocio ni consultas directas a la base de datos**.

2. **Capa de Servicio (`internal/tareas/servicio.go`)**:
   - **Qué contiene**: Las reglas de negocio, validaciones de transiciones de estados y lógica de aplicación.
   - **Responsabilidad**: Orquesta las operaciones entre la capa de transporte y los repositorios. **No sabe nada de HTTP** (no maneja `http.ResponseWriter` ni códigos de estado directos).

3. **Capa de Almacenamiento / Repositorio (`internal/almacen/`)**:
   - **Qué contiene**: La implementación del acceso a datos combinando GORM y consultas optimizadas generadas por `sqlc` (`tareas_sqlc.go`).
   - **Responsabilidad**: Ejecuta las consultas SQL contra la base de datos PostgreSQL. Es la única capa que interactúa directamente con la persistencia.

### ¿Por qué elegimos esta estructura?
Permite aislar completamente la base de datos de los controladores HTTP mediante la interfaz de repositorio. Si el día de mañana cambiamos de motor de base de datos o de framework HTTP, los cambios se limitan estrictamente a su capa correspondiente sin romper las reglas del negocio ni los controladores.