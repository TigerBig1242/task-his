package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bigthamm/task-his/internal/dto"
	"github.com/bigthamm/task-his/internal/service"
)

type StaffHandler struct {
	staffService service.StaffService
}

func NewStaffHandler(staffService service.StaffService) *StaffHandler {
	return &StaffHandler{staffService: staffService}
}

func (handler *StaffHandler) Create(c *gin.Context) {
	var request dto.CreateStaffRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := handler.staffService.Create(c.Request.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrHospitalNotFound):
			c.JSON(http.StatusNotFound, "Hospital not found")
			return
		case errors.Is(err, service.ErrUsernameAlreadyTaken):
			c.JSON(http.StatusConflict, "Username already taken")
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	c.JSON(http.StatusCreated, resp)
}

func (handler *StaffHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := handler.staffService.Login(c.Request.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, "Invalid username, password or hospital")
			return
		case errors.Is(err, service.ErrAccountInactive):
			c.JSON(http.StatusForbidden, "Account is active")
			return
		default:
			c.JSON(http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	c.JSON(http.StatusOK, resp)

}
