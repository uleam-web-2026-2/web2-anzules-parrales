package tareas

type Categoria struct {
	ID     uint    `json:"id"`
	Nombre string  `json:"nombre"`
	Tareas []Tarea `json:"tareas,omitempty"`
}

type Tarea struct {
	ID          uint   `json:"id"`
	CategoriaID uint   `json:"categoria_id"`
	Descripcion string `json:"descripcion"`
	Recompensa  int    `json:"recompensa"`
	Estado      string `json:"estado"`
}

var EstadosValidos = map[string]bool{
	"pendiente": true, "en_progreso": true, "completada": true, "cancelada": true,
}