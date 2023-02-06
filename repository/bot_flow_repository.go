package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BotFlowRepository interface {
	FindByChatID(ctx context.Context, chatID int64) (domain.BotFlowState, error)
	Save(ctx context.Context, chatID int64, state domain.BotFlowState) error
	Delete(ctx context.Context, chatID int64) error
}

type DynamoBotFlowRepository struct {
	db              dynamodb.Client
	dynamoTableName string
}

func NewDynamoBotFlowRepository(db dynamodb.Client, dynamoTableName string) *DynamoBotFlowRepository {
	return &DynamoBotFlowRepository{
		db:              db,
		dynamoTableName: dynamoTableName,
	}
}

func (r *DynamoBotFlowRepository) FindByChatID(ctx context.Context, chatID int64) (domain.BotFlowState, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(r.dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ChatID": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
		},
		ConsistentRead: aws.Bool(false),
	}

	result, err := r.db.GetItem(ctx, params)

	if err != nil {
		fmt.Printf("Got error calling dynamodb GetItem: %s\n", err)
		return "", err
	}

	if len(result.Item) == 0 {
		return domain.NotStarted, nil
	}

	state := result.Item["State"].(*types.AttributeValueMemberS).Value
	return domain.BotFlowState(state), nil
}

func (r *DynamoBotFlowRepository) Save(ctx context.Context, chatID int64, state domain.BotFlowState) error {
	now := time.Now().UTC()
	formattedDate := now.Format(time.RFC3339)

	input := dynamodb.PutItemInput{
		TableName: aws.String(r.dynamoTableName),
		Item: map[string]types.AttributeValue{
			"ChatID":      &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
			"UpdatedDate": &types.AttributeValueMemberS{Value: formattedDate},
			"State":       &types.AttributeValueMemberS{Value: string(state)},
		},
	}
	_, err := r.db.PutItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling dynamodb PutItem: %s\n", err)
		return err
	}

	return nil
}

func (r *DynamoBotFlowRepository) Delete(ctx context.Context, chatID int64) error {
	input := dynamodb.DeleteItemInput{
		TableName: aws.String(r.dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ChatID": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
		},
	}
	_, err := r.db.DeleteItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling dynamodb DeleteItem: %s\n", err)
		return err
	}

	return nil
}
