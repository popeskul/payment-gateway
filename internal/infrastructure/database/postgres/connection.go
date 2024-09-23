package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/popeskul/payment-gateway/internal/core/ports"
	"github.com/popeskul/payment-gateway/internal/infrastructure/uuid"
)

type Database struct {
	Pool               *pgxpool.Pool
	MerchantRepository ports.MerchantRepository
	PaymentRepository  ports.PaymentRepository
	RefundRepository   ports.RefundRepository
	UserRepository     ports.UserRepository
}

func NewDatabase(config ports.DatabaseConfig) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.GetHost(), config.GetPort(), config.GetUser(), config.GetPassword(), config.GetDBName())

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	poolConfig.MaxConns = int32(config.GetMaxConnections())
	poolConfig.MinConns = int32(config.GetMinConnections())
	poolConfig.MaxConnLifetime = time.Duration(config.GetMaxConnLifetime()) * time.Second
	poolConfig.MaxConnIdleTime = time.Duration(config.GetMaxConnIdleTime()) * time.Second

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	db := &Database{Pool: pool}

	uuidGenerator := uuid.NewUUIDGenerator()

	db.MerchantRepository = NewMerchantRepository(db, uuidGenerator)
	db.PaymentRepository = NewPaymentRepository(db, uuidGenerator)
	db.RefundRepository = NewRefundRepository(db, uuidGenerator)
	db.UserRepository = NewUserRepository(db, uuidGenerator)

	return db, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}

func (db *Database) BeginTx(ctx context.Context) (ports.Transaction, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %v", err)
	}
	return &Transaction{tx: tx}, nil
}

type Transaction struct {
	tx pgx.Tx
}

func (t *Transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *Transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
