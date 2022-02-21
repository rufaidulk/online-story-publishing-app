package collections

import (
	"context"
	"log"
	"math"
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

func CalculateAvgRatingOfSingleChapter(chapterId primitive.ObjectID) (float64, error) {
	coll := getChapterUserRatingCollection()
	matchStage := bson.D{
		{"$match", bson.D{{"chapter_id", chapterId}}},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$chapter_id"},
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

func getChapterUserRatingCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(ChapterUserRatingCollection)

	return coll
}
