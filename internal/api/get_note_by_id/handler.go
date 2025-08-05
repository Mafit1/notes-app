package getnotebyid

import (
	"errors"
	"net/http"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
	"github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/labstack/echo/v4"
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

type Request struct {
	ID int64 `param:"id" validate:"required"`
}

type Responce struct {
	Note Note `json:"note"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	note, err := h.noteService.GetByID(c.Request().Context(), in.ID)
	if err != nil {
		if errors.Is(err, notes.ErrNoteNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	noteResponce := Responce{
		Note: Note{
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
		},
	}

	return c.JSON(http.StatusOK, noteResponce)
}
