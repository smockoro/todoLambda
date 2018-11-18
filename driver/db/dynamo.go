package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"

	"github.com/smockoro/todoLambda/domain"
)

var (
	TableName = os.Getenv("DYNAMO_TABLE")
	Region    = os.Getenv("REGION")
)

type DB struct {
	Instance *dynamodb.DynamoDB
}

func New() DB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(Region)}),
	)
	return DB{Instance: dynamodb.New(sess)}
}

func (d DB) GetItem(user, id interface{}) (interface{}, error) {
	item, err := d.Instance.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"User": {
				S: aws.String(user.(string)),
			},
			"Id": {
				S: aws.String(id.(string)),
			},
		}})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get item")
	}
	if item.Item == nil {
		return nil, nil
	}
	todo := &model.Todo{}
	err = dynamodbattribute.UnmarshalMap(item.Item, &todo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal item")
	}
	return todo, nil
}

func (d DB) GetItems(user interface{}) (interface{}, error) {
	params := &dynamodb.QueryInput{
		TableName: aws.String(TableName),
		ExpressionAttributeNames: map[string]*string{
			"#User":    aws.String("User"),
			"#Id":      aws.String("Id"),
			"#Subject": aws.String("Subject"),
			"#Status":  aws.String("Status"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
				S: aws.String(user.(string)),
			},
		},
		KeyConditionExpression: aws.String("#User = :user"),
		ProjectionExpression:   aws.String("#User, #Id, #Subject, #Status"),
		ConsistentRead:         aws.Bool(false),
		Limit:                  aws.Int64(10),
	}
	result, err := d.Instance.Query(params)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get item")
	}
	if result.Items == nil {
		return nil, nil
	}
	todos := make([]*model.Todo, 0)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &todos)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal item")
	}
	return todos, nil
}

func (d DB) PutItem(i interface{}) (interface{}, error) {
	av, err := dynamodbattribute.MarshalMap(i)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TableName),
	}
	item, err := d.Instance.PutItem(input)
	if err != nil {
		return nil, err
	}
	return item, nil
}
