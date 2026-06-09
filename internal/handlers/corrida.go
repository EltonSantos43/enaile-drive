package handlers

import (
	"github.com/elton-santos/enaile-drive/internal/database"
	"github.com/elton-santos/enaile-drive/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func PostCorrida(c *gin.Context) {
	var novaCorrida models.Corrida
	if err := c.ShouldBindJSON(&novaCorrida); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novaCorrida.UsuarioID = 1 
	novaCorrida.VeiculoAtivo = "Argo"

	if novaCorrida.DataCustomizada != "" {
		parsedDate, err := time.Parse("2006-01-02", novaCorrida.DataCustomizada)
		if err == nil {
			novaCorrida.CreatedAt = parsedDate
		}
	}

	result := database.DB.Create(&novaCorrida)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, novaCorrida)
}

func GetResumoDiario(c *gin.Context) {
	UsuarioID := 1

	dataFiltro := c.Query("data")
	if dataFiltro == "" {
		dataFiltro = time.Now().Format("2006-01-02")
	}

	var corridas []models.Corrida

	database.DB.Where("usuario_id = ? AND DATE(created_at) = ?", UsuarioID, dataFiltro).Find(&corridas)

	var ganhosHoje float64
	for _, corrida := range corridas {
		ganhosHoje += corrida.ValorRecebido
	}

	totalCorridas := len(corridas)

	c.JSON(http.StatusOK, gin.H{
		"data": dataFiltro,
		"ganhos_hoje": ganhosHoje,
		"gastos_hoje": 0.0,
		"lucros_diario": ganhosHoje,
		"total_corridas": totalCorridas,
	})
}