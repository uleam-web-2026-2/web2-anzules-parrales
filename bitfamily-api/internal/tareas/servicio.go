package tareas
import "fmt"

type Servicio struct { repo Repositorio }
func NuevoServicio(repo Repositorio) *Servicio { return &Servicio{repo: repo} }
func (s *Servicio) Listar(estado string) ([]Tarea, error) { return s.repo.Listar(estado) }
func (s *Servicio) Buscar(id uint) (Tarea, error) { return s.repo.Buscar(id) }
func (s *Servicio) Borrar(id uint) error { return s.repo.Borrar(id) }
func (s *Servicio) ListarCategorias() ([]Categoria, error) { return s.repo.ListarCategorias() }

func (s *Servicio) Crear(t Tarea) (Tarea, error) {
	t.Estado = "pendiente" // Fase 2: siempre se fija inicial
	if msj := validar(t); msj != "" { return t, fmt.Errorf("%w: %s", ErrInvalido, msj) }
	
	existe, err := s.repo.ExisteCategoria(t.CategoriaID)
	if err != nil { return t, err }
	if !existe { return t, fmt.Errorf("%w: la categoria no existe", ErrInvalido) }
	
	err = s.repo.Crear(&t)
	return t, err
}

func (s *Servicio) Actualizar(id uint, cambios Tarea) (Tarea, error) {
	if msj := validar(cambios); msj != "" { return Tarea{}, fmt.Errorf("%w: %s", ErrInvalido, msj) }
	t, err := s.repo.Buscar(id)
	if err != nil { return t, err }
	
	t.Descripcion, t.Recompensa = cambios.Descripcion, cambios.Recompensa
	err = s.repo.Guardar(&t)
	return t, err
}

func (s *Servicio) CambiarEstado(id uint, nuevo string) (Tarea, error) {
	if !EstadosValidos[nuevo] { return Tarea{}, fmt.Errorf("%w: estado no permitido: %s", ErrInvalido, nuevo) }
	t, err := s.repo.Buscar(id)
	if err != nil { return t, err }
	if !PuedePasar(t.Estado, nuevo) { return t, fmt.Errorf("%w: de %s no se puede pasar a %s", ErrTransicion, t.Estado, nuevo) }
	
	t.Estado = nuevo
	err = s.repo.Guardar(&t)
	return t, err
}