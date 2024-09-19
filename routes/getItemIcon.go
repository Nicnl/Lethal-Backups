package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func GetItemIcon(c *gin.Context) {
	itemName := c.Param("itemName")

	// Prevent from path traversal
	if filepath.Base(itemName) != itemName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid icon name"})
		return
	}

	c.File(filepath.Join("resources", "items_icons", itemName))
}
