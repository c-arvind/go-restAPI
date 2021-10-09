package main

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userCollection *mongo.Collection
}

func NewUserHandler(usercol *mongo.Collection) *UserHandler {
	return &UserHandler{
		userCollection: usercol,
	}
}

var (
	//listUserRe = regexp.MustCompile(`^\/users[\/]*$`)
	getUserRe    = regexp.MustCompile(`^\/users[\/][0-9a-fA-F]{24}$`)
	createUserRe = regexp.MustCompile(`^\/users[\/]*$`)
)

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	//case r.Method == http.MethodGet && listUserRe.MatchString(r.URL.Path):
	//    h.List(w, r)
	//    return
	case r.Method == http.MethodGet && getUserRe.MatchString(r.URL.Path):
		h.GetUser(w, r)
		return
	case r.Method == http.MethodPost && createUserRe.MatchString(r.URL.Path):
		h.CreateUser(w, r)
		return
	default:
		notFound(w, r)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	check := Readjson(w, r, user)
	if !check {
		return
	}
	user.Password = Hash(user.Password)
	userResult, err := h.userCollection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprintf("Successfully created user with id: %v", userResult.InsertedID)))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	matches := getUserRe.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		notFound(w, r)
		return
	}

	user := &User{}
	userResult := h.userCollection.FindOne(context.Background(), bson.D{{"_id", matches[1]}})
	err := userResult.Decode(user)
	if err != nil {
		w.Write([]byte("unable to get data"))
	} else {
		Writejson(w, r, user)
	}

}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}
