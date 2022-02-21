package collections

import (
	"context"
	"log"
	"storyservice/adapters"
	"storyservice/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChapterUserRating struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserUuid  string             `bson:"user_uuid"`
	StoryId   primitive.ObjectID `bson:"story_id"`
	ChapterId primitive.ObjectID `bson:"chapter_id"`
	Rating    int8               `bson:"rating"`
}

func NewChapterUserRating() *ChapterUserRating {
	return &ChapterUserRating{}
}

func (c *ChapterUserRating) UpsertDocument() error {
	coll := getChapterUserRatingCollection()

	data := bson.D{
		{"rating", c.Rating},
	}
	filter := bson.D{
		{"user_uuid", c.UserUuid},
		{"story_id", c.StoryId},
		{"chapter_id", c.ChapterId},
	}
	update := bson.D{{"$set", data}}
	opts := options.Update().SetUpsert(true)
	result, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	log.Println(result)
	return nil
}

func getChapterUserRatingCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(ChapterUserRatingCollection)

	return coll
}
