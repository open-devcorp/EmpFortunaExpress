package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	uc "fortuna-express-web/pkg/domain/usecases"
	repository "fortuna-express-web/pkg/repositories"
	"fortuna-express-web/pkg/web"
	handlers "fortuna-express-web/pkg/web/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Importar el driver de SQLite
)

// NewDatabase crea y retorna una nueva conexión a la base de datos SQLite en memoria
func NewDatabase() (*sql.DB, error) {
	// Conectar a la base de datos SQLite en memoria
	connStr := ":memory:" // Usar SQLite en memoria

	// Conectar a la base de datos SQLite
	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verificar si la conexión es exitosa
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the SQLite in-memory database")

	// Crear las tablas si no existen
	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// createTables crea las tablas necesarias en la base de datos SQLite
func createTables(db *sql.DB) error {
	// SQL para crear las tablas
	liquidationsQuery := `CREATE TABLE IF NOT EXISTS liquidations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		departure TEXT NOT NULL,
		arrival TEXT NOT NULL,
		laundry REAL NOT NULL,
		garage REAL NOT NULL,
		guardianship REAL NOT NULL,
		cover REAL NOT NULL,
		sweeper REAL NOT NULL,
		driver TEXT NOT NULL,
		fuel REAL NOT NULL,
		date DATETIME NOT NULL,
		freight REAL NOT NULL,
		freight_liquid REAL NOT NULL,
		detraction REAL NOT NULL,
		gremission TEXT NOT NULL,
		gtransport TEXT NOT NULL,
		gtransport2 TEXT NOT NULL,
		invoice TEXT NOT NULL,
		driver_pay REAL NOT NULL,
		drive_description TEXT NOT NULL,
		fuel_description TEXT NOT NULL,
		liquid_trip REAL NOT NULL,
		truck TEXT NOT NULL,
		expense_total REAL NOT NULL,
		toll REAL NOT NULL,
		gast_adition BOOLEAN NOT NULL
	);`

	additionsQuery := `CREATE TABLE IF NOT EXISTS aditions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		price REAL NOT NULL,
		liquidation_id INTEGER NOT NULL,
		FOREIGN KEY (liquidation_id) REFERENCES liquidations(id)
		ON DELETE CASCADE
	);`

	// Ejecutar el script SQL para crear las tablas
	_, err := db.Exec(liquidationsQuery)
	if err != nil {
		return fmt.Errorf("failed to execute create liquidations table query: %w", err)
	}

	_, err = db.Exec(additionsQuery)
	if err != nil {
		return fmt.Errorf("failed to execute create additions table query: %w", err)
	}

	log.Println("Tables created or already exist")
	return nil
}

func main() {
	// Crear una instancia de Gin
	r := gin.Default()

	// Inicializar el logger
	logger := slog.Default()

	// Inicializar la base de datos SQLite en memoria
	database, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Crear repositorios y pasar el logger y la base de datos
	liquidationRepo := repository.NewLiquidationRepository(logger, database)
	aditionRepo := repository.NewAditionRepository(logger, database)

	// Crear el caso de uso
	uc := uc.NewLiquidationUseCase(liquidationRepo, aditionRepo)

	// Crear el handler
	handler := handlers.NewLiquidationsHandler(logger, uc)

	// Configurar el router desde el paquete 'web'
	web.SetupRouter(r, handler)

	// Iniciar el servidor en el puerto 8080
	if err := r.Run(":8080"); err != nil {
		logger.Error("Could not start server", err)
	}
}
