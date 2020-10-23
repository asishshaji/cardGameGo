package models

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Card struct {
	IMDB           float32 `json:"imdbrating" bson:"imdb,omitempty"`
	Rotten         float32 `json:"rottenrating" bson:"rotten,omitempty"`
	StateAwards    int32   `json:"stateawards" bson:"stateawards,omitempty"`
	NationalAwards int32   `json:"nationalawards" bson:"nationalawards,omitempty"`
}

func (c *Card) CreateCard(db *mongo.Database) error {

	fmt.Println("Hey")
	podcastsCollection := db.Collection("Cards")

	episodeResult, err := podcastsCollection.InsertMany(context.Background(), []interface{}{
		bson.D{
			{"title", "GraphQL for API Development"},
			{"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
			{"duration", 25},
		},
		bson.D{
			{"title", "Progressive Web Application Development"},
			{"description", "Learn about PWA development with Tara Manicsic."},
			{"duration", 32},
		},
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf("Inserted %v documents into episode collection!\n", len(episodeResult.InsertedIDs))

	return nil
}
