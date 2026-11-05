package transporte
import "bitfamily-api/internal/tareas"

type TareaPeticion struct {
	CategoriaID uint   `json:"categoria_id"`
	Descripcion string `json:"descripcion"`
	Recompensa  int    `json:"recompensa"`
}

func (p TareaPeticion) aEntidad() tareas.Tarea {
	return tareas.Tarea{CategoriaID: p.CategoriaID, Descripcion: p.Descripcion, Recompensa: p.Recompensa}
}

type TareaRespuesta struct {
	ID          uint   `json:"id"`
	CategoriaID uint   `json:"categoria_id"`
	Descripcion string `json:"descripcion"`
	Recompensa  int    `json:"recompensa"`
	Estado      string `json:"estado"`
}

func nuevaRespuesta(t tareas.Tarea) TareaRespuesta {
	return TareaRespuesta{ID: t.ID, CategoriaID: t.CategoriaID, Descripcion: t.Descripcion, Recompensa: t.Recompensa, Estado: t.Estado}
}