package connection

import (
	"belajar-golang-rest-api/lat/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func GetDatabase(conf config.Database) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable Timezone=%s",
		conf.Host,
		conf.Port,
		conf.User,
		conf.Pass,
		conf.Name,
		conf.Tz,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to opem connection:", err.Error())

	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping connection", err.Error())
	}
	return db

}
