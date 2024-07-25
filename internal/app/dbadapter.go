package app

import (
	"context"
	"fmt"
	"log"

	"github.com/VadimGossip/concoleChat-auth/internal/model"
	"github.com/VadimGossip/concoleChat-auth/internal/repository"
	"github.com/VadimGossip/concoleChat-auth/internal/repository/user"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type DBAdapter struct {
	cfg      model.DbCfg
	db       *pgx.Conn
	userRepo repository.UserRepository
}

func NewDBAdapter(cfg model.DbCfg) *DBAdapter {
	return &DBAdapter{cfg: cfg}
}

func (d *DBAdapter) Connect(ctx context.Context) error {
	dbDSN := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", d.cfg.Host, d.cfg.Port, d.cfg.Name, d.cfg.Username, d.cfg.Password, d.cfg.SSLMode)
	db, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err = db.Ping(ctx); err != nil {
		return err
	}
	d.db = db
	d.userRepo = user.NewRepository(d.db)

	return nil
}

func (d *DBAdapter) Disconnect(ctx context.Context) error {
	if err := d.db.Close(ctx); err != nil {
		logrus.Errorf("Error occured on db connection close: %s", err.Error())
	}

	return nil
}
