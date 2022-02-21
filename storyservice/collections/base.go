package collections

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CategoryCollection = "categories"
const StoryCollection = "stories"
const ChapterCollection = "chapters"
const ChapterUserRatingCollection = "chapter_user_ratings"
const ChapterReadLogCollection = "chapter_read_logs"

func ListTrendingStories(userUuid string, limit int64) ([]primitive.M, error) {
	coll := getStoryCollection()
	projection := bson.D{{"categories", 0}}
	opts := options.Find().SetLimit(limit).SetProjection(projection).SetSort(bson.D{{"avg_read_count", -1}})
	filter := bson.M{"user_uuid": bson.M{"$ne": userUuid}}
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
