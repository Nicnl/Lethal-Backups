package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func GetPlanets(c *gin.Context) {
	jsonResp, err := os.ReadFile(filepath.Join("resources", "planets.json"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResp)
}
