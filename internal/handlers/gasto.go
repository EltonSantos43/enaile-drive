package handlers

import (
	"net/http"
	"time"

	"github.com/elton-santos/enaile-drive/internal/database"
	"github.com/elton-santos/enaile-drive/internal/models"
	"github.com/gin-gonic/gin"
)

func PostGasto(c *gin.Context) {
	var  novoGasto models.Gasto

	if err := c.ShouldBindJSON(&novoGasto); err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novoGasto.UsuarioID = 1

	if  novoGasto.DataCustomizada != "" {
		parsedDate, err := time.Parse("2006-01-02", novoGasto.DataCustomizada)
		if err == nil {
			novoGasto.CreatedAt = parsedDate
		}
	}

	result := database.DB.Create(&novoGasto)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, novoGasto)
}