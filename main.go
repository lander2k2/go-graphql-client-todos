package main

import (
	"flag"
	"log"

	graphql "github.com/hasura/go-graphql-client"
)

func startSubscription(serverAddr string) error {
	client := graphql.NewSubscriptionClient(serverAddr).
		WithLog(log.Println).
		OnError(func(sc *graphql.SubscriptionClient, err error) error {
			log.Print("error", err)
			return err
		})

	defer client.Close()

	var subscription struct {
		Todos struct {
			Text string `graphql:"text"`
			Done bool   `graphql:"done"`
			User struct {
				Name string `graphql:"name"`
			} `graphql:"user"`
		} `graphql:"todoNotifs"`
	}

	_, err := client.Subscribe(subscription, nil, func(data []byte, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		if data == nil {
			return nil
		}
		log.Println("notification received")
		log.Println(string(data))
		return nil
	})

	if err != nil {
		return err
	}

	return client.Run()
}

func main() {
	serverAddr := flag.String("server", "ws://localhost:8080/query", "GraphQL server address")
	flag.Parse()

	if err := startSubscription(*serverAddr); err != nil {
		panic(err)
	}
}
