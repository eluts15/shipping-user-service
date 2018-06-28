package main

import (
	"log"

	pb "github.com/eluts15/shipping-user-service/auth"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/mdns"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {

	// Create a database connection and handle closing it before exiting.
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	db.AutoMigrate(&pb.User{})

	repo := &UserRepository{db}

	tokenService := &TokenService{repo}

	srv := k8s.NewService(
		micro.Name("shipping.auth"),
	)

	srv.Init()
	pb.RegisterAuthHandler(srv.Server(), &service{repo, tokenService})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
