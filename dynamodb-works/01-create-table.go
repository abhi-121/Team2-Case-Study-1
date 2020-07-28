package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region: aws.String("ap-south-1"),
	}))

	db := dynamodb.New(sess)

	awsParams := &dynamodb.CreateTableInput{
		TableName: aws.String("Music3"),
		KeySchema: [] *dynamodb.KeySchemaElement{
			{AttributeName: aws.String("Artist"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("SongTitle"), KeyType: aws.String("RANGE")},
		},

		AttributeDefinitions: [] *dynamodb.AttributeDefinition{
			{AttributeName: aws.String("Artist"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("SongTitle"), AttributeType: aws.String("S")},
		},

		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(10),
			WriteCapacityUnits: aws.Int64(5),
		},
	}
	resp, err := db.CreateTable(awsParams)
	if err!= nil {
		log.Fatal("Sorry error creating table : %v", err)
	}
	fmt.Println(resp)
}
