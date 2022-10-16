package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Note struct {
	Id       string    `json:"id"`
	Owner    string    `json:"owner"`
	Content  string    `json:"content"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

type Notes []Note

type dbConn interface {
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func GetNotesForOwner(ctx context.Context, conn dbConn, owner string) (Notes, error) {
	if owner == "" {
		return nil, errors.New("model: owner not supplied")
	}

	queryRows, err := conn.Query(ctx, "SELECT id, owner, content, created, modified FROM public.note")
	if err != nil {
		return nil, fmt.Errorf("model: could not query notes: %w", err)
	}
	defer queryRows.Close()

	notes := []Note{}
	for queryRows.Next() {
		row := Note{}
		err = queryRows.Scan(&row.Id, &row.Owner, &row.Content, &row.Created, &row.Modified)
		if err != nil {
			return nil, fmt.Errorf("model: query scan failed: %w", err)
		}
		if row.Owner == owner {
			notes = append(notes, row)
		}
	}

	if queryRows.Err() != nil {
		return nil, fmt.Errorf("model: query read failed: %w", queryRows.Err())
	}

	return notes, nil
}

func GetNoteById(ctx context.Context, conn dbConn, id string) (Note, error) {
	var note Note
	if id == "" {
		return note, errors.New("model: id not supplied")
	}

	row := conn.QueryRow(ctx, "SELECT id, owner, content, created, modified FROM public.note WHERE id = $1", id)

	err := row.Scan(&note.Id, &note.Owner, &note.Content, &note.Created, &note.Modified)
	if err != nil {
		return note, fmt.Errorf("model: query scan failed: %w", err)
	}
	return note, nil
}
