package main

import (
	"net/http"

	"./restAPI"
)

func main() {
	restAPI.InitMongoConnection()
	mux := http.NewServeMux()
	mux.HandleFunc("/", restAPI.HomePageHandler)
	mux.HandleFunc("/users", restAPI.CreateUsersHandler)
	mux.HandleFunc("/users/", restAPI.GetUsersHandler)
	mux.HandleFunc("/posts", restAPI.CreatePostsHandler)
	mux.HandleFunc("/posts/", restAPI.GetPostsHandler)
	mux.HandleFunc("/posts/users/", restAPI.GetPostsByUserHandler)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
