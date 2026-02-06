package geo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"

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
	var p repository.AddPointParams

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !p.IsValid() {
		http.Error(w, "Invalid coordinates or empty name. Lat: [-90, 90], Lon: [-180, 180]", http.StatusBadRequest)
		return
	}

	id, err := s.AddPoint(r.Context(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", path.Join(r.RequestURI, fmt.Sprint(id)))
	w.WriteHeader(http.StatusCreated)
}

func Float64(r *http.Request, key string, min, max float64) (float64, error) {
	value := chi.URLParam(r, key)

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("param %s is invalid: %w", key, err)
	}

	if val < min || val > max {
		return 0, fmt.Errorf("param %s out of range [%.f, %.f]", key, min, max)
	}

	return val, nil
}

func (s *Service) Box(w http.ResponseWriter, r *http.Request) {
	x1, err1 := Float64(r, "left", -180, 180)
	y1, err2 := Float64(r, "top", -90, 90)
	x2, err3 := Float64(r, "right", -180, 180)
	y2, err4 := Float64(r, "bottom", -90, 90)

	for _, err := range []error{err1, err2, err3, err4} {
		if err != nil {
			http.Error(w, fmt.Sprintln("Invalid coordinate range:", err.Error()), http.StatusBadRequest)
			return
		}
	}

	if x1 > x2 || y1 > y2 {
		http.Error(w, "Invalid coordinate range: left/top must be less than right/bottom", http.StatusBadRequest)
		return
	}

	points, err := s.GetPointsFromBox(r.Context(), repository.GetPointsFromBoxParams{
		Lon1: x1,
		Lat1: y1,
		Lon2: x2,
		Lat2: y2,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(points)
}

func Int32(r *http.Request, key string) (int32, error) {
	value := chi.URLParam(r, key)

	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("param %s is invalid: %w", key, err)
	}

	return int32(val), nil
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	id, err := Int32(r, "id")
	if err != nil {
		http.Error(w, fmt.Sprintln("Invalid id: ", err.Error()), http.StatusBadRequest)
		return
	}

	point, err := s.GetPoint(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(point)
}
