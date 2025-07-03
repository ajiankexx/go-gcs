package utils
import(
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

var (
	pgConn *pgx.Conn
	redisConn *redis.Client
)

func InitDB() {
	var err error
	pgConn, err = pgx.Connect(context.Background(), "postgres://admin:1234@localhost:5432/gcs_db")
	if err != nil {
		log.Fatalf("Unable to connect : %v", err)
	}
}

func InitRedis() {
	redisConn = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	if err := redisConn.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Unable to connect redis: %v", err)
	}
}

func GetDBConn() *pgx.Conn {
	if pgConn == nil {
		InitDB()
	}
	if pgConn == nil {
		log.Fatal("Postgres not initialized")
	}
	return pgConn
}

func GetRedisConn() *redis.Client {
	if redisConn == nil {
		InitRedis()
	}
	if redisConn == nil {
		log.Fatal("Redis not initialized")
	}
	return redisConn
}
