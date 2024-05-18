package util

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDatabase() *sql.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	connectionString := os.Getenv("DATABASE_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("database connection string is not set correctly")
		return nil
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

// util for testing //

type TestDatabase struct {
	DB *sql.DB
}

func (td TestDatabase) CreateTestTable() error {
	_, b, _, _ := runtime.Caller(0)
	utilPath := filepath.Dir(b)
	packagePath := filepath.Dir(utilPath)
	serverPath := filepath.Dir(packagePath)
	root := filepath.Dir(serverPath)
	infraPath := filepath.Join(root, "infra")
	migrationUpDirectory := filepath.Join(infraPath, "migration/up")
	files, err := os.ReadDir(migrationUpDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := migrationUpDirectory + "/" + file.Name()
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		_, err = td.DB.Exec(string(fileContent))
		if err != nil {
			return err
		}
	}

	return nil
}

func (td TestDatabase) DropTestTable() error {
	_, b, _, _ := runtime.Caller(0)
	utilPath := filepath.Dir(b)
	packagePath := filepath.Dir(utilPath)
	serverPath := filepath.Dir(packagePath)
	root := filepath.Dir(serverPath)
	infraPath := filepath.Join(root, "infra")
	migrationDownDirectory := filepath.Join(infraPath, "migration/down")
	files, err := os.ReadDir(migrationDownDirectory)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := migrationDownDirectory + "/" + file.Name()
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		_, err = td.DB.Exec(string(fileContent))
		if err != nil {
			return err
		}
	}

	return nil
}

func (td TestDatabase) DeleteItemsFromTable(tableName string) error {
	_, err := td.DB.Exec("DELETE FROM " + tableName)
	return err
}

func (td TestDatabase) Close() error {
	return td.DB.Close()
}

func NewTestDatabase() (TestDatabase, error) {
	_, b, _, _ := runtime.Caller(0)
	utilPath := filepath.Dir(b)
	packagePath := filepath.Dir(utilPath)
	serverPath := filepath.Dir(packagePath)
	envPath := filepath.Join(serverPath, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		return TestDatabase{}, err
	}

	connectionString := os.Getenv("DATABASE_CONNECTION_STRING")
	if connectionString == "" {
		errMessage := "database connection string is not set correctly"
		return TestDatabase{}, errors.New(errMessage)
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return TestDatabase{}, err
	}

	return TestDatabase{DB: db}, nil
}
