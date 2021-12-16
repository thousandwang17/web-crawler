package Connect

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mu sync.Mutex
var mongoClient *mongo.Client

func Mongo() {

	fmt.Println(`初始化連接 mongo !!`)
	for {
		if connect() {
			break
		}
		fmt.Println(`初始化連接失敗 5秒後重新連接 !!`)
		time.Sleep(5 * time.Second)
	}

	// 監聽並重連
	go func() {
		for {
			if !pin(mongoClient) {
				fmt.Println(`重新連接 mongo !!`)
				connect()
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func connect() bool {

	credential := options.Credential{
		Username: "root",
		Password: "root",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))

	fmt.Printf(`%v`, err)

	if err == nil && pin(client) {
		mu.Lock()
		mongoClient = client
		mu.Unlock()
		fmt.Println(` mongo 連接成功!!`)
		return true
	}

	return false

}

func GetMongo() *mongo.Client {
	mu.Lock()
	client := mongoClient
	mu.Unlock()

	return client
}

func pin(client *mongo.Client) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	err := client.Ping(ctx, readpref.Primary())

	if err == nil {
		return true
	}
	return false
}
