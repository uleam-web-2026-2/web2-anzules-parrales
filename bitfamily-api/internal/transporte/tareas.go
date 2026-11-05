package transporte
import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"bitfamily-api/internal/respuesta"
	"bitfamily-api/internal/tareas"
	"github.com/go-chi/chi/v5"
)

type Tareas struct { Servicio *tareas.Servicio }

func (m *Tareas) Rutas(r chi.Router) {
	r.Get("/categorias", m.listarCategorias)
	r.Get("/tareas", m.listar)
	r.Post("/tareas", m.crear)
	r.Get("/tareas/{id}", m.verUna)
	r.Put("/tareas/{id}", m.actualizar)
	r.Patch("/tareas/{id}", m.cambiarEstado)
	r.Delete("/tareas/{id}", m.borrar)
}

func responderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, tareas.ErrNoExiste):
		respuesta.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, tareas.ErrInvalido), errors.Is(err, tareas.ErrTransicion):
		respuesta.Error(w, http.StatusUnprocessableEntity, err.Error())
	default:
		log.Println("error interno:", err)
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo completar")
	}
}

func leerID(r *http.Request) (uint, bool) {
	n, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || n <= 0 { return 0, false }
	return uint(n), true
}

func (m *Tareas) listar(w http.ResponseWriter, r *http.Request) {
	lista, err := m.Servicio.Listar(r.URL.Query().Get("estado"))
	if err != nil { responderError(w, err); return }
	respuesta.JSON(w, http.StatusOK, lista)
}

func (m *Tareas) crear(w http.ResponseWriter, r *http.Request) {
	var t tareas.Tarea
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado, err := m.Servicio.Crear(t)
	if err != nil { responderError(w, err); return }
	respuesta.JSON(w, http.StatusCreated, creado)
}

func (m *Tareas) verUna(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(r)
	if !ok { respuesta.Error(w, http.StatusBadRequest, "id inválido"); return }
	t, err := m.Servicio.Buscar(id)
	if err != nil { responderError(w, err); return }
	respuesta.JSON(w, http.StatusOK, t)
}

func (m *Tareas) actualizar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(r)
	if !ok { respuesta.Error(w, http.StatusBadRequest, "id inválido"); return }
	var cambios tareas.Tarea
	if err := json.NewDecoder(r.Body).Decode(&cambios); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "JSON inválido"); return
	}
	t, err := m.Servicio.Actualizar(id, cambios)
	if err != nil { responderError(w, err); return }
	respuesta.JSON(w, http.StatusOK, t)
}

func (m *Tareas) cambiarEstado(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(r)
	if !ok { respuesta.Error(w, http.StatusBadRequest, "id inválido"); return }
	var pedido struct {
		Estado string `json:"estado"`
	}
	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "JSON inválido"); return
	}
	t, err := m.Servicio.CambiarEstado(id, pedido.Estado)
	if err != nil { responderError(w, err); return }
	respuesta.JSON(w, http.StatusOK, t)
}

func (m *Tareas) borrar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(r)
	if !ok { respuesta.Error(w, http.StatusBadRequest, "id inválido"); return }
	if err := m.Servicio.Borrar(id); err != nil { responderError(w, err); return }
	respuesta.JSON(w, http.StatusOK, map[string]uint{"borrada": id})
}

func (m *Tareas) listarCategorias(w http.ResponseWriter, r *http.Request) {
	lista, err := m.Servicio.ListarCategorias()
	if err != nil { responderError(w, err); return }
	respuesta.JSON(w, http.StatusOK, lista)
}