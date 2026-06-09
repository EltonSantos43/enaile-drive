package main

import (
	"github.com/elton-santos/enaile-drive/internal/database"
	"github.com/elton-santos/enaile-drive/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()

	database.InitDB()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/corridas", handlers.PostCorrida)
		v1.GET("/resumo_diario", handlers.GetResumoDiario)
		v1.POST("/usuarios", handlers.CadastroUsuario)
		v1.GET("/usuarios/confirmar", handlers.ConfirmarCadastro)
	}

	log.Println("Enaile Drive rodando em http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Erro ao rodar o servidor: ", err)
	}
}