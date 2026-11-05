package almacen
import (
	"errors"
	"bitfamily-api/internal/tareas"
	"gorm.io/gorm"
)

type Tareas struct { DB *gorm.DB }
var _ tareas.Repositorio = (*Tareas)(nil)

func (a *Tareas) Listar(estado string) ([]tareas.Tarea, error) {
	var lista []tareas.Tarea
	consulta := a.DB
	if estado != "" { consulta = consulta.Where("estado = ?", estado) }
	err := consulta.Find(&lista).Error
	return lista, err
}

func (a *Tareas) Buscar(id uint) (tareas.Tarea, error) {
	var t tareas.Tarea
	err := a.DB.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { return t, tareas.ErrNoExiste }
	return t, err
}

func (a *Tareas) Crear(t *tareas.Tarea) error { return a.DB.Create(t).Error }
func (a *Tareas) Guardar(t *tareas.Tarea) error { return a.DB.Save(t).Error }
func (a *Tareas) Borrar(id uint) error {
	resultado := a.DB.Delete(&tareas.Tarea{}, id)
	if resultado.Error != nil { return resultado.Error }
	if resultado.RowsAffected == 0 { return tareas.ErrNoExiste }
	return nil
}

func (a *Tareas) ExisteCategoria(id uint) (bool, error) {
	var n int64
	err := a.DB.Model(&tareas.Categoria{}).Where("id = ?", id).Count(&n).Error
	return n > 0, err
}

func (a *Tareas) ListarCategorias() ([]tareas.Categoria, error) {
	var lista []tareas.Categoria
	err := a.DB.Preload("Tareas").Find(&lista).Error
	return lista, err
}