package collections

import (
	"context"
	"errors"
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
	ReadStatus                      UserStoryReadStatus  `bson:"read_status"`
	InterestedCategories            []primitive.ObjectID `bson:"interested_categories"`
	FollowingAuthors                []string             `bson:"following_authors"`
	CategoriesBasedOnReadingHistory []primitive.ObjectID `bson:"categories_based_on_reading_history"`
}

type UserStoryReadStatus struct {
	StoryId   primitive.ObjectID `bson:"story_id"`
	ChapterId primitive.ObjectID `bson:"chapter_id"`
}

func NewStoryFeed() *StoryFeed {
	return &StoryFeed{}
}

func (s *StoryFeed) LoadByUser(userUuid string) error {
	coll := getStoryFeedCollection()
	err := coll.FindOne(context.TODO(), bson.M{"user_uuid": userUuid}).Decode(&s)

	if err != nil && err == mongo.ErrNoDocuments {
		return errors.New("requested story not found")
	} else if err != nil {
		return err
	}

	return nil
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

func (s *StoryFeed) AddCategoriesBasedOnReadingHistoryToDocument(categories []primitive.ObjectID) error {
	coll := getStoryFeedCollection()
	var objIds bson.A
	for _, val := range categories {
		objIds = append(objIds, val)
	}
	data := bson.D{
		{"categories_based_on_reading_history", bson.D{
			{"$each", objIds},
		}},
	}
	filter := bson.D{
		{"_id", s.Id},
	}
	update := bson.D{{"$addToSet", data}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoryFeed) AddFollowingAuthorToDocument() error {
	coll := getStoryFeedCollection()
	data := bson.D{
		{"following_authors", s.FollowingAuthors},
	}
	filter := bson.D{
		{"_id", s.Id},
	}
	update := bson.D{{"$set", data}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoryFeed) RemoveFollowingAuthorFromDocument(authorUuid string) error {
	var newAuthors []string
	for _, v := range s.FollowingAuthors {
		if v == authorUuid {
			continue
		}
		newAuthors = append(newAuthors, v)

	}
	s.FollowingAuthors = newAuthors

	coll := getStoryFeedCollection()
	data := bson.D{
		{"following_authors", s.FollowingAuthors},
	}
	filter := bson.D{
		{"_id", s.Id},
	}
	update := bson.D{{"$set", data}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
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

func (s *StoryFeed) FetchRecommendedStories() ([]bson.M, error) {
	m := make(map[primitive.ObjectID]bool)
	var categories []primitive.ObjectID
	for _, v := range s.InterestedCategories {
		m[v] = true
		categories = append(categories, v)
	}
	for _, v := range s.CategoriesBasedOnReadingHistory {
		if m[v] {
			continue
		}
		m[v] = true
		categories = append(categories, v)
	}
	coll := getCategoryCollection()
	var objIds bson.A
	for _, val := range categories {
		objIds = append(objIds, val)
	}
	matchStage := bson.D{
		{"$match", bson.D{{"_id", bson.D{{"$in", objIds}}}}},
	}
	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "stories"},
			{"as", "stories"},
			{"pipeline", bson.A{
				bson.D{{"$sort", bson.D{{"_id", -1}}}},
				bson.D{{"$limit", 1}},
			}},
		}},
	}
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, lookupStage})
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

func getStoryFeedCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(StoryFeedCollection)

	return coll
}
