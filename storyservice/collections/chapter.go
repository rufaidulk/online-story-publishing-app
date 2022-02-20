package collections

import (
	"context"
	"storyservice/adapters"
	"storyservice/helper"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chapter struct {
	Id               primitive.ObjectID `bson:"_id"`
	UserUuid         string             `bson:"user_uuid"`
	StoryId          primitive.ObjectID `bson:"story_id"`
	Title            string             `bson:"title"`
	Body             string             `bson:"body"`
	PromotionalTitle string             `bson:"promotional_title"`
	PromotionalImage string             `bson:"promotional_image"`
	Rating           int8               `bson:"rating"`
	ReadCount        int64              `bson:"read_count"`
}

func NewChapter() *Chapter {
	return &Chapter{}
}

func (c *Chapter) CreateDocument() error {
	coll := getChapterCollection()
	c.Id = primitive.NewObjectID()
	_, err := coll.InsertOne(context.TODO(), &c)
	if err != nil {
		return err
	}

	return nil
}

func getChapterCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(ChapterCollection)

	return coll
}
