package main

import (
	"log"
	"net/http"
	"os"

	"bitfamily-api/internal/almacen"
	"bitfamily-api/internal/tareas"
	"bitfamily-api/internal/transporte"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const origenFront = "http://localhost:5173"

func porIP(r *http.Request) (string, error) {
	ip := middleware.GetClientIP(r.Context())
	return httprate.CanonicalizeIP(ip), nil
}

func main() {
	// 1. Leemos la cadena de conexión desde las variables de entorno
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("configuración: falta DATABASE_URL")
	}

	// 2. Conectamos a la base de datos
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar: ", err)
	}

	err = db.AutoMigrate(&tareas.Categoria{}, &tareas.Tarea{})
	if err != nil {
		log.Fatal("no se pudo migrar: ", err)
	}

	// Semilla básica obligatoria para probar
	db.FirstOrCreate(&tareas.Categoria{ID: 1, Nombre: "Hogar"})

	// 3. Rutas y Middlewares
	r := chi.NewRouter()

	// CORS arriba de todo: responde las peticiones previas (OPTIONS).
	// Mínimo privilegio: se nombra el front, no "*".
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{origenFront},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.ClientIPFromRemoteAddr)

	// Obtener la conexión subyacente para sqlc
	conexion, err := db.DB()
	if err != nil {
		log.Fatal("no se pudo obtener la conexion: ", err)
	}

	// Inicializamos el repositorio de sqlc que incrusta a GORM
	repo := almacen.NuevoTareasSqlc(&almacen.Tareas{DB: db}, conexion)
	servicio := tareas.NuevoServicio(repo)

	avisos := transporte.NuevosAvisos(origenFront)
	servicio.AvisarCon(avisos)

	(&transporte.Tareas{Servicio: servicio}).Rutas(r)
	r.Get("/ws", avisos.Conectar)

	// 4. Leemos el puerto desde las variables de entorno
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8080" // Valor por defecto
	}

	log.Println("escuchando en el puerto", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, r))
}
