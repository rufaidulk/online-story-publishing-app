package collections

import (
	"context"
	"errors"
	"fmt"
	"log"
	"storyservice/adapters"
	"storyservice/helper"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "categories"

type Category struct {
	Id      primitive.ObjectID   `bson:"_id"`
	Slug    string               `bson:"slug"`
	Name    string               `bson:"name"`
	Stories []primitive.ObjectID `bson:"stories"`
}

func NewCategory() *Category {
	return &Category{}
}

func (c *Category) LoadById(id primitive.ObjectID) error {
	coll := getCollection()
	err := coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&c)

	if err != nil && err == mongo.ErrNoDocuments {
		return errors.New("requested category not found")
	} else if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c *Category) CreateDocument(name string) error {
	coll := getCollection()
	c.Id = primitive.NewObjectID()
	c.Name = name
	result, err := coll.InsertOne(context.TODO(), &c)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	fmt.Println(result)
	if err := c.updateSlug(); err != nil {
		return err
	}

	return nil
}

func (c *Category) updateSlug() error {
	c.Slug = strings.ReplaceAll(strings.ToLower(c.Name), " ", "-")
	c.Slug = fmt.Sprintf("%s-%s", c.Slug, c.Id.Hex())

	coll := getCollection()
	filter := bson.D{{"_id", c.Id}}
	update := bson.D{{"$set", bson.D{{"slug", c.Slug}}}}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func CreateCategoryCollectionIndexes() {
	coll := getCollection()
	indexView := coll.Indexes()
	model := mongo.IndexModel{
		Keys:    bson.D{{"name", 1}},
		Options: options.Index().SetUnique(true),
	}
	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	names, err := indexView.CreateOne(context.TODO(), model, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created indexes %v\n", names)
}

func getCollection() *mongo.Collection {
	dbClient := adapters.GetDbClient()
	coll := dbClient.Database(helper.GetEnv("MONGO_DB")).Collection(collectionName)

	return coll
}
