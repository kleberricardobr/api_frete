package database

import (
	"api_frete/models"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var Conn *sql.DB

type Database struct {
	Config *models.ConfigModel
}

func (d *Database) OpenPostgres() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		d.Config.Db.Host,
		strconv.Itoa(d.Config.Db.Port),
		d.Config.Db.User,
		d.Config.Db.Pass,
		d.Config.Db.Name,
	)

	var err error
	Conn, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	if err := Conn.Ping(); err != nil {
		return err
	}

	return nil
}

func (d *Database) RunMigrations() error {
	if Conn == nil {
		return fmt.Errorf("conexão com banco não foi estabelecida")
	}

	driver, err := postgres.WithInstance(Conn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("erro ao criar driver de migração: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("erro ao inicializar migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("erro ao executar migrations: %v", err)
	}

	fmt.Println("Migrations executadas com sucesso!")
	return nil
}

func (d *Database) ClosePostgres() {
	Conn.Close()
}
