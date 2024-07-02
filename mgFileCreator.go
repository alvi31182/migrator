package migrator

import (
	"fmt"
	"os"
	"time"
)

func GenerateFileName() string {
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
	fileName := GenerateFileName()
	filePath := fmt.Sprintf("migrations/%s", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("не удалось создать файл миграции: %w", err)
	}
	defer file.Close()

	fmt.Printf("Файл миграции создан: %s\n", filePath)
	return nil
}
