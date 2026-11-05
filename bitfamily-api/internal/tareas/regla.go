package tareas
import (
	"slices"
	"strings"
)

var siguientes = map[string][]string{
	"pendiente":   {"en_progreso", "cancelada"},
	"en_progreso": {"completada", "cancelada"},
	"completada":  {},
	"cancelada":   {},
}

func PuedePasar(de, a string) bool { return slices.Contains(siguientes[de], a) }

func validar(t Tarea) string {
	if strings.TrimSpace(t.Descripcion) == "" { return "la descripcion es obligatoria" }
	if t.Recompensa < 1 { return "la recompensa debe ser mayor a 0" }
	if !EstadosValidos[t.Estado] { return "estado no permitido: " + t.Estado }
	return ""
}