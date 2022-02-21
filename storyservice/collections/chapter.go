package collections

import (
	"context"
	"errors"
	"math"
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

func CalculateAvgRatingOfStory(storyId primitive.ObjectID) (float64, error) {
	coll := getChapterCollection()
	matchStage := bson.D{
		{"$match", bson.D{{"story_id", storyId}}},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$story_id"},
			{"avg_rating", bson.D{
				{"$avg", "$rating"},
			}},
		}},
	}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return 0, err
	}
	for _, result := range results {
		rating := result["avg_rating"].(float64)
		avgRating := math.Round(rating*10) / 10
		return avgRating, nil
	}

	return 0, nil
}

func getChapterCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(ChapterCollection)

	return coll
}
