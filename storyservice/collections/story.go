package collections

import (
	"context"
	"errors"
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
	Chapters         []ChapterInfo        `bson:"chapters"`
	IsPremium        bool                 `bson:"is_premium"`
	IsCompleted      bool                 `bson:"is_completed"`
	Rating           float64              `bson:"rating"`
	AvgReadCount     int64                `bson:"avg_read_count"`
}

type ChapterInfo struct {
	ChapterId    primitive.ObjectID `bson:"chapter_id"`
	ChapterTitle string             `bson:"chapter_title"`
}

func NewStory() *Story {
	return &Story{}
}

func (s *Story) LoadById(id string) error {
	coll := getStoryCollection()
	objId, _ := primitive.ObjectIDFromHex(id)
	err := coll.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&s)

	if err != nil && err == mongo.ErrNoDocuments {
		return errors.New("requested story not found")
	} else if err != nil {
		return err
	}

	return nil
}

func (s *Story) SetCategories(categories []string) error {
	s.Categories = nil
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
	_, err := coll.InsertOne(context.TODO(), &s)
	if err != nil {
		return err
	}

	if err := s.updateSlug(); err != nil {
		return err
	}

	return nil
}

func (s *Story) AddChapter(chapterId primitive.ObjectID, chapterTitle string) {
	chapterInfo := ChapterInfo{
		ChapterId:    chapterId,
		ChapterTitle: chapterTitle,
	}
	s.Chapters = append(s.Chapters, chapterInfo)
}

func (s *Story) EditChapter(chapterId primitive.ObjectID, chapterTitle string) {
	for k, v := range s.Chapters {
		if v.ChapterId == chapterId {
			v.ChapterTitle = chapterTitle
			s.Chapters[k] = v
			break
		}
	}
}

func (s *Story) RemoveChapter(chapterId primitive.ObjectID) {
	var newChapters []ChapterInfo
	for _, v := range s.Chapters {
		if v.ChapterId == chapterId {
			continue
		}
		newChapters = append(newChapters, v)

	}

	s.Chapters = newChapters
}

func (s *Story) UpdateDocument() error {
	s.Slug = strings.ReplaceAll(strings.ToLower(s.Title), " ", "-")
	s.Slug = fmt.Sprintf("%s-%s", s.Slug, s.Id.Hex())
	data := bson.D{
		{"slug", s.Slug},
		{"title", s.Title},
		{"language_code", s.LanguageCode},
		{"promotional_title", s.PromotionalTitle},
		{"promotional_image", s.PromotionalImage},
		{"categories", s.Categories},
		{"chapters", s.Chapters},
		{"is_premium", s.IsPremium},
		{"is_completed", s.IsCompleted},
	}
	coll := getStoryCollection()
	filter := bson.D{{"_id", s.Id}}
	update := bson.D{{"$set", data}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *Story) UpdateRating() error {
	avgStoryRating, err := CalculateAvgRatingOfStory(s.Id)
	if err != nil {
		return err
	}
	coll := getStoryCollection()
	data := bson.D{
		{"rating", avgStoryRating},
	}
	filter := bson.D{{"_id", s.Id}}
	update := bson.D{{"$set", data}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *Story) UpdateAvgReadCount() error {
	avgReadCount, err := CalculateAvgReadCountOfStory(s.Id)
	if err != nil {
		return err
	}
	coll := getStoryCollection()
	data := bson.D{
		{"avg_read_count", avgReadCount},
	}
	filter := bson.D{{"_id", s.Id}}
	update := bson.D{{"$set", data}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
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

func getStoryCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(StoryCollection)

	return coll
}
