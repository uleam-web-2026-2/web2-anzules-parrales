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


semana 11
### D5 (Conexión)
La API ubica a la base de datos utilizando el nombre del servicio definido en Docker Compose (`db`) como host en la cadena de conexión. Esto es posible gracias a la red interna que crea Docker Compose, permitiendo la resolución de nombres entre contenedores sin necesidad de utilizar `host.docker.internal`.

### D6 (Orquestación y orden)
Se implementó la directiva `depends_on` en el servicio de la API para indicar que depende de la base de datos. Esto asegura que Docker ordene el inicio de los contenedores correctamente, levantando primero a PostgreSQL.

### D7 (Persistencia)
Se configuró un volumen persistente (`pgdata`) en el servicio de la base de datos. Esto garantiza que la información almacenada en las tablas de tareas y categorías no se pierda al detener o eliminar los contenedores.

### D8 (Secretos)
El manejo de secretos y credenciales (como el usuario, la contraseña y la URL de conexión) se configuró mediante variables de entorno inyectadas directamente en la sección `environment` del archivo `compose.yaml`.

### D9 (Puertos)
Se expuso el puerto `8090` para que la API sea accesible desde el sistema operativo anfitrión (permitiendo realizar las peticiones desde Bruno), y se mapeó el puerto `5433:5432` para la base de datos, evitando conflictos con instancias locales preexistentes de PostgreSQL.