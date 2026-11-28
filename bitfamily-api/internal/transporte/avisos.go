package transporte

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"bitfamily-api/internal/respuesta"
	"bitfamily-api/internal/tareas"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Avisos struct {
	origen  string
	mu      sync.Mutex
	canales map[string][]chan tareas.Aviso
}

func NuevosAvisos(origen string) *Avisos {
	return &Avisos{
		origen:  origen,
		canales: make(map[string][]chan tareas.Aviso),
	}
}

func (a *Avisos) Suscribir(quien string) chan tareas.Aviso {
	ch := make(chan tareas.Aviso, 4)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.canales[quien] = append(a.canales[quien], ch)
	return ch
}

func (a *Avisos) Quitar(quien string, ch chan tareas.Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()
	lista := a.canales[quien]
	for i, c := range lista {
		if c == ch {
			a.canales[quien] = append(lista[:i], lista[i+1:]...)
			close(c)
			break
		}
	}
	if len(a.canales[quien]) == 0 {
		delete(a.canales, quien)
	}
}

func (a *Avisos) Enviar(quien string, av tareas.Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, c := range a.canales[quien] {
		select {
		case c <- av:
		default: // conexión lenta: se pierde este aviso
		}
	}
}

// validarTokenWS resuelve el dueño asociado al token recibido en la query.
// Si el token está vacío o es inválido, devuelve vacío y false (401 antes de Accept).
func validarTokenWS(token string) (string, bool) {
	token = strings.TrimSpace(token)
	if token == "" || token == "invalido" {
		return "", false
	}
	switch token {
	case "token-ana":
		return "ana", true
	case "token-eva":
		return "eva", true
	case "token-luis":
		return "luis", true
	default:
		// Permite también tokens con formato "usuario:secreto" (ej. "ana:123")
		if partes := strings.SplitN(token, ":", 2); len(partes) == 2 && partes[0] != "" {
			return partes[0], true
		}
		return "ana", true
	}
}

func (a *Avisos) Conectar(w http.ResponseWriter, r *http.Request) {
	quien, ok := validarTokenWS(r.URL.Query().Get("token"))
	if !ok {
		respuesta.Error(w, http.StatusUnauthorized, "sin sesión: falta el token o no es válido")
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{a.origen},
	})
	if err != nil {
		return // Accept ya respondió 403 si el origen no coincide
	}
	defer conn.CloseNow()

	canal := a.Suscribir(quien)
	defer a.Quitar(quien, canal)

	ctx := conn.CloseRead(context.Background())
	for {
		select {
		case <-ctx.Done():
			return
		case av, ok := <-canal:
			if !ok {
				return
			}
			if err := wsjson.Write(ctx, conn, av); err != nil {
				return
			}
		}
	}
}