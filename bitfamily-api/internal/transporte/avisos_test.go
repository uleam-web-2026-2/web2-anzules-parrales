package transporte

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bitfamily-api/internal/tareas"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func TestWebSocketSinTokenDa401(t *testing.T) {
	av := NuevosAvisos("http://localhost:5173")
	srv := httptest.NewServer(http.HandlerFunc(av.Conectar))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, resp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://localhost:5173"}},
	})
	if err == nil {
		t.Fatal("se esperaba error sin token")
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("se esperaba 401, se obtuvo: %v", resp)
	}
}

func TestWebSocketOtroOrigenDa403(t *testing.T) {
	av := NuevosAvisos("http://localhost:5173")
	srv := httptest.NewServer(http.HandlerFunc(av.Conectar))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws?token=token-ana"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, resp, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://malo.example"}},
	})
	if err == nil {
		t.Fatal("se esperaba error por origen no permitido")
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("se esperaba 403, se obtuvo: %v", resp)
	}
}

func TestAvisoLlegaSoloADueno(t *testing.T) {
	av := NuevosAvisos("http://localhost:5173")
	srv := httptest.NewServer(http.HandlerFunc(av.Conectar))
	defer srv.Close()

	wsBase := "ws" + strings.TrimPrefix(srv.URL, "http")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	connAna, _, err := websocket.Dial(ctx, wsBase+"/ws?token=token-ana", &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://localhost:5173"}},
	})
	if err != nil {
		t.Fatalf("no se pudo conectar ana: %v", err)
	}
	defer connAna.CloseNow()

	connEva, _, err := websocket.Dial(ctx, wsBase+"/ws?token=token-eva", &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://localhost:5173"}},
	})
	if err != nil {
		t.Fatalf("no se pudo conectar eva: %v", err)
	}
	defer connEva.CloseNow()

	time.Sleep(100 * time.Millisecond)

	av.Enviar("ana", tareas.Aviso{
		Tipo:   "estado_cambiado",
		ID:     10,
		Estado: "en_progreso",
	})

	var recibido tareas.Aviso
	if err := wsjson.Read(ctx, connAna, &recibido); err != nil {
		t.Fatalf("ana no recibió el aviso: %v", err)
	}
	if recibido.Tipo != "estado_cambiado" || recibido.ID != 10 {
		t.Fatalf("aviso inesperado para ana: %+v", recibido)
	}

	evaCtx, evaCancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer evaCancel()
	var msgEva tareas.Aviso
	if err := wsjson.Read(evaCtx, connEva, &msgEva); err == nil {
		t.Fatalf("eva no debía recibir el aviso de ana, pero recibió: %+v", msgEva)
	}
}
