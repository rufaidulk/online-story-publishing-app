package collections

import (
	"context"
	"errors"
	"storyservice/adapters"
	"storyservice/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chapter struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserUuid  string             `bson:"user_uuid"`
	StoryId   primitive.ObjectID `bson:"story_id"`
	Title     string             `bson:"title"`
	Body      string             `bson:"body"`
	Rating    float64            `bson:"rating"`
	ReadCount int64              `bson:"read_count"`
}

func NewChapter() *Chapter {
	return &Chapter{}
}

func (c *Chapter) LoadById(id string) error {
	coll := getChapterCollection()
	objId, _ := primitive.ObjectIDFromHex(id)
	err := coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&c)

	if err != nil && err == mongo.ErrNoDocuments {
		return errors.New("requested chapter not found")
	} else if err != nil {
		return err
	}

	return nil
}

func (c *Chapter) LoadByObjectId(id primitive.ObjectID) error {
	coll := getChapterCollection()
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&c)

	if err != nil && err == mongo.ErrNoDocuments {
		return errors.New("requested chapter not found")
	} else if err != nil {
		return err
	}

	return nil
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

func (c *Chapter) UpdateDocument() error {
	data := bson.D{
		{"title", c.Title},
		{"body", c.Body},
	}
	coll := getChapterCollection()
	filter := bson.D{{"_id", c.Id}}
	update := bson.D{{"$set", data}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chapter) UpdateRating() error {
	avgChapterRating, err := CalculateAvgRatingOfSingleChapter(c.Id)
	if err != nil {
		return err
	}
	coll := getChapterCollection()
	data := bson.D{
		{"rating", avgChapterRating},
	}
	filter := bson.D{{"_id", c.Id}}
	update := bson.D{{"$set", data}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chapter) DeleteDocument() error {
	coll := getChapterCollection()
	filter := bson.D{{"_id", c.Id}}
	_, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Chapter) IncrementReadCount() error {
	coll := getChapterCollection()
	data := bson.D{
		{"read_count", 1},
	}
	filter := bson.D{{"_id", c.Id}}
	update := bson.D{{"$inc", data}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
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
