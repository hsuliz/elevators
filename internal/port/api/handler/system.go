package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hsuliz/elevators/internal/domain"
)

type System struct {
	domainSystem *domain.System
}

func NewSystem(domainSystem *domain.System) *System {
	return &System{domainSystem}
}

func (s System) CallElevator(c *gin.Context) {
	floorParam := c.Param("floor")
	floorNumber, err := strconv.Atoi(floorParam)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid floor number",
		})
		return
	}

	s.domainSystem.Call(floorNumber)
	c.Status(http.StatusOK)
}
