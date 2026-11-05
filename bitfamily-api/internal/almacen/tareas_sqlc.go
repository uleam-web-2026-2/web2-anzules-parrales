package almacen

import (
	"context"
	"database/sql"
	"errors"
	"bitfamily-api/internal/almacen/filas"
	"bitfamily-api/internal/tareas"
)

type TareasSqlc struct {
	*Tareas
	Conexion *sql.DB
	q        *filas.Queries
}

var _ tareas.Repositorio = (*TareasSqlc)(nil)

func NuevoTareasSqlc(g *Tareas, c *sql.DB) *TareasSqlc {
	return &TareasSqlc{Tareas: g, Conexion: c, q: filas.New(c)}
}

func (a *TareasSqlc) Buscar(id uint) (tareas.Tarea, error) {
	f, err := a.q.BuscarTarea(context.Background(), int64(id))
	if errors.Is(err, sql.ErrNoRows) { return tareas.Tarea{}, tareas.ErrNoExiste }
	if err != nil { return tareas.Tarea{}, err }
	return tareas.Tarea{
		ID: uint(f.ID), CategoriaID: uint(f.CategoriaID),
		Descripcion: f.Descripcion, Recompensa: int(f.Recompensa), Estado: f.Estado,
	}, nil
}

func (a *TareasSqlc) Guardar(t *tareas.Tarea) error {
	n, err := a.q.GuardarTarea(context.Background(), filas.GuardarTareaParams{
		ID: int64(t.ID), CategoriaID: int64(t.CategoriaID),
		Descripcion: t.Descripcion, Recompensa: int64(t.Recompensa), Estado: t.Estado,
	})
	if err == nil && n == 0 { return tareas.ErrNoExiste }
	return err
}
