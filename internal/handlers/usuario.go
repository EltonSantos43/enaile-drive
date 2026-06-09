package handlers

import (
	"crypto/rand"
	"fmt"
	"github.com/elton-santos/enaile-drive/internal/database"
	"github.com/elton-santos/enaile-drive/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func gerarToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func  CadastroUsuario(c *gin.Context)  {
	var novoUsuario models.Usuario

	if err := c.ShouldBindJSON(&novoUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "erro",
			"mensagem": "Dados inválidos. Verifique o formato do e-mail e do telefone.",
			"detalhe": err.Error(),
		})
		return
	}

	novoUsuario.TokenAV = gerarToken()
	novoUsuario.Ativo = false

	result := database.DB.Create(&novoUsuario)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": "erro",
			"mensagem": "Email ou Telefonejá cadastrados no sistema.",
		})
		return
	}

	novoUsuario.Senha = "****"
	c.JSON(http.StatusCreated, gin.H{
		"status": "sucesso",
		"mensagem": "Cadastro realizado com sucesso. Verifique seu email para ativar sua conta.",
		"usuario": novoUsuario,
	})	
}