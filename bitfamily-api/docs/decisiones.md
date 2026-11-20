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
## Semana 11 el contenedor

### D1 Qué entra en la imagen
La imagen final solo contiene el binario compilado estáticamente gracias al uso de una estructura de dos etapas con `gcr.io/distroless/static-debian13:nonroot`. Todo el código fuente, herramientas de Go y archivos sensibles se quedan en la etapa de compilación. Si no existiera el archivo `.dockerignore`, el archivo `.env` local se enviaría como parte del contexto de construcción hacia Docker, comprometiendo los secretos.

### D2 Dónde está la base, vista desde el contenedor
Falló `localhost` porque dentro del contenedor `localhost` hace referencia al propio contenedor aislado y no a la máquina host. Se solucionó utilizando `host.docker.internal` (o el nombre del servicio `db` en Compose). Para un servidor real de producción, se debe cambiar por la IP o dominio real del servidor de base de datos accesible por red.

### D3 Los puertos
El mapeo de puertos `-p 8090:8090` indica que el número de la izquierda (8090) corresponde al puerto expuesto en la máquina anfitriona (al que llama Bruno), y el de la derecha (8090) es el puerto interno donde escucha el servidor dentro del contenedor. Este número debe coincidir exactamente con la variable `PUERTO` definida en la configuración; de lo contrario, el tráfico externo no llegará al servidor.

### D4 La hora
Se decidió mantener la hora en formato UTC tal como la maneja por defecto el contenedor basado en distroless. En el contexto de negocio de BitFamily, esto afecta la sincronización de las marcas de tiempo (timestamps) de las tareas, por lo que la aplicación debe encargarse de formatear las zonas horarias si el cliente lo requiere explícitamente.

## Semana 11 reto: la API y la base, juntas

### D5 Cómo encuentra la API a la base
La API localiza a la base de datos utilizando el nombre del servicio (`db`) definido en el archivo `compose.yaml` a través de la red interna de Docker Compose, usando el puerto interno estándar `5432` y sin requerir el puerto `5433` de la máquina.

### D6 El orden de arranque
La primera vez con la base de datos vacía, la API podría intentar conectar antes de que Postgres esté listo. Esto se resolvió implementando la directiva `depends_on` junto con un bloque de `healthcheck` en la base de datos para asegurar que acepte conexiones antes de levantar la API.

### D7 Dónde quedan los datos
Los datos de las tablas se persisten gracias a la declaración de un volumen de Docker (`pgdata`). Esto previene que la información se pierda al ejecutar un `docker compose down`, y los datos solo se eliminan si se fuerza explícitamente un `docker compose down -v`.

### D8 Los secretos
Los secretos y credenciales de conexión no están escritos en texto plano dentro del archivo `compose.yaml`, sino que se inyectan dinámicamente mediante variables de entorno y archivos de configuración seguros del entorno local.

### D9 Lo que se publica en la máquina
El archivo `compose.yaml` publica únicamente el puerto de la API (`8090:8090`) hacia el anfitrión. La base de datos **no publica** ningún puerto hacia afuera (se elimina el mapeo `5433:5432`), asegurando que solo sea accesible de forma interna por la API. Si se requiere consultar la base con `sqlc` o `psql` externamente, se realiza ejecutando comandos directos dentro del contenedor de la base (`docker exec`).