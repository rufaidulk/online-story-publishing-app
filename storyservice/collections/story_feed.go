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

type StoryFeed struct {
	Id                              primitive.ObjectID   `bson:"_id"`
	UserUuid                        string               `bson:"user_uuid"`
	InterestedCategories            []primitive.ObjectID `bson:"interested_categories"`
	FollowingAuthors                []string             `bson:"following_authors"`
	CategoriesBasedOnReadingHistory []primitive.ObjectID `bson:"categories_based_on_reading_history"`
}

func NewStoryFeed() *StoryFeed {
	return &StoryFeed{}
}

func (s *StoryFeed) UpsertDocument() error {
	coll := getStoryFeedCollection()
	data := bson.D{
		{"user_uuid", s.UserUuid},
		{"interested_categories", s.InterestedCategories},
	}
	filter := bson.D{
		{"user_uuid", s.UserUuid},
	}
	update := bson.D{{"$set", data}}
	opts := options.Update().SetUpsert(true)
	_, err := coll.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoryFeed) SetCategories(categories []string) error {
	s.InterestedCategories = nil
	for _, val := range categories {
		objId, err := primitive.ObjectIDFromHex(val)
		if err != nil {
			return err
		}
		s.InterestedCategories = append(s.InterestedCategories, objId)
	}

	return nil
}

func getStoryFeedCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(StoryFeedCollection)

	return coll
}
