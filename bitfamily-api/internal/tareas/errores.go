package tareas
import "errors"

var (
	ErrNoExiste   = errors.New("no existe la tarea")
	ErrInvalido   = errors.New("datos no válidos")
	ErrTransicion = errors.New("transición no permitida")
)