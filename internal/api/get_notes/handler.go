package getnotes

import (
	"net/http"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
	"github.com/Mafit1/notes-app/internal/models"
	"github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type handler struct {
	noteService notes.Service
}

func New(noteService notes.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{noteService})
}

type Note struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Request struct{}

type Responce struct {
	Notes []Note `json:"notes"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	out, err := h.noteService.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	notes := lo.Map(out, func(note models.Note, i int) Note {
		return Note{
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
		}
	})

	return c.JSON(http.StatusOK, Responce{notes})
}
