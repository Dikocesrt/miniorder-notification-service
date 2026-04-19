package utils

import (
	"context"
	"database/sql"
	"fmt"
	"miniorder-order-service/internal/domain"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func InitConnectDatabase(svcName string, debugMode bool, dbConfig domain.DatabaseConfig) *bun.DB {
	addr := fmt.Sprintf("%s:%d", dbConfig.IP, dbConfig.Port)
	fmt.Printf("Connecting to Database : '%v'\n", addr)
	openDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(addr),
		pgdriver.WithUser(dbConfig.User),
		pgdriver.WithPassword(dbConfig.Password),
		pgdriver.WithDatabase(dbConfig.Name),
		pgdriver.WithApplicationName(svcName),
		pgdriver.WithTLSConfig(nil),
		pgdriver.WithDialTimeout(dbConfig.DialTimeout),
		// pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
	))
	db := bun.NewDB(openDB, pgdialect.New())
	if err := db.PingContext(context.Background()); err != nil {
		panic(err)
	}
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxOpenConns)
	if debugMode {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	fmt.Printf("Connected to Database : '%v'\n", addr)
	return db
}
