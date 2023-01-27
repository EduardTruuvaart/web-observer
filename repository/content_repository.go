package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "Content"

type ContentRepository interface {
	FindByID(ctx context.Context, url string) (*domain.Content, error)
	Save(ctx context.Context, content domain.Content) error
}

type DynamoContentRepository struct {
	db dynamodb.Client
}

func NewDynamoContentRepository(db dynamodb.Client) *DynamoContentRepository {
	return &DynamoContentRepository{
		db: db,
	}
}

func (r *DynamoContentRepository) FindByID(ctx context.Context, url string) (*domain.Content, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"URL": &types.AttributeValueMemberS{Value: url},
		},
	}

	result, err := r.db.GetItem(ctx, params)

	if err != nil {
		return nil, err
	}

	fmt.Printf("GetItem consumed units: %d\n", result.ConsumedCapacity.CapacityUnits)

	if len(result.Item) == 0 {
		return nil, nil
	}

	var content *domain.Content = &domain.Content{}
	opt := func(opt *attributevalue.DecoderOptions) {
		opt.TagKey = "json"
	}
	err = attributevalue.UnmarshalMapWithOptions(result.Item, content, opt)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func (r *DynamoContentRepository) Save(ctx context.Context, content domain.Content) error {
	now := time.Now().UTC()
	formattedDate := now.Format(time.RFC3339)

	input := dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"URL":         &types.AttributeValueMemberS{Value: content.URL},
			"Data":        &types.AttributeValueMemberB{Value: content.Data},
			"IsActive":    &types.AttributeValueMemberBOOL{Value: content.IsActive},
			"UpdatedDate": &types.AttributeValueMemberS{Value: formattedDate},
		},
	}

	result, err := r.db.PutItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling PutItem: %s\n", err)
		return err
	}

	fmt.Printf("PutItem consumed units: %d\n", result.ConsumedCapacity.CapacityUnits)

	return nil
}
