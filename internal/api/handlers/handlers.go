package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/SemenShakhray/list-of-song/internal/config"
	"github.com/SemenShakhray/list-of-song/internal/models"
	"github.com/SemenShakhray/list-of-song/internal/service"

	"go.uber.org/zap"
)

type Handler struct {
	Log     *zap.Logger
	Service service.Servicer
	Cfg     config.Config
}

func NewHandler(log *zap.Logger, serv service.Servicer) Handler {
	return Handler{
		Log:     log,
		Service: serv,
	}
}

// AddSong adds a new song to the database
//
//	@Summary		Add a new song
//	@Description	Add a new song by providing its details. It also fetches additional info from an external API
//	@Tags			Songs
//	@Accept			json
//	@Produce		json
//	@Param			song	body		models.Song	true		"Song details"
//	@Success		200		{object}	map[string]string	"Song added successfully"
//	@Failure		400		{object}	map[string]string	"Invalid input"
//	@Failure		500		{object}	map[string]string	"Failed add song"
//	@Router			/songs [post]
func (h *Handler) AddSong(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("Incoming request to AddSong endpoint")

	var song models.Song
	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.Log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, `{"error":"failed to decode request"}`, http.StatusBadRequest)
		return
	}
	h.Log.Debug("Song data", zap.String("song", song.Song), zap.String("group", song.Group))

	h.Log.Debug("calling an external API")
	if h.Cfg.API.Call {
		hostPort := net.JoinHostPort(h.Cfg.API.Host, h.Cfg.API.Port)
		apiURL := fmt.Sprintf("http://%s/info?group=%s&song=%s", hostPort, url.QueryEscape(song.Group), url.QueryEscape(song.Song))

		resp, err := http.Get(apiURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			h.Log.Error("Failed to fetch song details from external API", zap.Error(err))
			http.Error(w, `{"error":"failed to fetch song details from external API"}`, http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&song); err != nil {
			h.Log.Error("Failed to decode API response", zap.Error(err))
			http.Error(w, `{"error":"failed to decode API response"}`, http.StatusInternalServerError)
			return
		}
	}

	err = h.Service.AddSong(r.Context(), song)
	if err != nil {
		h.Log.Error("service AddSong", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Song added successfully", zap.Any("song", song))

	res := map[string]string{"message": "song added successfully"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.Log.Error("Failed to encode response", zap.Error(err))
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

// GetAll returns a list of Song's
//
//	@Summary		Get all songs
//	@Description	Retrieve a list of songs with optional filters
//	@Tags			Songs
//	@Accept			json
//	@Produce		json
//	@Param			song			query		string		false	"Filter by song name (partial match)"
//	@Param			group			query		string		false	"Filter by group name (partial match)"
//	@Param			text			query		string		false	"Filter by lyrics (partial match)"
//	@Param			link			query		string		false	"Filter by link (partial match)"
//	@Param			date_release	query		string		false	"Filter by release date (format: YYYY-MM-DD)"
//	@Param			limit			query		integer		false	"Limit the number of results"
//	@Param			offset			query		integer		false	"Offset for pagination"
//	@Success		200				{array}		models.Song	"List of songs"
//	@Failure		400				{object}	map[string]string	"Invalid filters provided"
//	@Failure		500				{object}	map[string]string	"Internal server error"
//	@Router			/songs [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("Incoming request to GetAll endpoint")

	filtres, err := ValidFiltres(r)
	if err != nil {
		h.Log.Error("Failed to validate filters", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	h.Log.Debug("Filters validated successfully", zap.Any("filters", filtres))

	songs, err := h.Service.GetAll(r.Context(), filtres)
	if err != nil {
		h.Log.Error("service GetAll", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(songs)
	if err != nil {
		h.Log.Error("Failed to encode response", zap.Error(err))
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}

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
		filters.Offset = offset
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

// Update updates an existing song
//
//	@Summary		Update a song
//	@Description	Update details of an existing song by its ID.
//	@Tags			Songs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Song ID"
//	@Param			song	body		models.Song			true	"Updated song details"
//	@Success		200		{object}	map[string]string	"Song updated successfully"
//	@Failure		400		{object}	map[string]string	"Invalid request body or ID"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/songs/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("Incoming request to Update endpoint")

	var song models.Song

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.Log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, `{"error":"failed to decode request"}`, http.StatusBadRequest)
		return
	}
	// val := chi.URLParam(r, "id")
	id, err := ValidID(r)
	if err != nil {
		h.Log.Error("Converion id", zap.Error(fmt.Errorf("failed conversion id into int: %w", err)))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	h.Log.Debug("Request body decoded successfully", zap.Any("song", song))

	song.Id = id
	err = h.Service.Update(r.Context(), song)
	if err != nil {
		h.Log.Error("service Update", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	res := map[string]string{"message": "song updated successfully"}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.Log.Error("Failed to encode response", zap.Error(err))
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Response sent successfully")
}

// Delete deletes a song from the database by its ID
//
//	@Summary		Delete a song
//	@Description	Delete an existing song by its ID
//	@Tags			Songs
//	@Param			id	path		int					true	"Song ID"
//	@Success		200	{object}	map[string]string	"Song deleted successfully"
//	@Failure		400	{object}	map[string]string	"Invalid song ID"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/songs/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("Incoming request to Delete endpoint")

	// val := chi.URLParam(r, "id")
	id, err := ValidID(r)
	if err != nil {
		h.Log.Error("Converion id", zap.Error(fmt.Errorf("failed conversion id into int: %w", err)))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	h.Log.Debug("Request body decoded successfully", zap.Int("song", id))

	err = h.Service.Delete(r.Context(), id)
	if err != nil {
		h.Log.Error("service Delete", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	res := map[string]string{"message": "song deleted successfully"}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		h.Log.Error("Failed to encode response", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Response sent successfully")
}

// GetText возвращает текст песни по её ID и фильтрам.
//
//	@Summary		Get song text
//	@Description	Retrieve the text of a song based on filters and its ID
//	@Tags			Songs
//	@Param			id		path		int					true	"Song ID"
//	@Param			group	query		string				false	"Filter by group"
//	@Param			song	query		string				false	"Filter by song"
//	@Success		200		{object}	map[string]string	"Song text retrieved successfully"
//	@Failure		400		{object}	map[string]string	"Invalid input"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/songs/{id} [get]
func (h *Handler) GetText(w http.ResponseWriter, r *http.Request) {
	h.Log.Debug("Incoming request to GetText endpoint")

	filters, err := ValidFiltres(r)
	if err != nil {
		h.Log.Error("Failed to validate filters", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	id, err := ValidID(r)
	if err != nil {
		h.Log.Error("Converion id", zap.Error(fmt.Errorf("failed conversion id into int: %w", err)))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	text, err := h.Service.GetText(r.Context(), filters, id)
	if err != nil {
		h.Log.Error("service GetText", zap.Error(err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
	}

	response := map[string]string{"text": text}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.Log.Error("Failed to encode response", zap.Error(err))
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	h.Log.Debug("Response sent successfully")
}
