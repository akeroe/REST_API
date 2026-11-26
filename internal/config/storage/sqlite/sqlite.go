package sqlite

import (
	"database/sql"
	"fmt"

	"url-shorneter/internal/config/storage"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New" //  в константе op(operation) хранится имя функции new

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err) // врап ошибки
	}
	// stmt - подготовленный запрос
	stmt, err := db.Prepare(` 
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES(?, ?)") // stmt подготавливает запрос на
	// вставку данных в таблицу url,  ? - места для значений
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close() // stmt.Close() закрывает подготовленный запрос после выполнения функции SaveURL

	res, err := stmt.Exec(alias, urlToSave) // stmt.Exec выполняет подготовленный запрос с подставленными значениями alias и urlToSave
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId() // res.LastInsertId() возвращает ID последней вставленной записи
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil

}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ? ")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string                           // в эту переменную положим возвращаемый URL
	err = stmt.QueryRow(alias).Scan(&resURL) // stmt.QueryRow(alias) выполняет запрос к базе данных с подставленным значением alias и возвращает одну строку результата,
	// Scan(&resURL) сканирует результат запроса и сохраняет значение в переменную resURL
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)

	}
	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) (int64,error) {
const op = "storage.sqlite.DeleteURL" 
result, err := s.db.Exec("DELETE FROM url WHERE alias = ?", alias) // s.db.Exec выполняет запрос к базе данных на удаление записи с подставленным значением alias
if err != nil {
	return 0, fmt.Errorf("%s: get rows affected: %w", op, err)
}

rowsAffected, err := result.RowsAffected() // считывает количество удаленных строк из результата выполнения запроса
if err != nil {
	return 0, fmt.Errorf("%s: get rows affected: %w", op, err)
}

return rowsAffected, nil

}