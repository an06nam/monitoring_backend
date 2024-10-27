package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/dieselMonitoring/internal/db"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func SetupCRUDRoutes(dbConn *db.DB) {
	r := mux.NewRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	r.HandleFunc("/sensor/{uid}", GetDataByUID(dbConn)).Methods("GET")

	log.Printf("Starting API server on %s...", os.Getenv("CRUD_SERVER_PORT"))
	log.Fatal(http.ListenAndServe(os.Getenv("CRUD_SERVER_PORT"), r))
}

func GetDataByUID(dbConn *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid := vars["uid"]

		data, err := dbConn.GetDataByUID(uid)
		if err != nil {
			http.Error(w, "Data not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
