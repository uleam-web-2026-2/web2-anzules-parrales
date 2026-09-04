# Servidor Mínimo - Taller de Montaje de Entorno

Proyecto base desarrollado para la asignatura Aplicaciones Web II / Aplicación para el Servidor Web (ULEAM 2026-2).

---

## 👥 Integrantes
* **Dhylan Ariel Parrales Pachay**
* **LUIS ANZULES**

**Paralelo:** Web II (2026-2)  
**Repositorio:** `web2-anzules-parrales`

---

## 🛠️ Requisitos Previos

1. **Go:** Versión 1.25 o superior instalada.
2. **PostgreSQL:**
   * **Opción elegida:** Instalación nativa (PostgreSQL 18).
   * **Puerto:** `5432`
   * **Usuario:** `postgres`
   * **Base de datos requerida:** `mesa_ayuda`

---

## ⚙️ Configuración de la Base de Datos

Antes de ejecutar el servidor, asegúrate de tener creada la base de datos:

1. Abrir **pgAdmin** o la consola SQL de PostgreSQL (`psql`).
2. Crear la base de datos del curso ejecutando:
   ```sql
   CREATE DATABASE mesa_ayuda;












# [Nombre del producto]

> Reemplacen todo lo que está entre corchetes en su primer commit.
> Este README es la puerta de entrada del repositorio: en la semana 1 otra pareja debe poder levantar el servidor siguiendo solo lo que dice aquí, y desde la semana 4 es la base de la integración continua.

**Hilo Servidor · ULEAM · Período 2026-2**
Aplicaciones Web II (TDI-610) · Aplicación para el Servidor Web (IS-503)

## Integrantes

| Integrante | Usuario de GitHub | Paralelo |
| --- | --- | --- |
| [Apellidos Nombres] | [usuario] | [Servidor Web A / Web II A / Web II B] |
| [Apellidos Nombres] | [usuario] | [Servidor Web A / Web II A / Web II B] |

## El producto

[Una o dos líneas: qué negocio es y qué vende. Se completa cuando la ficha del negocio pase la compuerta de la semana 3. Mientras tanto, el dominio de trabajo es la mesa de ayuda.]

## Cómo levantar el servidor

```bash
go run .
```

El servidor responde en `http://localhost:8080`.

Requisitos: Go (versión estable) y PostgreSQL.

## Base de datos

- Opción usada por la pareja: [nativa / contenedor]
- Con contenedor: `docker run --name pg -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:16`
- Base de datos del proyecto: `[nombre]`

## Cómo correr las pruebas

```bash
go test ./...
```

Las pruebas corren también en la integración continua (pestaña Actions). Desde la semana 4, un entregable cuyas pruebas no pasan en la integración no se recibe.

## Convenciones del repositorio

- Un commit de cada integrante como mínimo por taller; el commit de cierre se hace en clase.
- Mensajes de commit: qué cambió y por qué, entendibles sin el autor presente.
- Uso de IA declarado en el cuerpo del commit: una línea con qué herramienta y para qué parte.
- Ningún secreto en el código ni en el historial: la configuración se externaliza (semana 4).

## Estructura

```
docs/ficha_negocio.md   → ficha del negocio (se entrega antes del día A de la semana 3)
main.go                 → punto de entrada del servidor
```
