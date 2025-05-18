package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Password  string `json:"password"`
	}

	updateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
		Password  *string `json:"password"`
	}
	ErrRes struct {
		Error string `json:"error"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"Invalid request format"})
			return
		}

		if req.FirstName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"first name is required"})
			return
		}
		if len(req.FirstName) > 50 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"first name must be at most 50 characters"})
			return
		}

		if req.LastName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"last name is required"})
			return
		}
		if len(req.LastName) > 50 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"last name must be at most 50 characters"})
			return
		}

		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"email is required"})
			return
		}
		if len(req.Email) > 50 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"email must be at most 50 characters"})
			return
		}

		if req.Phone == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"phone is required"})
			return
		}
		if len(req.Phone) > 11 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"phone must be at most 11 characters"})
			return
		}

		if req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"password is required"})
			return
		}
		if len(req.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrRes{"password must be at least 8 characters"})
			return
		}

		user, err := s.Create(req.FirstName, req.LastName, req.Email, req.Phone, req.Password)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrRes{err.Error()})
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.GetAll()
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrRes{err.Error()})
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrRes{"user does not exist"})
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req updateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrRes{"Invalid request format"})
			return
		}
		if req.FirstName != nil && *req.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrRes{"first name is required"})
			return
		}

		if req.LastName != nil && *req.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrRes{"last name is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Update(id, req.FirstName, req.LastName, req.Email, req.Phone, req.Password); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrRes{"user does not exist"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"data": "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		if err := s.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrRes{"user does not exist"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"data": "success"})
	}
}
