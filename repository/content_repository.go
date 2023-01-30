package repository

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/EduardTruuvaart/web-observer/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ContentRepository interface {
	FindByID(ctx context.Context, url string) (*domain.ObserverTrace, error)
	Save(ctx context.Context, content domain.ObserverTrace) error
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

func (r *DynamoContentRepository) FindByID(ctx context.Context, url string) (*domain.ObserverTrace, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(r.dynamoTableName),
		Key: map[string]types.AttributeValue{
			"URL": &types.AttributeValueMemberS{Value: url},
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

	s3Result, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.s3BucketName),
		Key:    aws.String(content.FileName),
	})

	if err != nil {
		fmt.Printf("Got error calling s3 GetObject: %s\n", err)
		return nil, err
	}

	content.Data, err = ioutil.ReadAll(s3Result.Body)

	if err != nil {
		fmt.Printf("Got error reading bytes from s3 Body: %s\n", err)
		return nil, err
	}

	return content, nil
}

func (r *DynamoContentRepository) Save(ctx context.Context, content domain.ObserverTrace) error {
	now := time.Now().UTC()
	formattedDate := now.Format(time.RFC3339)
	content.FileName = fmt.Sprintf("%x.html", md5.Sum([]byte(content.URL)))

	reader := bytes.NewReader(content.Data)
	_, err := r.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.s3BucketName),
		Key:    aws.String(content.FileName),
		Body:   reader,
	})
	if err != nil {
		fmt.Printf("Got error calling s3 PutObject: %s\n", err)
		return err
	}

	input := dynamodb.PutItemInput{
		TableName: aws.String(r.dynamoTableName),
		Item: map[string]types.AttributeValue{
			"URL":         &types.AttributeValueMemberS{Value: content.URL},
			"UpdatedDate": &types.AttributeValueMemberS{Value: formattedDate},
			"FileName":    &types.AttributeValueMemberS{Value: content.FileName},
		},
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}
	_, err = r.db.PutItem(ctx, &input)

	if err != nil {
		fmt.Printf("Got error calling dynamodb PutItem: %s\n", err)
		return err
	}

	return nil
}
