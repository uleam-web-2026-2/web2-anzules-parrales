# Auditoría de Inyección SQL — BitFamily API

## Revisión de hallazgos (`auditoria.txt`)

| Archivo y línea | Fragmento | ¿Arma SQL? | ¿Llega texto del cliente? | ¿Cómo viaja el valor? | Veredicto y acción tomada |
| :--- | :--- | :---: | :---: | :--- | :--- |
| `internal/almacen/filas/consultas.sql.go:17` | `q.db.QueryRowContext(ctx, buscarTarea, id)` | Sí (consulta estática `sqlc`) | Sí (`id` de ruta) | Parámetro posicional (`$1`) | **Seguro.** Código generado por `sqlc` con placeholders; sin cambios requeridos. |
| `internal/almacen/filas/consultas.sql.go:42` | `q.db.ExecContext(ctx, guardarTarea, ...)` | Sí (consulta estática `sqlc`) | Sí (campos del JSON) | Parámetros posicionales (`$1`, `$2`, ...) | **Seguro.** Código generado por `sqlc` con placeholders; sin cambios requeridos. |
| `internal/almacen/filas/db.go:13-16` | Definición de interfaz `DBTX` (`ExecContext`, `QueryContext`, `QueryRowContext`) | No (solo firma de interfaz) | No | N/A | **Falso positivo.** Interfaz generada por `sqlc`; sin cambios requeridos. |
| `internal/almacen/tareas.go:14` | `consulta = consulta.Where("estado = ?", estado)` | Sí (GORM) | Sí (`estado` por query param) | Parámetro (`?`) | **Seguro.** GORM vincula el valor mediante placeholder `?`; sin cambios requeridos. |
| `internal/almacen/tareas.go:15` | `err := consulta.Find(&lista).Error` | Sí (GORM) | Indirecto (usa `consulta` previa) | Parámetros previos (`?`) | **Seguro.** Ejecuta la consulta parametrizada construida en la línea 14. |
| `internal/almacen/tareas.go:21` | `err := a.DB.First(&t, id).Error` | Sí (GORM) | Sí (`id` numérico de ruta) | Parámetro (`?`) | **Seguro.** GORM parametriza la llave primaria; sin cambios requeridos. |
| `internal/almacen/tareas.go:37` | `a.DB.Model(&tareas.Categoria{}).Where("id = ?", id).Count(&n)` | Sí (GORM) | Sí (`categoria_id` del body) | Parámetro (`?`) | **Seguro.** Filtrado parametrizado con `?`; sin cambios requeridos. |
| `internal/almacen/tareas.go:43` | `err := a.DB.Preload("Tareas").Find(&lista).Error` | Sí (GORM) | No | Sin entrada de usuario | **Seguro.** Consulta fija sin parámetros externos; sin cambios requeridos. |
| `internal/tareas/regla.go:19` | `return "estado no permitido: " + t.Estado` | No | Sí (`estado` del JSON) | Concatenación de mensaje de error | **Falso positivo.** Construye un texto de error de dominio en Go, nunca una consulta SQL. |
| `internal/transporte/tareas.go:44` | `m.Servicio.Listar(r.URL.Query().Get("estado"))` | No | Sí (`estado` por query param) | Argumento de función en Go | **Falso positivo.** Solo extrae el parámetro HTTP y lo delega al servicio/almacén donde se parametriza. |
| `main.go:62` | `http.ListenAndServe(":"+puerto, r)` | No | No (variable de entorno) | Concatenación de dirección TCP | **Falso positivo.** Define el puerto de escucha HTTP del servidor. |

## Conclusión de la auditoría
No se encontraron concatenaciones de entrada del cliente ni uso de `fmt.Sprintf` para armar consultas SQL. Todas las operaciones contra PostgreSQL utilizan consultas parametrizadas mediante `sqlc` (`$1`, `$2`, ...) y GORM (`?`), por lo que no fue necesario modificar código de consultas en la capa de almacén.