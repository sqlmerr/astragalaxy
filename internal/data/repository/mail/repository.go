package mail_repository

import (
	"context"

	"github.com/google/uuid"
	database "github.com/sqlmerr/astragalaxy/internal/data/postgres/database/sqlc"
	postgres_pool "github.com/sqlmerr/astragalaxy/internal/data/postgres/pool"
	"github.com/sqlmerr/astragalaxy/internal/model"
)

type MailRepository interface {
	CreateMailPackage(ctx context.Context, data CreatePackage) (model.MailPackage, error)
	GetMailPackage(ctx context.Context, id uuid.UUID) (model.MailPackage, error)
	GetMailPackagesByFromAgent(ctx context.Context, fromAgentID *uuid.UUID) ([]model.MailPackage, error)
	GetMailPackagesByToAgent(ctx context.Context, toAgentID uuid.UUID) ([]model.MailPackage, error)
	GetMailPackagesCountByFromAgent(ctx context.Context, fromAgentID *uuid.UUID) (int, error)
	DeleteMailPackage(ctx context.Context, id uuid.UUID) error
}

type MailRepositoryImpl struct {
	q  database.Queries
	db postgres_pool.DBTx
}

func New(q database.Queries, db postgres_pool.DBTx) *MailRepositoryImpl {
	return &MailRepositoryImpl{q: q, db: db}
}
