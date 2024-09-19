package routes

import (
	"github.com/gin-gonic/gin"
	"lethal_company_save_manager/save_historizer"
	"net/http"
)

func PostRestoreHash(c *gin.Context) {
	hash := c.Param("hash")

	save, ok := save_historizer.ObtainSave(hash)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "save not found"})
		return
	}

	save.Restore()
}
