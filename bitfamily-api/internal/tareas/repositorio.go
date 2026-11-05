package tareas

type Repositorio interface {
	Listar(estado string) ([]Tarea, error)
	Buscar(id uint) (Tarea, error)
	Crear(t *Tarea) error
	Guardar(t *Tarea) error
	Borrar(id uint) error
	ExisteCategoria(id uint) (bool, error)
	ListarCategorias() ([]Categoria, error)
}