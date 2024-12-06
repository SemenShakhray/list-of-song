package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"listsongs/internal/models"
	"listsongs/internal/service"

	"go.uber.org/zap"
)

type Handler struct {
	Log     *zap.Logger
	Service service.Servicer
}

func NewHandler(log *zap.Logger, serv service.Servicer) Handler {
	return Handler{
		Log:     log,
		Service: serv,
	}
}
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {

	var song models.Song

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.Log.Debug("Failed to decode request body", zap.Error(err))
		http.Error(w, `{"error":"failed to decode request"}`, http.StatusBadRequest)
		return
	}

	h.Log.Debug("Song data", zap.String("song", song.Song), zap.String("group", song.Group))

	err = h.Service.AddSong(r.Context(), song)
	if err != nil {
		h.Log.Debug("service AddSong", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Song added successfully", zap.String("song", song.Song), zap.String("group", song.Group))

	res := map[string]string{"message": "song added successfully"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.Log.Debug("Failed to encode response", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	h.Log.Debug("Incoming request to GetAll endpoint")

	filtres, err := ValidFiltres(r)
	if err != nil {
		h.Log.Debug("Failed to validate filters", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.Log.Debug("Filters validated successfully", zap.Any("filters", filtres))

	songs, err := h.Service.GetAll(r.Context(), filtres)
	if err != nil {
		h.Log.Debug("service GetAll", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// for _, item := range songs {
	// res := fmt.Sprintf("song: %s\ngroup: %s\ntext: %s\nlink: %s\ndata_release: %s\n", item.Song, item.Group, item.Text, item.Link, item.Date)
	// w.Write([]byte(res))
	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		h.Log.Debug("Failed to encode response", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	// }

	h.Log.Debug("Response sent successfully")
}

func ValidFiltres(r *http.Request) (models.Filters, error) {

	var filters models.Filters

	val := r.FormValue("limit")
	if val == "" {
		filters.Limit = 5
	} else {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return models.Filters{}, fmt.Errorf("failed conversion limit into int: %w", err)
		}
		if limit < 0 {
			return models.Filters{}, fmt.Errorf("invalid limit: negative")
		}
		filters.Limit = limit
	}

	val = r.FormValue("offset")
	if val == "" {
		filters.Offset = 0
	} else {
		offset, err := strconv.Atoi(val)
		if err != nil {
			return models.Filters{}, fmt.Errorf("failed conversion offset: %w", err)
		}
		if offset < 0 {
			return models.Filters{}, fmt.Errorf("invalid offset: negative")
		}
	}

	filters.Song = r.FormValue("song")
	filters.Group = r.FormValue("group")
	filters.Text = r.FormValue("text")
	filters.Link = r.FormValue("link")
	filters.Date = r.FormValue("date")

	return filters, nil
}

func ValidID(r *http.Request) (int, error) {
	path := r.URL.Path
	paramPath := strings.Split(path, "/")
	val := paramPath[2]
	valByte := []byte(val)
	id, err := strconv.Atoi(string(valByte[0]))
	if err != nil {
		return 0, fmt.Errorf("failed conversion id into int: %w", err)

	}
	return id, nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	h.Log.Debug("Incoming request to Update endpoint")

	var song models.Song

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.Log.Debug("Failed to decode request body", zap.Error(err))
		http.Error(w, `{"error":"failed to decode request"}`, http.StatusBadRequest)
		return
	}
	// val := chi.URLParam(r, "id")
	id, err := ValidID(r)
	if err != nil {
		h.Log.Debug("Converion id", zap.Error(fmt.Errorf("failed conversion id into int: %w", err)))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Log.Debug("Request body decoded successfully", zap.Any("song", song))

	err = h.Service.Update(r.Context(), song, id)
	if err != nil {
		h.Log.Debug("service Update", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := map[string]string{"message": "song updated successfully"}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.Log.Debug("Failed to encode response", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Response sent successfully")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	h.Log.Debug("Incoming request to Delete endpoint")

	// val := chi.URLParam(r, "id")
	id, err := ValidID(r)
	if err != nil {
		h.Log.Debug("Converion id", zap.Error(fmt.Errorf("failed conversion id into int: %w", err)))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Log.Debug("Request body decoded successfully", zap.Int("song", id))

	err = h.Service.Delete(r.Context(), id)
	if err != nil {
		h.Log.Debug("service Delete", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := map[string]string{"message": "song deleted successfully"}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.Log.Debug("Failed to encode response", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Response sent successfully")
}

func (h *Handler) GetText(w http.ResponseWriter, r *http.Request) {

	h.Log.Debug("Incoming request to GetText endpoint")

	filters, err := ValidFiltres(r)
	if err != nil {
		h.Log.Debug("Failed to validate filters", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := ValidID(r)
	if err != nil {
		h.Log.Debug("Converion id", zap.Error(fmt.Errorf("failed conversion id into int: %w", err)))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	text, err := h.Service.GetText(r.Context(), filters, id)
	if err != nil {
		h.Log.Debug("service GetText", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := map[string]string{"text": text}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Log.Debug("Failed to encode response", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Response sent successfully")
}
