package app

import (
	"fmt"
	"github.com/SemmiDev/lets-tests/controllers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func routes(router *chi.Mux) {
	router.Options("/*", func(w http.ResponseWriter, r *http.Request) {})
	api := router.Route("/api/v1", func(router chi.Router) {})

	api.Route("/chats", func(r chi.Router) {
		r.Post("/", controllers.CreateChat)
		r.Get("/", controllers.GetAllChats)
		r.Get("/{chat_id}", controllers.GetChat)
		r.Put("/{chat_id}", controllers.UpdateChat)
		r.Delete("/{chat_id}", controllers.DeleteChat)
	})

	fmt.Println()
	registeredEndpointLog("/", "POST", "CreateChat")
	registeredEndpointLog("/", "GET", "GetAllChat")
	registeredEndpointLog("/{chat_id}", "GET", "GetChat")
	registeredEndpointLog("/{chat_id}", "PUT", "UpdateChat")
	registeredEndpointLog("/{chat_id}", "DELETE", "DeleteChat")

}

func registeredEndpointLog(path, act, handler string) {
	log.Println("Registered Endpoint :: " + act + " ::  " + " localhost:3333/api/v1/chats" + path + " :: " + "Handler -> " + handler)
}
