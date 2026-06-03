package handlers

import (
	"github.com/seu-usuario/enaile-drive/internal/database"
	"github.com/seu-usuario/enaile-drive/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostCorrida(c *gin.Context) {
	var novaCorrida models.Corrida
	if err := c.ShouldBindJSON(&novaCorrida); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novaCorrida.UsuarioID = 1 
	novaCorrida.VeiculoAtivo = "Sedan" 

	result := database.DB.Create(&novaCorrida)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, novaCorrida)
}