package repository

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ContentRepository interface {
	FindByID(ctx context.Context, chatID int64) (*domain.ObserverTrace, error)
	Create(ctx context.Context, chatID int64) error
	UpdateWithData(ctx context.Context, chatID int64, url string, data []byte) error
	UpdateWithUrl(ctx context.Context, chatID int64, url string) error
	UpdateWithSelectorAndActivate(ctx context.Context, chatID int64, cssSelector string) error
	Delete(ctx context.Context, chatID int64) error
}

type DynamoContentRepository struct {
	db              dynamodb.Client
	s3Client        s3.Client
	dynamoTableName string
	s3BucketName    string
}

func NewDynamoContentRepository(db dynamodb.Client, s3Client s3.Client, dynamoTableName string, s3BucketName string) *DynamoContentRepository {
	return &DynamoContentRepository{
		db:              db,
		s3Client:        s3Client,
		dynamoTableName: dynamoTableName,
		s3BucketName:    s3BucketName,
	}
}

func (r *DynamoContentRepository) FindByID(ctx context.Context, chatID int64) (*domain.ObserverTrace, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(r.dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ChatID": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
		},
		ConsistentRead:         aws.Bool(false),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	result, err := r.db.GetItem(ctx, params)

	if err != nil {
		fmt.Printf("Got error calling dynamodb GetItem: %s\n", err)
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, nil
	}

	var content *domain.ObserverTrace = &domain.ObserverTrace{}
	opt := func(opt *attributevalue.DecoderOptions) {
		opt.TagKey = "json"
	}
	err = attributevalue.UnmarshalMapWithOptions(result.Item, content, opt)

	if err != nil {
		return nil, err
	}

	if content.FileName != nil {
		s3Result, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(r.s3BucketName),
			Key:    content.FileName,
		})

		if err != nil {
			fmt.Printf("Got error calling s3 GetObject: %s\n", err)
			return nil, err
		}

		bytesData, err := ioutil.ReadAll(s3Result.Body)
		content.Data = &bytesData

		if err != nil {
			fmt.Printf("Got error reading bytes from s3 Body: %s\n", err)
			return nil, err
		}
	}

	return content, nil
}

func (r *DynamoContentRepository) Create(ctx context.Context, chatID int64) error {
	now := time.Now().UTC()
	formattedDate := now.Format(time.RFC3339)

	input := dynamodb.PutItemInput{
		TableName: aws.String(r.dynamoTableName),
		Item: map[string]types.AttributeValue{
			"ChatID":      &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
			"IsActive":    &types.AttributeValueMemberBOOL{Value: false},
			"CreatedDate": &types.AttributeValueMemberS{Value: formattedDate},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	_, err := r.db.PutItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling dynamodb PutItem: %s\n", err)
		return err
	}

	return nil
}

func (r *DynamoContentRepository) UpdateWithData(ctx context.Context, chatID int64, url string, data []byte) error {
	now := time.Now().UTC()
	formattedDate := now.Format(time.RFC3339)
	fileName := fmt.Sprintf("%x.html", md5.Sum([]byte(url)))

	input := dynamodb.UpdateItemInput{
		TableName: aws.String(r.dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ChatID": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
		},
		UpdateExpression: aws.String("SET URL = :url, FileName = :fileName, UpdatedDate = :updatedDate"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":url":         &types.AttributeValueMemberS{Value: url},
			":fileName":    &types.AttributeValueMemberS{Value: fileName},
			":updatedDate": &types.AttributeValueMemberS{Value: formattedDate},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	reader := bytes.NewReader(data)
	_, err := r.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.s3BucketName),
		Key:    aws.String(fileName),
		Body:   reader,
	})
	if err != nil {
		fmt.Printf("Got error calling s3 PutObject: %s\n", err)
		return err
	}

	_, err = r.db.UpdateItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling dynamodb PutItem: %s\n", err)
		return err
	}

	return nil
}

func (r *DynamoContentRepository) UpdateWithUrl(ctx context.Context, chatID int64, url string) error {
	now := time.Now().UTC()
	formattedDate := now.Format(time.RFC3339)

	input := dynamodb.UpdateItemInput{
		TableName: aws.String(r.dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ChatID": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
		},
		UpdateExpression: aws.String("SET URL = :url, UpdatedDate = :updatedDate"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":url":         &types.AttributeValueMemberS{Value: url},
			":updatedDate": &types.AttributeValueMemberS{Value: formattedDate},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	_, err := r.db.UpdateItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling dynamodb UpdateItem: %s\n", err)
		return err
	}

	return nil
}

func (r *DynamoContentRepository) UpdateWithSelectorAndActivate(ctx context.Context, chatID int64, cssSelector string) error {
	now := time.Now().UTC()
	formattedDate := now.Format(time.RFC3339)

	input := dynamodb.UpdateItemInput{
		TableName: aws.String(r.dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ChatID": &types.AttributeValueMemberN{Value: strconv.FormatInt(chatID, 10)},
		},
		UpdateExpression: aws.String("SET CssSelector = :cssSelector, IsActive = :isActive, UpdatedDate = :updatedDate"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cssSelector": &types.AttributeValueMemberS{Value: cssSelector},
			":updatedDate": &types.AttributeValueMemberS{Value: formattedDate},
			":isActive":    &types.AttributeValueMemberBOOL{Value: true},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	_, err := r.db.UpdateItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling dynamodb UpdateItem: %s\n", err)
		return err
	}

	return nil
}

func (r *DynamoContentRepository) Delete(ctx context.Context, chatID int64) error {
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
