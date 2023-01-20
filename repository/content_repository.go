package repository

import (
	"context"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const tableName = "Content"

type ContentRepository interface {
	FindByID(ctx context.Context) (*domain.Content, error)
}

type DynamoContentRepository struct {
	db dynamodb.Client
}

func NewDynamoContentRepository(db dynamodb.Client) *DynamoContentRepository {
	return &DynamoContentRepository{
		db: db,
	}
}

func (r *DynamoContentRepository) FindByID(ctx context.Context) (*domain.Content, error) {
	panic("not implemented")
}
