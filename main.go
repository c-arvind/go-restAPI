package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	mongoClient := ConnectDB()

	userCollection := mongoClient.Database("Insta").Collection("Users")
	postCollection := mongoClient.Database("Insta").Collection("Posts")

	userHandler := NewUserHandler(userCollection)
	postHandler := NewPostHandler(postCollection)

	http.Handle("/users/", userHandler)
	http.Handle("/posts/", postHandler)

	fmt.Println("server started on 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}
