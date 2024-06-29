package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Server struct {
	//ListenAddr string
	ChoreStorage ChoreStorage
}

func (s *Server) ListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	items, err := s.ChoreStorage.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) TodoListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	items, err := s.ChoreStorage.ListTodo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) ItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := s.ChoreStorage.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusAccepted)
}

func (s *Server) CreateHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	var item Chore
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := s.ChoreStorage.Create(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getId(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var item Chore
	if err = json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.ChoreStorage.Update(id, item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	//w.Write([]byte("{}"))
}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.ChoreStorage.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	//w.Write([]byte("{}"))
}

func (s *Server) CompleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var execution Execution
	if err := json.NewDecoder(r.Body).Decode(&execution); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = s.ChoreStorage.Complete(id, &execution)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	//w.Write([]byte("{}"))
}

func (s *Server) GetHandler() http.Handler {
	srv := http.NewServeMux()
	srv.HandleFunc("GET /api/chore", s.ListHandler)
	srv.HandleFunc("GET /api/chores", s.TodoListHandler)
	srv.HandleFunc("GET /api/chore/{id}", s.ItemHandler)
	srv.HandleFunc("POST /api/chore", s.CreateHandler)
	srv.HandleFunc("PUT /api/chore/{id}", s.UpdateHandler)
	srv.HandleFunc("PATCH /api/chore/{id}", s.CompleteHandler)
	srv.HandleFunc("DELETE /api/chore/{id}", s.DeleteHandler)
	return srv
}
