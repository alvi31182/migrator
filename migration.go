package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func CreateMigrationTable(db *sql.DB) error {
	fmt.Println("Подключение к базе данных:", db)
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS migration_version (
			id SERIAL PRIMARY KEY,
			version BIGINT NOT NULL,
			executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу миграций: %w", err)
	}
	return nil
}

func GetLatestMigrationVersion(db *sql.DB) (int64, error) {
	var version int64
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM migration_version").Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить версию последней миграции", err)
	}
	return version, nil
}

func ApplyMigration(db *sql.DB, filePath string, version int64) error {
	// Открываем файл миграции
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл миграции: %s: %w", filePath, err)
	}
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("не удалось применить миграцию: %s: %w", filePath, err)
	}

	_, err = db.Exec("INSERT INTO migration_version (version) VALUES ($1)", version)
	if err != nil {
		return fmt.Errorf("не удалось добавить версию миграции в базу данных: %w", err)
	}

	fmt.Printf("Миграция %s успешно применена\n", filePath)
	return nil
}

func generateFileName() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("Version%d.sql", timestamp)
}

func CreateMigrationFile() error {
	// Убедимся, что папка migrations существует
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		err := os.Mkdir("migrations", os.ModePerm)
		if err != nil {
			return fmt.Errorf("не удалось создать папку migrations: %w", err)
		}
	}

	// Создаем файл миграции
	fileName := generateFileName()
	filePath := fmt.Sprintf("migrations/%s", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл миграции: %w", err)
	}
	defer file.Close()

	fmt.Printf("Файл миграции создан: %s\n", filePath)
	return nil
}
