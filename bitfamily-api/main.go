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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
