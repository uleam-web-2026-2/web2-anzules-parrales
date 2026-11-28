package tareas

// Aviso es lo que se le cuenta a una persona cuando cambia algo suyo.
// Sale como JSON por el WebSocket: el hito 3 pide al menos tipo e id.
type Aviso struct {
	Tipo   string `json:"tipo"`
	ID     uint   `json:"id"`
	Estado string `json:"estado"`
}

// Avisador es lo que el servicio necesita para avisar. El servicio no
// sabe si el aviso sale por WebSocket, por SSE o por correo: solo dice
// a quién y qué.
type Avisador interface {
	Enviar(quien string, av Aviso)
}