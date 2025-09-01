package getnotesbyuserid

import (
	"net/http"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
	"github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	notesService notes.Service
}

func New(notesService notes.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{notesService})
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
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid userID in context")
	}

	notes, err := h.notesService.GetAllByUserID(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, notes)
}
