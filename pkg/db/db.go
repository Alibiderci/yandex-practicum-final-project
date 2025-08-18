// Пакет db управляет всем взаимодействием с базой данных SQLite.
// Он включает инициализацию соединения, создание схемы и закрытие соединения.
package db

import (
	"os"
	"database/sql"

	_ "modernc.org/sqlite"
)

// DB - это глобальный пул соединений с базой данных.
// Он инициализируется при старте приложения функцией Init() и используется
// во всех функциях для работы с БД.
var DB *sql.DB

// schema определяет структуру таблиц базы данных.
// Запрос выполняется только один раз при первом создании файла БД.
const schema = `CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(128) NOT NULL,
	comment TEXT,
	repeat VARCHAR(64) NOT NULL
);

CREATE INDEX IF NOT EXISTS date_index ON scheduler (date)`

// Close корректно закрывает соединение с базой данных, если оно было открыто.
func Close() error {

	if DB != nil {
		return DB.Close()
	}
	return nil
}

// Init инициализирует соединение с файлом базы данных SQLite.
// Если файл БД не существует, он будет создан вместе с необходимой схемой таблиц.
func Init() error {

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db" // Файл БД по умолчанию.
	}

	// Проверяем, существует ли файл БД, чтобы определить, нужно ли создавать схему.
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	// Открываем (или создаем) файл базы данных.
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если база данных была только что создана, выполняем запрос на создание схемы.
	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
