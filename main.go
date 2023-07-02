package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aaradhyakul/goserver/internal/database"
	chi "github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
 
type apiConfig struct {
	DB *database.Queries
}

func main(){
	godotenv.Load()
	portString:=os.Getenv("PORT")
	fmt.Println("hello")
	if portString==""{
		log.Fatal("PORT not found in env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL==""{
		log.Fatal("dbURL not found in env")
	}
	conn,err := sql.Open("postgres", dbURL)

	if err != nil{
		log.Fatal("Connection Failed!")
	}	
	apiCfg := apiConfig{
		DB:database.New(conn),
	}
	

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET","POST","PUT","DELETE","PATCH"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,


	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err",handlerErr)
	v1Router.Post("/users",apiCfg.handlerCreateUser)
	v1Router.Get("/users",apiCfg.handlerGetUser)

	router.Mount("/v1",v1Router)


	srv := &http.Server{
		Handler:router,
		Addr:":" + portString,
	}

	log.Printf("Server starting on port:%v",portString)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
