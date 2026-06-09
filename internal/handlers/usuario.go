package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/elton-santos/enaile-drive/internal/database"
	"github.com/elton-santos/enaile-drive/internal/models"
	"github.com/elton-santos/enaile-drive/internal/service"
	"github.com/gin-gonic/gin"
)

func gerarToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func CadastroUsuario(c *gin.Context) {
	var novoUsuario models.Usuario

	if err := c.ShouldBindJSON(&novoUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":   "erro",
			"mensagem": "Dados inválidos. Verifique o formato do e-mail e do telefone.",
			"detalhe":  err.Error(),
		})
		return
	}

	novoUsuario.TokenAV = gerarToken()
	novoUsuario.Ativo = false

	result := database.DB.Create(&novoUsuario)
	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":   "erro",
			"mensagem": "Email ou Telefone já cadastrados no sistema.",
		})
		return
	}

	_ = service.EnviarEmailConfirmacao(novoUsuario.Email, novoUsuario.Nome, novoUsuario.TokenAV)

	novoUsuario.Senha = "****"
	c.JSON(http.StatusCreated, gin.H{
		"status":   "sucesso",
		"mensagem": "Cadastro realizado com sucesso. Verifique seu email para ativar sua conta.",
		"usuario":  novoUsuario,
	})
} // <-- CORRIGIDO: Agora fechamos a função CadastroUsuario corretamente!

func ConfirmarCadastro(c *gin.Context) {
	tokenFiltro := c.Query("token")

	if tokenFiltro == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte("<h1>Erro: Token inválido ou ausente!</h1>"))
		return
	}

	// CORRIGIDO: Criamos a variável e buscamos o usuário dono do token no banco de dados
	var usuario models.Usuario
	result := database.DB.Where("token_av = ?", tokenFiltro).First(&usuario)
	if result.Error != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte("<h1>Erro: Cadastro não encontrado ou token expirado!</h1>"))
		return
	}

	usuario.Ativo = true
	usuario.TokenAV = "" // Limpa o token por segurança
	database.DB.Save(&usuario)

	htmlSucesso := fmt.Sprintf(`
    <html>
        <body style="font-family: Arial, sans-serif; text-align: center; padding: 50px; color: #333;">
            <div style="max-width: 500px; margin: 0 auto; border: 1px solid #e0e0e0; padding: 30px; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.1);">
                <h1 style="color: #25d366;">✓ Conta Ativada com Sucesso!</h1>
                <p style="font-size: 18px; margin: 20px 0;">Olá, <strong>%s</strong>!</p>
                <p>Sua conta no <strong>Enaile Drive</strong> está confirmada e liberada.</p>
                <p>Você já pode fechar esta aba e começar a cadastrar suas corridas.</p>
            </div>
        </body>
    </html>
    `, usuario.Nome)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlSucesso))
}