package routes

import (
	"github.com/gin-gonic/gin"
	"lethal_company_save_manager/save_historizer"
	"net/http"
)

func GetSaves(c *gin.Context) {
	jsonResp, err := save_historizer.ListKnownSaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonResp)
}
