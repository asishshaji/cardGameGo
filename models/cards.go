package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Card type
type Card struct {
	ID             primitive.ObjectID `bson:"_id"`
	IMDB           float32            `json:"imdbrating" bson:"imdb,omitempty"`
	Rotten         float32            `json:"rottenrating" bson:"rotten,omitempty"`
	StateAwards    int32              `json:"stateawards" bson:"stateawards,omitempty"`
	NationalAwards int32              `json:"nationalawards" bson:"nationalawards,omitempty"`
}

// CreateCard will create and add a new card
// to the mongodb
func (c *Card) CreateCard(db *mongo.Database) error {

	cardsCollection := db.Collection("Cards")
	cardResult, err := cardsCollection.InsertOne(context.TODO(), c)

	if err != nil {
		return err
	}
	log.Println(cardResult)

	return nil
}
