# Hoja del Hito 2 – BitFamily API

## 1. Usuarios de ejemplo
| Usuario | Rol | Contraseña de ejemplo |
| :--- | :--- | :--- |
| dhylan | padre | dhylan-2026 |
| luis | hijo | luis-2026 |

## 2. Las seis peticiones
| # | Petición | Con qué credencial | Debe responder |
| :-: | :--- | :--- | :--- |
| 1 | GET /tareas | Ninguna | 401 |
| 2 | POST /login con dhylan | Ninguna | 200 y el token |
| 3 | GET /tareas | Token de Dhylan | 200 |
| 4 | PATCH /tareas/1 (estado inválido o sin permisos) | Token de Luis (hijo) | 403 |
| 5 | PATCH /tareas/1 (cambio a en_progreso) | Token de Dhylan (padre) | 200 |
| 6 | PATCH /tareas/1 (transición prohibida) | Token de Dhylan (padre) | 422 |

## 3. Matriz ruta x rol
| Ruta | Sin sesión | Hijo | Padre |
| :--- | :--- | :--- | :--- |
| POST /login | Sí | Sí | Sí |
| GET /categorias | No (401) | Sí | Sí |
| GET /tareas | No (401) | Sí | Sí |
| POST /tareas | No (401) | No (403) | Sí |
| GET /tareas/{id} | No (401) | Sí | Sí |
| PUT /tareas/{id} | No (401) | No (403) | Sí |
| PATCH /tareas/{id} | No (401) | No (403) | Sí |
| DELETE /tareas/{id} | No (401) | No (403) | Sí |

## 4. Recuadro de cumplimiento
- [x] El validador termina sin errores ni avisos
- [x] Los structs de todas las entidades están en el código y sus tablas existen
- [x] Las seis peticiones están en la colección, en orden, y cada una responde lo que dice esta hoja
- [x] go test ./... en verde, con una prueba de integración
- [x] El código de la entidad principal está en capas y la estructura está explicada en docs/decisiones.md
- [x] Hay una captura de respaldo de cada paso