package geo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/shvedko/geo-service/internal/repository"
)

type Service struct {
	*repository.Queries
}

func New(db repository.DBTX) (*Service, error) {
	return &Service{
		Queries: repository.New(db),
	}, nil
}

func (s *Service) Post(w http.ResponseWriter, r *http.Request) {
	var point repository.AddPointParams

	err := json.NewDecoder(r.Body).Decode(&point)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !point.IsValid() {
		http.Error(w, "empty name or invalid lat/lon", http.StatusBadRequest)
		return
	}

	id, err := s.AddPoint(r.Context(), point)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", path.Join(r.RequestURI, fmt.Sprint(id)))
	w.WriteHeader(http.StatusCreated)
}

func (s *Service) Box(w http.ResponseWriter, r *http.Request) {
	var box repository.GetPointsFromBoxParams

	err := Decode(r, &box)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !box.IsValid() {
		http.Error(w, "min_lon/min_lat must be valid and less than max_lon/max_lat", http.StatusBadRequest)
		return
	}

	points, err := s.GetPointsFromBox(r.Context(), box)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(points)
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	var point repository.GetPointRow

	err := Decode(r, &point)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	point, err = s.GetPoint(r.Context(), point.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(point)
}
