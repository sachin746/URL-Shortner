package database

import (
	"context"
	"time"

	"URL-Shortner/entities"
	"URL-Shortner/flags"
	"URL-Shortner/log"
	"URL-Shortner/utils/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Get() *gorm.DB {
	return db
}

func InitDatabase(ctx context.Context) {
	port := configs.Get().Database.Postgres.Port

	host := flags.DatabaseHost()
	pass := flags.DatabasePassword()
	dbname := flags.DatabaseName()
	user := flags.DatabaseUser()
	log.Sugar.Infof("Connecting to database %s on port %d with user %s", dbname, port, user)

	if host == "" || pass == "" {
		log.Sugar.Fatal("DATABASE_HOST and DATABASE_PASSWORD environment variables must be set")
	}

	dialector := postgres.Open("postgres://" + user + ":" + pass + "@" + host + "/" + dbname + "?sslmode=disable")

	var err error
	db, err = gorm.Open(
		dialector,
		&gorm.Config{
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
			NowFunc: func() time.Time {
				return time.Now().In(time.UTC)
			},
		},
	)
	if err != nil {
		log.Sugar.Fatalf("Failed to connect to database: %v", err)
	}
	log.Sugar.Info("Database connection established successfully")

	// Auto-migrate the schema
	log.Sugar.Info("Running database migrations...")
	err = db.AutoMigrate(
		&entities.User{},
		&entities.URLIDGenerator{},
		&entities.ShortenUrl{},
	)
	if err != nil {
		log.Sugar.Fatalf("Failed to migrate database: %v", err)
	}
	log.Sugar.Info("Database migration completed")
}
