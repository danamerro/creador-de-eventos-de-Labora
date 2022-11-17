package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

// -------Structs------
// bd ficticia
type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}
type allEvents []event

var events = allEvents{
	{ID: "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

// ------Eventos------
//crear evento

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

//traer un evento

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

//traer todos los eventos

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

//editar evento

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

//borrar evento

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1])
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

//-----------Main-----------------

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenido al Blog!")
}
func main() {
	router := mux.NewRouter().StrictSlash(true)
	//ruta principal
	router.HandleFunc("/", Index)
	//traer todos los eventos
	router.HandleFunc("/event/getAll", getAllEvents).Methods("GET")
	//crea eventos
	router.HandleFunc("/event/create", createEvent).Methods("POST")
	//muestra eventos por id
	router.HandleFunc("/event/get/{id}", getOneEvent).Methods("GET")
	//modifica eventos por id
	router.HandleFunc("/event/update/{id}", updateEvent).Methods("PATCH")
	//elimina eventos por id
	router.HandleFunc("/event/delete/{id}", deleteEvent).Methods("DELETE")
	server := http.ListenAndServe(":8080", router)
	log.Fatal(server)
}
