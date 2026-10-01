# Hito 1 · Ficha del negocio

**Pareja:** Dhylan Parrales · Luis Anzules
**Paralelo:** Aplicaciones Web II B
**Negocio en una línea:** Plataforma web de gestión de tareas domésticas y micro-recompensas en dólares para motivar y educar financieramente a los hijos en el hogar.

## 1. Negocio de referencia

**Enlace:** https://www.starterstory.com/ideas/family-chore-app
App de hábitos para familias. Cobra suscripción freemium de $2.99/mes para acceder a funciones avanzadas de reporte familiar. Cifra que declara el fundador: $18,000 MRR (Ingreso Mensual Recurrente).

## 2. Caso de contraste

**Fuente:** https://failory.com/choresapp-post-mortem
Negocio: ChoresApp (Se estancó y redujo su base de usuarios activos un 40%). Hipótesis: Fracasó por carecer de un sistema de estímulo económico real y tangible (dólares) y por no involucrar un flujo de aprobación estricto por parte de los padres, volviéndolo un juego aburrido sin valor práctico.

## 3. Adaptación al Ecuador

1. Alta informalidad y uso masivo de efectivo y transferencias. Las familias en Ecuador manejan dinero físico diario; la app debe permitir registrar abonos manuales o saldos exactos en la cuenta virtual del niño.
2. Cultura de mesada y ahorro familiar en dólares. Los niños no entienden de criptomonedas abstractas; se debe manejar un saldo directo en dólares reales ($ USD) para que el premio sea tangible en permisos o compras reales.
3. Conectividad intermitente en sectores periféricos de Manta. Se requiere que la arquitectura de la aplicación almacene temporalmente las tareas en el caché local para validarlas aunque la red falle por un momento.

**Qué cambió en el modelo por estas restricciones:** Se eliminó cualquier moneda digital ficticia, reemplazándola por un atributo numérico decimal estricto (`saldo_usd`). Esto permite que cada tarea completada y aprobada abone centavos de dólar reales al perfil del hijo.

## 4. Modelo de datos

### Entidad: Usuario

| Atributo | Tipo | Obligatorio | Ejemplo |
|----------|------|-------------|---------|
| id | número entero | sí | 1 |
| nombre | texto | sí | Juan |
| rol | uno de: padre, madre, hijo | sí | hijo |
| saldo_usd | número decimal | sí | 4.50 |

### Entidad: Tarea

| Atributo | Tipo | Obligatorio | Ejemplo |
|----------|------|-------------|---------|
| id | número entero | sí | 1042 |
| titulo | texto | sí | Limpiar su cuarto |
| recompensa_usd | número decimal | sí | 0.50 |
| estado | uno de: disponible, completada, aprobada | sí | disponible |
| asignado_a | referencia a Usuario | sí | 1 |

### Entidad: HistorialMovimiento

| Atributo | Tipo | Obligatorio | Ejemplo |
|----------|------|-------------|---------|
| id | número entero | sí | 55 |
| descripcion | texto | sí | Pago por limpiar cuarto |
| monto_usd | número decimal | sí | 0.50 |
| fecha | fecha y hora | sí | 2026-09-30 10:32 |
| usuario_id | referencia a Usuario | sí | 1 |

### Relaciones

| Entidades | Cardinalidad | Frase |
|-----------|--------------|-------|
| Usuario - Tarea | 1 - N | Un usuario tiene muchas tareas asignadas |
| Usuario - HistorialMovimiento | 1 - N | Un usuario registra muchos movimientos de saldo |
| Tarea - HistorialMovimiento | 1 - N | Una tarea genera un movimiento de historial al pagarse |

### Structs en Go

```go
type Usuario struct {
    ID       uint    `gorm:"primaryKey" json:"id"`
    Nombre   string  `gorm:"not null" json:"nombre"`
    Rol      string  `gorm:"not null" json:"rol"`
    SaldoUSD float64 `gorm:"type:decimal(10,2);default:0.0" json:"saldo_usd"`
}

type Tarea struct {
    ID            uint    `gorm:"primaryKey" json:"id"`
    Titulo        string  `gorm:"not null" json:"titulo"`
    RecompensaUSD float64 `gorm:"type:decimal(10,2);not null" json:"recompensa_usd"`
    Estado        string  `gorm:"not null;default:disponible" json:"estado"`
    AsignadoA     uint    `json:"asignado_a"`
}

type HistorialMovimiento struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Descripcion string    `gorm:"not null" json:"descripcion"`
    MontoUSD    float64   `gorm:"type:decimal(10,2);not null" json:"monto_usd"`
    Fecha       time.Time `json:"fecha"`
    UsuarioID   uint      `json:"usuario_id"`
}