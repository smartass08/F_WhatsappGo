package utils

import (
	"context"
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DB struct {
	client *mongo.Client
}

type session struct {
	ClientId    string `json:"clientId"`
	ClientToken string `json:"clientToken"`
	EncKey      []byte `json:"encKey"`
	Mackey      []byte `json:"macKey"`
	ServerToken string `json:"serverToken"`
	Wid         string `json:"wid"`
}

func (C *DB) Access(url string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error while connecting DB : %v", err)
		return
	}
	C.client = client
}

func (C DB) GetKey() (bool, whatsapp.Session) {
	Ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	var value whatsapp.Session
	defer cancel()
	collection := C.client.Database(GetDbName()).Collection(GetDbCollection())
	var sus *session
	err := collection.FindOne(Ctx, bson.M{"key": "access"}).Decode(&sus)
	if err != nil {
		fmt.Println(err.Error())
		if err == mongo.ErrNoDocuments {
			return false, value
		}
	}
	value = whatsapp.Session{
		ClientId:    sus.ClientId,
		ClientToken: sus.ClientToken,
		ServerToken: sus.ServerToken,
		EncKey:      sus.EncKey,
		MacKey:      sus.Mackey,
		Wid:         sus.Wid,
	}
	return true, value
}

func (C DB) Addkey(kek whatsapp.Session) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := C.client.Database(GetDbName()).Collection(GetDbCollection())
	_, err := collection.InsertOne(ctx, bson.M{
		"key":         "access",
		"clientId":    kek.ClientId,
		"clientToken": kek.ClientToken,
		"encKey":      kek.EncKey,
		"macKey":      kek.MacKey,
		"serverToken": kek.ServerToken,
		"wid":         kek.Wid})
	if err != nil {
		return false
	}
	return true
}

func (C DB) DelKeys() {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := C.client.Database(GetDbName()).Collection(GetDbCollection())
	res, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
}
