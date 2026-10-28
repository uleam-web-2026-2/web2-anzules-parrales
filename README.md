# Servidor Mesa de Ayuda - Semana 2

Servidor HTTP desarrollado en Go con el enrutador Chi para la gestión centralizada de tickets de soporte técnico.

---

## 👥 Integrantes

| Integrante | Usuario de GitHub | Paralelo |
| --- | --- | --- |
| Dhylan Ariel Parrales Pachay | [usuario] | Web II |
| Luis Manuel Anzules | [usuario] | Web II |

**Repositorio:** `web2-anzules-parrales`  
**Módulo:** `github.com/uleam-web-2026-2/ab-anzules-parrales`

---

## 🚀 Cómo levantar el servidor

1. Asegúrate de tener Go instalado (v1.22 o superior).
2. Instala/sincroniza las dependencias:
   ```bash
   go mod tidy
   ```
3. Inicia el servidor:
   ```bash
   go run main.go
   ```

El servidor responderá en `http://localhost:8080`.

> **Nota de almacenamiento:** En esta etapa los tickets viven en memoria interna (protegidos contra condiciones de carrera mediante `sync.Mutex`) y se reinician al apagar la aplicación.

---

## 🧪 Cómo correr las pruebas del contrato

Ejecuta el script de pruebas automatizadas con Git Bash o en un entorno Linux/macOS:

```bash
chmod +x pruebas.sh
./pruebas.sh
```

Para probar contra un servidor externo o en red local:
```bash
./pruebas.sh http://IP_DEL_SERVIDOR:8080
```

---

## 📋 Tabla de Rutas y Errores

| Método | Ruta | Éxito | Errores posibles |
| --- | --- | --- | --- |
| `GET` | `/tickets` | `200 OK` (lista de tickets) | - |
| `POST` | `/tickets` | `201 Created` (ticket creado) | `400 Bad Request` (`cuerpo_malformado`) <br> `422 Unprocessable Entity` (`datos_invalidos`) |
| `GET` | `/tickets/{id}` | `200 OK` (detalle del ticket) | `400 Bad Request` (`id_invalido`) <br> `404 Not Found` (`no_encontrado`) |
| `GET` | `/explotar` | - | `500 Internal Server Error` (`error_interno`) |
| `*` | Ruta inexistente | - | `404 Not Found` (`ruta_inexistente`) |
| `*` | Método no soportado | - | `405 Method Not Allowed` (`metodo_no_permitido`) |

---

## 🔍 Revisión Cruzada

Revisión cruzada: Sin hallazgos — auto-evaluado con `pruebas.sh` (100% de casos del contrato verificados).