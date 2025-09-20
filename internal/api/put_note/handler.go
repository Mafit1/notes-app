package patchnote

import (
	"errors"
	"net/http"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
	"github.com/Mafit1/notes-app/internal/models"
	"github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/google/uuid"
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
	ID      int64  `param:"id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type Response struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid userID in context")
	}

	note := models.Note{
		ID:      in.ID,
		Title:   in.Title,
		Content: in.Content,
	}

	updatedNote, err := h.noteService.Update(c.Request().Context(), userID, note)
	if err != nil {
		if errors.Is(err, notes.ErrNoteNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if errors.Is(err, notes.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	noteResponse := Response{
		ID:      updatedNote.ID,
		Title:   updatedNote.Title,
		Content: updatedNote.Content,
	}

	return c.JSON(http.StatusOK, noteResponse)
}
