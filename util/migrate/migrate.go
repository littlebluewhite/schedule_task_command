package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"log"
	"path/filepath"
	"runtime"
	"schedule_task_command/util/config"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	rootPath string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath = filepath.Dir(filepath.Dir(filepath.Dir(b)))
}

type Migration struct {
	client *migrate.Migrate
}

func New(config config.DBConfig) *Migration {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		config.User, config.Password, config.Host, config.Port, config.DB)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("could not connect to PostgreSQL database... %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}
	sourceURL := fmt.Sprintf("file://%s/migrations", rootPath)
	fmt.Println(sourceURL)
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}
	return &Migration{
		client: m,
	}
}

func (m *Migration) To(targetVersion uint) {
	if err := m.client.Migrate(targetVersion); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	afterVersion, _, _ := m.client.Version()
	fmt.Printf("Migration to version:%d success", afterVersion)
}

func (m *Migration) Up() {
	if err := m.client.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	afterVersion, _, _ := m.client.Version()
	fmt.Printf("Migration up version:%d success", afterVersion)
}

func (m *Migration) Down() {
	if err := m.client.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	version, _, _ := m.client.Version()
	fmt.Printf("Migration down version:%d success", version)
}
