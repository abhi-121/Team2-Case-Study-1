package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:8000"),
		Region: aws.String("ap-south-1"),
	}))

	db := dynamodb.New(sess)

	music := Music{
		Artist: "Sonu Nigam",
		SongTitle: "Hello Again",
		AlbumTitle: "My Album",
		Awards: 2,
	}

	musicMap, err := dynamodbattribute.Marshal(music)
	if err != nil {
		panic("Cannot map the values given in music struct")
	}

	params := &dynamodb.PutItemInput{
		TableName: aws.String("Music3"),
		Item: musicMap,
	}

	resp, err := db.PutItem(params)

	if err != nil {
		panic("Error in putting the item")
	}
	fmt.Println(resp)
}

type Music struct {
	Artist string
	SongTitle string
	AlbumTitle string
	Awards int
}