package models

import (
	"fmt"
	"time"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type Chapter struct {
	ID                   uuid.UUID `json:"id"`
	CreatedAt            time.Time `json:"_"`
	UpdatedAt            time.Time `json:"_"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	NumberOfPages        int64     `json:"number_of_pages"`
	ChapterNumber        int64     `json:"chapter_number"`
	CollectionID         uuid.UUID `json:"chapter"`
}

func (i *Chapter) Create(conn *pgx.Conn) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty.")
	}

	now := time.Now()
	row := conn.QueryRow(context.Background(), "INSERT INTO chapter (title, description, number_of_pages, chapter_number, collection_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", i.Title, i.Description, i.NumberOfPages, i.ChapterNumber, i.CollectionID, now, now)

	err := row.Scan(&i.ID)

	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("There was an error creating the chapter")
	}

	return nil
}

func (i *Chapter) Update(conn *pgx.Conn) error {
	i.Title = strings.Trim(i.Title, " ")
	if len(i.Title) < 1 {
		return fmt.Errorf("Title must not be empty")
	}

	now := time.Now()
	_, err := conn.Exec(context.Background(), "UPDATE chapter SET title=$1, description=$2, number_of_pages=$3, chapter_number=$4, updated_at=$5 WHERE id=$6", i.Title, i.Description, i.NumberOfPages, i.ChapterNumber, now, i.ID)

	if err != nil {
		fmt.Printf("Error updating chapter: (%v)", err)
		return fmt.Errorf("Error updating chapter")
	}

	return nil
}

func FindChapterById(id uuid.UUID, conn *pgx.Conn) (Chapter, error) {
	row := conn.QueryRow(context.Background(), "SELECT title, description, number_of_pages, chapter_number, collection_id FROM chapter WHERE id=$1", id)
	chapter := Chapter{
		ID: id,
	}
	err := row.Scan(&chapter.Title, &chapter.Description, &chapter.NumberOfPages, &chapter.CollectionID, &chapter.ChapterNumber)
	if err != nil {
		return chapter, fmt.Errorf("The chapter doesn't exist")
	}

	return chapter, nil
}

func (i *Chapter) Delete(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM chapter WHERE id=$5", i.ID)

	if err != nil {
		fmt.Printf("Error deleting chapter: (%v)", err)
		return fmt.Errorf("Error deleting chapter")
	}

	return nil
}
