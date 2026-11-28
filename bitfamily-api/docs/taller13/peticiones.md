# Peticiones y evidencias — Taller 13

| Petición | Antes (taller 12) | Con el WebSocket | Con el reto |
| :--- | :---: | :---: | :---: |
| GET Listar Categorias | 200 | 200 | 200 |
| GET Listar Tareas | 200 | 200 | 200 |
| POST Crear Tarea | 201 | 201 | 201 |
| GET Ver Tarea por ID | 200 | 200 | 200 |
| PUT Actualizar Tarea | 400 | 400 | 400 |
| PATCH Cambiar Estado | 422 | 422 | 422 |
| GET /ws (sin token) | — | 401 | 401 |
| GET /ws (otro origen) | — | 403 | 403 |
| GET /ws (token + origen válido) | — | 101 | 101 |

## Evidencias generadas (puerto 8090)

- `prueba_servicio.txt`: pruebas unitarias del servicio al cambiar estado (envía aviso al dueño en transición válida y no avisa en transición inválida).
- `ws_sin_token.txt`: respuesta `401 Unauthorized` al conectar a `/ws` sin token.
- `ws_otro_origen.txt`: respuesta `403 Forbidden` al conectar a `/ws` con un `Origin` distinto a `http://localhost:5173`.
- `ws_permitido.txt`: respuesta `101 Switching Protocols` al conectar a `/ws` con token válido y origen `http://localhost:5173`.
- `prueba_websocket.txt`: pruebas de conexión WebSocket (`401`, `403` y entrega filtrada por dueño).
- `cliente_ws.html`: cliente WebSocket con reconexión automática, filtro por usuario/dueño y refresco de tareas.