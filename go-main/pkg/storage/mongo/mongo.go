package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName         = "goNews"
	collectionName = "posts"
)

type MongoDB struct {
	c *mongo.Client
}

func New(connstr string) (*MongoDB, error) {
	mongoOpts := options.Client().ApplyURI(connstr)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	db := MongoDB{
		c: client,
	}
	return &db, nil
}

func (db *MongoDB) Posts() ([]storage.Post, error) {
	var data []storage.Post
	collection := db.c.Database(dbName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}

	return data, cur.Err()
}

func (db *MongoDB) AddPost(p storage.Post) error {
	collection := db.c.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDB) UpdatePost(p storage.Post) error {
	filter := bson.D{{Key: "_id", Value: p.ID}}
	update := bson.D{primitive.E{Key: "$inc", Value: bson.D{ //or $set
		{Key: "title", Value: p.Title},
		{Key: "content", Value: p.Content},
	},
	}}

	_, err := db.c.Database(dbName).Collection(collectionName).UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *MongoDB) DeletePost(p storage.Post) error {
	filter := bson.D{primitive.E{Key: "_id", Value: p.ID}}
	_, err := db.c.Database(dbName).Collection(collectionName).DeleteOne(
		context.Background(),
		filter,
	)
	if err != nil {
		panic(err)
	}
	return nil
}
