package utils
import(
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn

func InitDB(dsn string) {
	var err error
	dbConn, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("Unable to connect database: %v", err)
	}
}

func GetDBConn() *pgx.Conn {
	return dbConn
}
