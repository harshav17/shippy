package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/micro/go-micro/metadata"

	"context"

	pb "github.com/harshav17/shippy/consignment-service/proto/consignment"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	microclient "github.com/micro/go-micro/client"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	// Create new greeter client
	client := pb.NewShippingServiceClient("consignment", microclient.DefaultClient)

	//Define out flags
	service := micro.NewService(
		micro.Name("consignment.cli"),
		micro.Flags(
			cli.StringFlag{
				Name:  "filename",
				Usage: "Name of the json",
			},
			cli.StringFlag{
				Name:  "token",
				Usage: "token genenrated from user service",
			},
		),
	)

	//Start as service
	service.Init(
		micro.Action(func(c *cli.Context) {
			filename := c.String("filename")
			token := c.String("token")

			consignment, err := parseFile(filename)

			if err != nil {
				log.Fatalf("Could not parse file: %v", err)
			}

			// Create a new context which contains our given token.
			// This same context will be passed into both the calls we make
			// to our consignment-service.
			ctx := metadata.NewContext(context.Background(), map[string]string{
				"token": token,
			})

			// First call using our tokenised context
			r, err := client.CreateConsignment(ctx, consignment)
			if err != nil {
				log.Fatalf("Could not create: %v", err)
			}
			log.Printf("Created: %t", r.Created)

			// Second call
			getAll, err := client.GetConsignments(ctx, &pb.GetRequest{})
			if err != nil {
				log.Fatalf("Could not list consignments: %v", err)
			}
			for _, v := range getAll.Consignments {
				log.Println(v)
			}

			os.Exit(0)
		}),
	)

	// Run the server
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
