package mongo

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	UserName string `bson:"username"`
}

const (
	uri = "mongodb://10.2.130.182:20000/test"
)

func initMongo() (context.Context, context.CancelFunc, *Mongo) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	mongo := &Mongo{URI: uri}

	mongo.ConfigWillLoad(ctx)
	mongo.ConfigDidLoad(ctx)

	if err := mongo.Serve(ctx); err != nil {
		panic(err)
	}

	return ctx, cancel, mongo
}

func TestConnect(t *testing.T) {
	ctx, cancel, m := initMongo()
	defer m.Shutdown(ctx)
	defer cancel()
}

func TestListDatabaseNames(t *testing.T) {
	ctx, cancel, m := initMongo()
	defer m.Shutdown(ctx)
	defer cancel()

	names, err := m.Client().ListDatabaseNames(ctx, bson.D{})
	t.Log(names)
	t.Log(err)
}

func TestInsertAndFind(t *testing.T) {
	ctx, cancel, m := initMongo()
	defer m.Shutdown(ctx)
	defer cancel()

	db := m.Client().Database("test")
	coll := db.Collection("user")

	insRet, err := coll.InsertOne(ctx, &User{UserName: "boxgo"})
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(insRet)
	}

	findRet := coll.FindOne(ctx, bson.M{"_id": insRet.InsertedID})
	if findRet.Err() != nil {
		t.Fatal(findRet.Err())
	} else {

		m := User{}
		if err := findRet.Decode(&m); err != nil {
			t.Fatal(err)
		}

		if m.UserName != "boxgo" {
			t.Fatal("username should be 'boxgo'")
		}
	}

	coll.Drop(ctx)
}
