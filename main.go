package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	chi "github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)
 

func main(){
	godotenv.Load()
	portString:=os.Getenv("PORT")
	fmt.Println("hello")
	if portString==""{
		log.Fatal("PORT not found")
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

	router.Mount("/v1",v1Router)


	srv := &http.Server{
		Handler:router,
		Addr:":" + portString,
	}

	log.Printf("Server starting on port:%v",portString)
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
