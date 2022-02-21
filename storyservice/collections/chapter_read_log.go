package collections

import (
	"context"
	"storyservice/adapters"
	"storyservice/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChapterReadLog struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserUuid  string             `bson:"user_uuid"`
	StoryId   primitive.ObjectID `bson:"story_id"`
	ChapterId primitive.ObjectID `bson:"chapter_id"`
}

func NewChapterReadLog() *ChapterReadLog {
	return &ChapterReadLog{}
}

func (c *ChapterReadLog) UpsertDocument() (isFirstTimeRead bool, err error) {
	coll := getChapterReadLogCollection()
	data := bson.D{
		{"user_uuid", c.UserUuid},
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
		return
	}

	if result.MatchedCount == 0 {
		isFirstTimeRead = true
	}

	return
}

func getChapterReadLogCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(ChapterReadLogCollection)

	return coll
}
