package http_handler

import (
	"context"
	"encoding/json"
	"github.com/germanov-v/go-rotator/internal/model"
	"github.com/germanov-v/go-rotator/internal/repository/postgres"
	"github.com/germanov-v/go-rotator/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

type jsonResponse struct {
	BannerID model.BannerId `json:"banner_id,omitempty"`
	Message  string         `json:"message,omitempty"`
}

type addBannerRequest struct {
	BannerID model.BannerId `json:"banner_id"`
}

type clickRequest struct {
	BannerID model.BannerId `json:"banner_id"`
	GroupID  model.GroupId  `json:"group_id"`
}

// addBannerHandler(svc *service.RotationService) http.HandlerFunc {
func AddBannerHandler(repo *postgres.PostgresRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slot := model.SlotId(vars["slot"])
		var req addBannerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//err := json.NewEncoder(w).Encode(jsonResponse{Message: "invalid request body"})
			//if err != nil {
			//	return
			//}
			_ = json.NewEncoder(w).Encode(jsonResponse{Message: "invalid request body"})

			return
		}
		if err := repo.AddBanner(context.Background(), slot, req.BannerID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(jsonResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(jsonResponse{Message: "banner added"})
	}
}

// func removeBannerHandler(svc *service.RotationService) http.HandlerFunc {
func RemoveBannerHandler(repo *postgres.PostgresRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slot := model.SlotId(vars["slot"])
		banner := model.BannerId(vars["banner"])
		if err := repo.RemoveBanner(context.Background(), slot, banner); err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(jsonResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jsonResponse{Message: "banner removed"})
	}
}

func RotateBannerHandler(svc *service.RotationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slot := model.SlotId(vars["slot"])
		group := model.GroupId(r.URL.Query().Get("group"))
		if group == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(jsonResponse{Message: "missing group query parameter"})
			return
		}
		banner, err := svc.Rotate(context.Background(), slot, group)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(jsonResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jsonResponse{BannerID: banner})
	}
}

// func recordClickHandler(svc *service.RotationService) http.HandlerFunc {
func RecordClickHandler(repo *postgres.PostgresRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		slot := model.SlotId(vars["slot"])
		var req clickRequest
		//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//	w.WriteHeader(500)
		//	json.NewEncoder(w).Encode(jsonResponse{Message: "invalid request body"})
		//	return
		//}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(jsonResponse{Message: "invalid request body"})
			return
		}
		if err := repo.IncrementClick(context.Background(), slot, req.BannerID, req.GroupID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(jsonResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jsonResponse{Message: "click recorded"})
	}
}
