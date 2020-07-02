package cloudlocker

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"io/ioutil"
	"net/http"
)

func newRouter(server *LockerServer) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/set", handleSet(server)).Methods("POST")
	router.HandleFunc("/get", handleGet(server)).Methods("POST")
	router.HandleFunc("/delete", handleDelete(server)).Methods("POST")
	return router
}

func handleSet(server *LockerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var e entry
		err = json.Unmarshal(body, &e)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = server.db.Put(e.K, e.V, &opt.WriteOptions{Sync: false})
	}
}

func handleGet(server *LockerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		v, err := server.db.Get(body, nil)
		if err != nil {
			if err != leveldb.ErrNotFound {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		_, _ = w.Write(v)
	}
}

func handleDelete(server *LockerServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//log.Println("r.Body", string(body))
		_ = server.db.Delete(body, nil)
	}
}
