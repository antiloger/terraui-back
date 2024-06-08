package main

import (
	"fmt"
	"log"

	"github.com/Terracode-Dev/terraui-back/api"
	"github.com/Terracode-Dev/terraui-back/database"
	"github.com/Terracode-Dev/terraui-back/util"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DBerr := database.InitDynamoDBClient()
	if DBerr != nil {
		fmt.Println("Error initializing DynamoDB client:", DBerr)
	}

	//TODO: key/token generation for testing Auth -----------------------
	// key, kerr := util.GenerateSecretKey()
	// if kerr != nil {
	// 	fmt.Println("Error generating key:", kerr)
	// }
	//Generetade key from above : N9dnx3hLakwCvns5hY0aEjihuBqtALpBDahXyRRMiS4= use it in getToken method baby...

	TK, terr := util.GetToken("1", "1", "N9dnx3hLakwCvns5hY0aEjihuBqtALpBDahXyRRMiS4=")
	if terr != nil {
		fmt.Println("Error generating token:", terr)
	}
	fmt.Println("Token:", TK)
	//-------------------------------------------------------------------

	server := api.NewServer(":8080")
	server.Run()

}
