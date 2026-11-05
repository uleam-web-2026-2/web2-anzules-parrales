package tareas

type Categoria struct {
	ID     uint    `json:"id"`
	Nombre string  `json:"nombre"`
	Tareas []Tarea `json:"tareas,omitempty"`
}

// Fase 2: La entidad ya no tiene etiquetas json
type Tarea struct {
	ID          uint
	CategoriaID uint
	Descripcion string
	Recompensa  int
	Estado      string
}

var EstadosValidos = map[string]bool{
	"pendiente": true, "en_progreso": true, "completada": true, "cancelada": true,
}