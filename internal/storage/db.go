package storage

// import 	(
// 	"database/sql"
// 	"fsender/config"
// )

// type Storage struct {
// 	Cfg *config.Config
// 	*sql.DB
// }

// func NewStorage(cfg *config.Config, db *sql.DB) *Storage {
// 	return &Storage{
// 		Cfg: cfg,
// 		DB: db,
// 	}
// }

// func (s *Storage) InitDB() {
// 	query := `
// 	CREATE TABLE IF NOT EXISTS users (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		file_path TEXT NOT NULL,
// 		email TEXT NOT NULL UNIQUE
//         )`

// 	_, err := s.DB.Exec(query)
// 	if err != nil {
// 		panic(err)
// 	}
// }
