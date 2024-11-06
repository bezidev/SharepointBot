package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

type sqlImpl struct {
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func (db *sqlImpl) Init() {
	db.db.MustExec(schema)
}

type SQL interface {
	Init()

	GetSharepointNotification(id string) (notification SharepointNotification, err error)
	GetSharepointNotifications() (notification []SharepointNotification, err error)
	InsertSharepointNotification(notification SharepointNotification) (err error)
	UpdateSharepointNotification(notification SharepointNotification) error
	DeleteSharepointNotification(id string) error
}

func NewSQL(driver string, drivername string, logger *zap.SugaredLogger) (SQL, error) {
	db, err := sqlx.Connect(driver, drivername)
	return &sqlImpl{
		db:     db,
		logger: logger,
	}, err
}
