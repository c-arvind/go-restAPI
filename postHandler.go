package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
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

var (
	//listUserRe = regexp.MustCompile(`^\/users[\/]*$`)
	getPostRe     = regexp.MustCompile(`^\/posts[\/][0-9a-fA-F]{24}$`)
	createPostRe  = regexp.MustCompile(`^\/posts[\/]*$`)
	getPostUserRe = regexp.MustCompile(`^\/posts[\/]users[\/][0-9a-fA-F]{24}$`)
)

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && getPostRe.MatchString(r.URL.Path):
		h.GetPost(w, r)
		return
	case r.Method == http.MethodPost && createPostRe.MatchString(r.URL.Path):
		h.CreatePost(w, r)
		return
	case r.Method == http.MethodGet && getPostUserRe.MatchString(r.URL.Path):
		h.GetPostUser(w, r)
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
	err := h.postCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Writejson(w, r, post)
}

func (h *PostHandler) GetPostUser(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/posts/users/"):]
	fmt.Println(id)

	postCursor, err := h.postCollection.Find(context.TODO(), bson.D{{"userId", id}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	posts := Post{}
	err = postCursor.All(context.TODO(), posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Writejson(w, r, posts)
}
