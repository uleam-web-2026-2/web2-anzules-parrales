package main
import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"bitfamily-api/internal/almacen"
	"bitfamily-api/internal/tareas"
	"bitfamily-api/internal/transporte"
)

func main() {
	// Reemplaza esto con tu URL de conexión real de PostgreSQL
	dsn := "host=localhost user=postgres password=dhylan12 dbname=bitfamily port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { log.Fatal("no se pudo conectar: ", err) }

	err = db.AutoMigrate(&tareas.Categoria{}, &tareas.Tarea{})
	if err != nil { log.Fatal("no se pudo migrar: ", err) }

	// Semilla básica obligatoria para probar
	db.FirstOrCreate(&tareas.Categoria{ID: 1, Nombre: "Hogar"})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	repo := &almacen.Tareas{DB: db}
	servicio := tareas.NuevoServicio(repo)
	(&transporte.Tareas{Servicio: servicio}).Rutas(r)

	log.Println("Servidor escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}