package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/mongodb/mongo-go-driver/bson"
// 	"gitlab.com/akita/util/tracing"
// )

// func httpTask(w http.ResponseWriter, r *http.Request) {
// 	id := r.FormValue("id")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	filter := bson.D{bson.E{Key: "id", Value: id}}

// 	cursor, err := collection.Find(ctx, filter)
// 	dieOnErr(err)
// 	defer cursor.Close(ctx)

// 	tasks := make([]tracing.Task, 0)

// 	for cursor.Next(ctx) {
// 		task := tracing.Task{}
// 		err = cursor.Decode(&task)
// 		dieOnErr(err)

// 		tasks = append(tasks, task)
// 	}

// 	rsp, err := json.Marshal(tasks)
// 	dieOnErr(err)

// 	fmt.Println(rsp)
// 	_, err = w.Write([]byte(rsp))
// 	dieOnErr(err)
// }
