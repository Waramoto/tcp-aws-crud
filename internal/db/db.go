package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"tcp-aws-crud/config"
)

type DB struct {
	client *dynamodb.Client
	cfg    config.AWS
}

func New(ctx context.Context, cfg config.AWS) (*DB, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &DB{
		client: dynamodb.NewFromConfig(awsCfg),
		cfg:    cfg,
	}, nil
}

func (db *DB) CreateItem(ctx context.Context, id, data string) error {
	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: db.cfg.DynamoDB.TableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to read item: %w", err)
	}
	if result.Item != nil {
		return fmt.Errorf("item already exists")
	}

	_, err = db.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: db.cfg.DynamoDB.TableName,
		Item: map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberS{Value: id},
			"data": &types.AttributeValueMemberS{Value: data},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

func (db *DB) ReadItem(ctx context.Context, id string) (string, error) {
	result, err := db.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: db.cfg.DynamoDB.TableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to read item: %w", err)
	}
	if result.Item == nil {
		return "", fmt.Errorf("item not found")
	}

	itemStrBuilder := strings.Builder{}
	itemStrBuilder.WriteString("Item:")
	for k, v := range result.Item {
		itemStrBuilder.WriteString(fmt.Sprintf("\n\t%s: %v", k, v))
	}

	return itemStrBuilder.String(), nil
}

func (db *DB) UpdateItem(ctx context.Context, id, data string) error {
	_, err := db.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: db.cfg.DynamoDB.TableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET data = :data"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":data": &types.AttributeValueMemberS{Value: data},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (db *DB) DeleteItem(ctx context.Context, id string) error {
	_, err := db.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: db.cfg.DynamoDB.TableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
