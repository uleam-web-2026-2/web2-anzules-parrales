package tareas

import (
	"errors"
	"testing"
)

type avisadorFalso struct {
	enviados []struct {
		quien string
		aviso Aviso
	}
}

func (a *avisadorFalso) Enviar(quien string, av Aviso) {
	a.enviados = append(a.enviados, struct {
		quien string
		aviso Aviso
	}{quien: quien, aviso: av})
}

type repoMemoria struct {
	tareas map[uint]Tarea
}

func newRepoMemoria() *repoMemoria {
	return &repoMemoria{tareas: make(map[uint]Tarea)}
}

func (r *repoMemoria) Listar(estado string) ([]Tarea, error) {
	var out []Tarea
	for _, t := range r.tareas {
		if estado == "" || t.Estado == estado {
			out = append(out, t)
		}
	}
	return out, nil
}

func (r *repoMemoria) Buscar(id uint) (Tarea, error) {
	t, ok := r.tareas[id]
	if !ok {
		return Tarea{}, ErrNoExiste
	}
	return t, nil
}

func (r *repoMemoria) Crear(t *Tarea) error {
	if t.ID == 0 {
		t.ID = uint(len(r.tareas) + 1)
	}
	r.tareas[t.ID] = *t
	return nil
}

func (r *repoMemoria) Guardar(t *Tarea) error {
	if _, ok := r.tareas[t.ID]; !ok {
		return ErrNoExiste
	}
	r.tareas[t.ID] = *t
	return nil
}

func (r *repoMemoria) Borrar(id uint) error {
	if _, ok := r.tareas[id]; !ok {
		return ErrNoExiste
	}
	delete(r.tareas, id)
	return nil
}

func (r *repoMemoria) ExisteCategoria(id uint) (bool, error) {
	return id > 0, nil
}

func (r *repoMemoria) ListarCategorias() ([]Categoria, error) {
	return nil, nil
}

func TestCambiarEstadoAvisaAlDueno(t *testing.T) {
	repo := newRepoMemoria()
	av := &avisadorFalso{}
	srv := NuevoServicio(repo)
	srv.AvisarCon(av)

	creada, err := srv.Crear(Tarea{
		CategoriaID: 1,
		Dueno:       "ana",
		Descripcion: "Revisar presupuesto",
		Recompensa:  10,
	})
	if err != nil {
		t.Fatalf("error inesperado al crear: %v", err)
	}

	if _, err := srv.CambiarEstado(creada.ID, "en_progreso"); err != nil {
		t.Fatalf("error inesperado al cambiar estado: %v", err)
	}

	if len(av.enviados) != 1 {
		t.Fatalf("se esperaba 1 aviso, llegaron %d", len(av.enviados))
	}
	if av.enviados[0].quien != "ana" {
		t.Fatalf("se esperaba aviso para 'ana', llegó para %q", av.enviados[0].quien)
	}
	if av.enviados[0].aviso.Tipo != "estado_cambiado" || av.enviados[0].aviso.ID != creada.ID {
		t.Fatalf("aviso inesperado: %+v", av.enviados[0].aviso)
	}
}

func TestTransicionInvalidaNoAvisa(t *testing.T) {
	repo := newRepoMemoria()
	av := &avisadorFalso{}
	srv := NuevoServicio(repo)
	srv.AvisarCon(av)

	creada, err := srv.Crear(Tarea{
		CategoriaID: 1,
		Dueno:       "ana",
		Descripcion: "Tarea de prueba",
		Recompensa:  5,
	})
	if err != nil {
		t.Fatalf("error inesperado al crear: %v", err)
	}

	_, err = srv.CambiarEstado(creada.ID, "completada")
	if err == nil {
		t.Fatal("se esperaba error por transición inválida")
	}
	if !errors.Is(err, ErrTransicion) {
		t.Fatalf("se esperaba ErrTransicion, se obtuvo: %v", err)
	}
	if len(av.enviados) != 0 {
		t.Fatalf("no debían enviarse avisos, llegaron %d", len(av.enviados))
	}
}
