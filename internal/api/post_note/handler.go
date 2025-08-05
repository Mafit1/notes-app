package postnote

import (
	"net/http"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
	"github.com/Mafit1/notes-app/internal/models"
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
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type Responce struct {
	ID int64 `json:"id"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	note := models.Note{
		Title:   in.Title,
		Content: in.Content,
	}

	id, err := h.noteService.Create(c.Request().Context(), note)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Responce{ID: id})
}
