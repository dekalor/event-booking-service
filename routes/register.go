package routes

import (
	"net/http"
	"strconv"

	"example.com/event-booking/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Event ID is not valid"})
		return
	}

	event, err := models.GetEventById(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Event not found"})
		return
	}

	err = event.Register(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Register"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Register Success"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("user_id")
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)

	var event models.Event
	event.ID = eventId

	err = event.CancelRegister(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to Cancel Registration"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Cancel Registration Success"})
}
