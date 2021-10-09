package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostHandler struct {
	postCollection *mongo.Collection
}

func NewPostHandler(postcol *mongo.Collection) *PostHandler {
	return &PostHandler{
		postCollection: postcol,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && getUserRe.MatchString(r.URL.Path):
		h.GetPost(w, r)
		return
	case r.Method == http.MethodPost && createUserRe.MatchString(r.URL.Path):
		h.CreatePost(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	ok := Readjson(w, r, post)
	if !ok {
		return
	}

	post.Timestamp = time.Now().Format(time.RFC850)

	_, err := h.postCollection.InsertOne(context.TODO(), post)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("successfully created post"))
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/posts/"):]
	fmt.Println(id)

	post := Post{}
	err := h.postCollection.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Writejson(w, r, post)
}
