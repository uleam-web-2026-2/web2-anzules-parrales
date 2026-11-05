package respuesta
import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, codigo int, datos any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(datos)
}

func Error(w http.ResponseWriter, codigo int, mensaje string) {
	JSON(w, codigo, map[string]string{"error": mensaje})
}