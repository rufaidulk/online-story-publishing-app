package collections

import (
	"context"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CategoryCollection = "categories"
const StoryCollection = "stories"
const ChapterCollection = "chapters"
const ChapterUserRatingCollection = "chapter_user_ratings"
const ChapterReadLogCollection = "chapter_read_logs"

func ListMostRatedStories(userUuid string, limit int64) ([]primitive.M, error) {
	coll := getStoryCollection()
	projection := bson.D{{"categories", 0}}
	opts := options.Find().SetLimit(limit).SetProjection(projection).SetSort(bson.D{{"rating", -1}})
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

func ListCategories(limit int64) ([]primitive.M, error) {
	coll := getCategoryCollection()
	opts := options.Find().SetLimit(limit)
	filter := bson.M{}
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

func CalculateAvgReadCountOfStory(storyId primitive.ObjectID) (int64, error) {
	coll := getChapterCollection()
	matchStage := bson.D{
		{"$match", bson.D{{"story_id", storyId}}},
	}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$story_id"},
			{"avg_read_count", bson.D{
				{"$avg", "$read_count"},
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
		readCount := result["avg_read_count"].(float64)
		return int64(readCount), nil
	}

	return 0, nil
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

func ListByCategory(categoryId primitive.ObjectID, skip int64, limit int64) ([]primitive.M, error) {
	coll := getStoryCollection()
	projection := bson.D{{"categories", 0}}
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetProjection(projection)
	filter := bson.M{"categories": categoryId}
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
