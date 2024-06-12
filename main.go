package main

import (
	"fmt"
	"log"

	"github.com/Terracode-Dev/terraui-back/api"
	"github.com/Terracode-Dev/terraui-back/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// --- INIT DynamoDB client ---
	DB, err := database.StartDB()
	if err != nil {
		fmt.Println("Error initializing DynamoDB client:", err)
	}

	// TODO: key/token generation for testing Auth -----------------------
	// key, kerr := util.GenerateSecretKey()
	// if kerr != nil {
	// 	fmt.Println("Error generating key:", kerr)
	// }
	// Generetade key from above and store in .env adn use it in getToken method baby...

	// TK, terr := util.GetToken("1", "1", os.Getenv("JKEY"))
	// if terr != nil {
	// 	fmt.Println("Error generating token:", terr)
	// }
	// fmt.Println("Token:", TK)
	//-------------------------------------------------------------------

	// --- INIT NewServer ---
	server := api.NewServer(":8088", DB)
	server.Run()
}
