package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

func obtenerOrigenesPermitidos() []string {
	raw := os.Getenv("CORS_ORIGINS")
	if strings.TrimSpace(raw) == "" {
		return []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"http://localhost:8090",
		}
	}

	partes := strings.Split(raw, ",")
	origenes := make([]string, 0, len(partes))
	for _, p := range partes {
		limpio := strings.TrimSpace(p)
		if limpio != "" {
			origenes = append(origenes, limpio)
		}
	}
	if len(origenes) == 0 {
		return []string{"http://localhost:5173"}
	}
	return origenes
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
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Fase 2: Configuración CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   obtenerOrigenesPermitidos(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Fase 3: Rate Limiting por IP (ej. 60 peticiones por minuto)
	r.Use(httprate.LimitByIP(60, time.Minute))

	// Obtener la conexión subyacente para sqlc
	conexion, err := db.DB()
	if err != nil {
		log.Fatal("no se pudo obtener la conexion: ", err)
	}

	// Inicializamos el repositorio de sqlc que incrusta a GORM
	repo := almacen.NuevoTareasSqlc(&almacen.Tareas{DB: db}, conexion)
	servicio := tareas.NuevoServicio(repo)
	(&transporte.Tareas{Servicio: servicio}).Rutas(r)

	// 4. Leemos el puerto desde las variables de entorno
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8080" // Valor por defecto
	}

	log.Println("escuchando en el puerto", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, r))
}
