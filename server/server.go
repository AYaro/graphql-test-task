package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	graphql_test_task "github.com/AYaro/graphql-test-task"
	"github.com/gorilla/websocket"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql_test_task.NewExecutableSchema(graphql_test_task.Config{Resolvers: &graphql_test_task.Resolver{}}), handler.WebsocketKeepAliveDuration(19*time.Second),
		handler.WebsocketUpgrader(websocket.Upgrader{
			CheckOrigin: func(request *http.Request) bool {
				return true
			},
			HandshakeTimeout: 5 * time.Second,
		})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Println(http.ListenAndServe(":"+port, nil))
}
