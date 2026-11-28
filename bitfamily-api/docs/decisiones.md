# Decisiones de Arquitectura – BitFamily API
## Semana 13 los avisos en tiempo real

### D1 Por qué WebSocket y no SSE ni sondeo
En `bitfamily-api` avisamos cuando cambia el estado de una `Tarea` (`estado_cambiado`), en un flujo principalmente desde el servidor hacia el cliente conectado. El sondeo (`polling`) no basta porque obliga a consultar periódicamente y retrasa el aviso según el intervalo configurado; SSE (`Server-Sent Events`) es unidireccional y depende de `EventSource` del navegador, mientras que WebSocket (`github.com/coder/websocket`) funciona de forma bidireccional y tiene soporte estándar en clientes web, móviles o de escritorio. Lo que nos cuesta el WebSocket es mantener conexiones abiertas en memoria, sincronizar el mapa de canales con `sync.Mutex` y validar token y origen en el *handshake*.

### D2 Quién recibe cada aviso
La llave del mapa `canales` en `Avisos` (`internal/transporte/avisos.go`) es el identificador/dueño (`string`) de quien le corresponde el aviso, y el dueño del registro sale del campo `Dueno` de `tareas.Tarea` (tomado del token/sesión al crear y no del cuerpo). Para que coincidan, tanto al registrar el canal en `Conectar` como al llamar `s.avisos.Enviar(res.Dueno, ...)` en `internal/tareas/servicio.go` (después de `s.repo.Guardar(&res)`), se usa el mismo valor de dueño. Si una conexión está lenta, el `select` con `default` en `Enviar` descarta ese aviso puntual para no bloquear al resto de suscriptores; si hubiera dos copias de la API, cada instancia tendría su propio mapa en memoria y los avisos no cruzarían entre instancias sin un bus/broker compartido.

### D3 El token en la dirección
La ruta `GET /ws` va fuera del guardia de cabecera `Authorization` porque el API WebSocket estándar del navegador no permite enviar cabeceras arbitrarias de autenticación al abrir el socket, y el token viaja en la query (`?token=...`). Revisamos el token antes de `websocket.Accept` para responder `401` con `respuesta.Error` sin abrir el socket si falta o es inválido, y luego `Accept` valida `OriginPatterns` con `origenFront` (`http://localhost:5173`) para devolver `403` si el origen no coincide. El origen sigue siendo necesario aunque haya token para prevenir conexiones *cross-site* desde orígenes no autorizados. En `docs/taller13/logs_avisos.txt` el token queda registrado en la URL de la petición `GET /ws?token=...` (por lo que se reemplaza por `token=BORRADO` antes del commit), y si un token venciera mientras la conexión ya está abierta, el socket ya aceptado no se corta automáticamente hasta que se cierre o se revalide en una nueva conexión.
## Semana 12 la seguridad de la API

### D1 Lo que encontramos y cerramos
Al revisar `docs/taller12/auditoria.md` (generado con `git grep -n -E -f docs/taller12/patrones.txt "*.go" ":(exclude)*_test.go"`), se identificaron las líneas que manipulan consultas o comandos SQL y las que son falsos positivos (mensajes de error o configuración). Las consultas reales se ejecutan mediante `sqlc` (con parámetros posicionales `$1`) y GORM (con marcadores `?` e identificadores numéricos parseados previamente con `strconv.Atoi`), sin concatenar datos del usuario con `+` ni `fmt.Sprintf`. Por ello no se requirieron cambios en las consultas SQL, y en `docs/taller12/inyeccion.txt` se comprobó que las entradas maliciosas no alteran la lógica ni exponen datos adicionales.

### D2 El origen que puede leer nuestras respuestas
Se configuró `AllowedOrigins` con `http://localhost:5173` (el origen del frontend de desarrollo) en lugar de `"*"` siguiendo el principio de mínimo privilegio para evitar que sitios no autorizados lean respuestas desde el navegador. El middleware `cors.Handler` se colocó arriba de todo, justo después de `r := chi.NewRouter()` y antes de otros middlewares y rutas, para que las peticiones previas `OPTIONS` se respondan de inmediato sin ser bloqueadas ni requerir autenticación. En `docs/taller12/cors_curl_ejecuta.txt` se observa que `curl` desde `http://malo.example` sí envía y recibe respuesta porque CORS es una política que aplican los navegadores, no un cliente de línea de comandos como `curl`. Cuando el frontend pase a producción, se deberá reemplazar `http://localhost:5173` por el origen `https://` definitivo de producción.

### D3 El límite del login
Se configuró el rate limit por IP con `httprate` (tope de 5 peticiones por ventana de 1 minuto) junto con `middleware.ClientIPFromRemoteAddr` y `porIP`, aplicándolo de forma acotada y no como `r.Use(...)` global sobre toda la API para no interrumpir las pruebas normales en Bruno. En BitFamily aún no existe una ruta `POST /login` (las rutas activas corresponden a categorías y tareas), por lo que este criterio queda definido para proteger el endpoint de inicio de sesión cuando se implemente, devolviendo `429 Too Many Requests` y la cabecera `Retry-After` al exceder el límite.

### D4 Detrás de un proxy
`middleware.ClientIPFromRemoteAddr` obtiene la IP desde `r.RemoteAddr` (la dirección remota de la conexión TCP directa) y la almacena en el contexto. Si la API queda detrás de un proxy inverso o balanceador de carga sin una configuración adecuada de confianza, `r.RemoteAddr` corresponderá a la IP del proxy y no a la de cada cliente real, haciendo que todos compartan el mismo contador y se bloqueen entre sí. Aun así, no debe confiarse ciegamente en `X-Forwarded-For` porque cualquier cliente puede falsificar esa cabecera para saltarse el límite; solo debe consultarse si el servidor verifica explícitamente que la petición proviene de un proxy de confianza.
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