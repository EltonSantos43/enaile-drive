package main

import (
	"log"

	"github.com/elton-santos/enaile-drive/internal/database"
	"github.com/elton-santos/enaile-drive/internal/handlers"
	"github.com/elton-santos/enaile-drive/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	database.InitDB()

	r := gin.Default()

	r.Use(cors.Default())

	v1 := r.Group("/api/v1")
	{
		// 🔓 ROTAS PÚBLICAS: Acessíveis por qualquer um sem token
		v1.POST("/usuarios", handlers.CadastroUsuario)
		v1.POST("/login", handlers.LoginUsuario)
		v1.GET("/usuarios/confirmar", handlers.ConfirmarCadastro)

		// 🔒 ROTAS PROTEGIDAS: Todas usam a variável 'protegido' para exigir o Token JWT
		protegido := v1.Group("/")
		protegido.Use(middleware.Autenticacao())
		{
			protegido.GET("/resumo_diario", handlers.GetResumoDiario)
			protegido.POST("/corridas", handlers.PostCorrida)
			protegido.POST("/gastos", handlers.PostGasto)
			protegido.PUT("/usuarios/veiculo", handlers.AtualizarVeiculo)
		}
	}

	log.Println("Enaile Drive rodando em http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Erro ao rodar o servidor: ", err)
	}
}