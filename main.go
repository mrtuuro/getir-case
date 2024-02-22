package main

import (
	"getir-arac/cmd"
	"getir-arac/config"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	env()
	cfg := config.NewConfig(os.Getenv("MONGODB_URI"))
	app := cmd.NewApp(cfg)

	// Handlers
	http.HandleFunc("/api/records", app.HandleGetRecords)    // POST Method
	http.HandleFunc("/api/in-memory", app.HandleGetInMemory) // GET Method
	http.HandleFunc("/api/memory", app.HandleInsertInMemory) // POST Method

	log.Print("listening on port:", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

// env() is used to load environment variables like database connection uri or server port.
func env() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}
