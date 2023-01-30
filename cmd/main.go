package main

import (
	"fmt"
	"log"
	"net"

	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/config"
	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/db"
	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/pb"
	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/services"
	"github.com/alejandrogh97/rent-hive-auth-svc/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "rent-hive-auth-svc",
		ExpirationHours: 1,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed at listen:", err)
	}

	fmt.Println("Auth svc listening on port", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
