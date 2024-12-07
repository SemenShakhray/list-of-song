package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/SemenShakhray/list-of-song/internal/models"

	"go.uber.org/zap"
)

func (s *Store) AddSong(ctx context.Context, song models.Song) error {

	s.Log.Debug("Attempting to add song",
		zap.String("song", song.Song),
		zap.String("group", song.Group),
	)

	query := `INSERT INTO songs (song, group_name, text, link, date_release) VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (song, group_name) DO NOTHING;`

	row, err := s.DB.ExecContext(ctx, query, song.Song, song.Group, song.Text, song.Link, song.Date)
	if err != nil {
		s.Log.Debug("Add song", zap.Error(err))
		return fmt.Errorf("failed to add song in storage: %w", err)
	}
	n, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if n == 0 {
		s.Log.Warn("Song already exists in the storage",
			zap.String("song", song.Song),
			zap.String("group", song.Group),
		)
	}
	s.Log.Debug("Song successfully added")
	return nil

}

func (s *Store) Update(ctx context.Context, song models.Song) error {

	s.Log.Debug("Updating info about song",
		zap.Int("song", song.Id),
	)

	query := `UPDATE songs SET 
	text = COALESCE(NULLIF($1, ''), text),
    link = COALESCE(NULLIF($2, ''), link),
    date_release = COALESCE(NULLIF($3, ''), date_release)
	WHERE id = $4`

	row, err := s.DB.ExecContext(ctx, query, song.Text, song.Link, song.Date, song.Id)
	if err != nil {
		return fmt.Errorf("failed to update info about song: %w", err)
	}
	n, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if n == 0 {
		s.Log.Debug("No changes made to the song info", zap.Any("song", song))
		return fmt.Errorf("song not found")
	} else {
		s.Log.Debug("Song info successfully updated")
	}
	return nil
}

func (s *Store) GetAll(ctx context.Context, filters models.Filters) ([]models.Song, error) {

	s.Log.Debug("Get songs", zap.Any("filters_song", filters))

	query := `SELECT id, song, group_name, text, link, date_release FROM songs
WHERE (song ILIKE '%' || $1 || '%') AND 
( group_name ILIKE '%' || $2 || '%') AND 
(text ILIKE '%' || $3 || '%') AND 
(link ILIKE '%' || $4 || '%') AND 
(date_release = $5)
LIMIT $6 OFFSET $7;`

	var songs []models.Song
	rows, err := s.DB.QueryContext(ctx, query,
		filters.Song,
		filters.Group,
		filters.Text,
		filters.Link,
		filters.Date,
		filters.Limit,
		filters.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query songs: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.Id, &song.Song, &song.Group, &song.Text, &song.Link, &song.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan song: %w", err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through songs: %w", err)
	}
	s.Log.Debug("List of songs has been successfully received", zap.Int("total", len(songs)))
	return songs, nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	s.Log.Debug("Attempting to delete song",
		zap.Int("song", id),
	)

	query := "DELETE FROM songs WHERE id = $1;"
	row, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	n, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %w", err)
	}
	if n == 0 {
		s.Log.Debug("No songs deleted",
			zap.Int("song", id),
		)
		return fmt.Errorf("failed delete the song")
	}
	return nil
}

func (s *Store) GetText(ctx context.Context, filters models.Filters, id int) (string, error) {
	s.Log.Debug("Attemting get text of song")

	query := "SELECT text FROM songs WHERE id = $1"
	row := s.DB.QueryRowContext(ctx, query, id)

	var text string
	err := row.Scan(&text)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Log.Warn("Song not found", zap.String("song", filters.Song))
			return "", fmt.Errorf("failed to get text of the song: %w", err)
		}
		return "", fmt.Errorf("failed to retrieve text of the song: %w", err)
	}

	verses := strings.Split(text, "\n")
	if filters.Offset >= len(verses) {
		s.Log.Debug("Offset exceeds number of verses", zap.Int("offset", filters.Offset), zap.Int("len(verses)", len(verses)))
		text = ""
		return text, nil
	}
	end := filters.Offset + filters.Limit
	if end > len(verses) {
		end = len(verses)
	}
	text = strings.Join(verses[filters.Offset:end], "\n")

	s.Log.Debug("Text of the song successfully received", zap.String("song", filters.Song))
	return text, nil
}
