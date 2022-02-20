package collections

import (
	"context"
	"fmt"
	"storyservice/adapters"
	"storyservice/helper"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Story struct {
	Id               primitive.ObjectID   `bson:"_id"`
	UserUuid         string               `bson:"user_uuid"`
	Slug             string               `bson:"slug"`
	Title            string               `bson:"title"`
	PromotionalTitle string               `bson:"promotional_title"`
	PromotionalImage string               `bson:"promotional_image"`
	LanguageCode     string               `bson:"language_code"`
	Categories       []primitive.ObjectID `bson:"categories"`
	Chapters         map[int]ChapterInfo  `bson:"chapters"`
	IsPremium        bool                 `bson:"is_premium"`
	IsCompleted      bool                 `bson:"is_completed"`
	Rating           int8                 `bson:"rating"`
	AvgReadCount     int64                `bson:"avg_read_count"`
}

type ChapterInfo struct {
	ChapterId    primitive.ObjectID
	ChapterTitle string
}

func NewStory() *Story {
	return &Story{}
}

func (s *Story) SetCategories(categories []string) error {
	for _, val := range categories {
		objId, err := primitive.ObjectIDFromHex(val)
		if err != nil {
			return err
		}
		s.Categories = append(s.Categories, objId)
	}

	return nil
}

func (s *Story) CreateDocument() error {
	coll := getStoryCollection()
	s.Id = primitive.NewObjectID()
	s.Chapters = make(map[int]ChapterInfo)
	_, err := coll.InsertOne(context.TODO(), &s)
	if err != nil {
		return err
	}

	if err := s.updateSlug(); err != nil {
		return err
	}

	return nil
}

func (s *Story) updateSlug() error {
	s.Slug = strings.ReplaceAll(strings.ToLower(s.Title), " ", "-")
	s.Slug = fmt.Sprintf("%s-%s", s.Slug, s.Id.Hex())

	coll := getStoryCollection()
	filter := bson.D{{"_id", s.Id}}
	update := bson.D{{"$set", bson.D{{"slug", s.Slug}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *Story) AddChapter(chapterId primitive.ObjectID, chapterTitle string) error {
	chapterInfo := ChapterInfo{
		ChapterId:    chapterId,
		ChapterTitle: chapterTitle,
	}
	chapterNo := len(s.Chapters) + 1
	s.Chapters[chapterNo] = chapterInfo
	coll := getStoryCollection()
	filter := bson.D{{"_id", s.Id}}
	update := bson.D{{"$set", bson.D{{"chapters", s.Chapters}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func getStoryCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(StoryCollection)

	return coll
}
