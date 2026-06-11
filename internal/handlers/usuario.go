package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/elton-santos/enaile-drive/internal/database"
	"github.com/elton-santos/enaile-drive/internal/models"
	"github.com/elton-santos/enaile-drive/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


var jwtKey = []byte("sua_chave_secreta_super_segura_enaile_2026")

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

	// Criptografa a senha antes de salvar no banco
	senhaCriptografada, err := bcrypt.GenerateFromPassword([]byte(novoUsuario.Senha), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a senha."})
		return 
	}
	novoUsuario.Senha = string(senhaCriptografada)

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

	// 🔥 ALTERADO: Captura o erro do e-mail e exibe no terminal do servidor
	errEmail := service.EnviarEmailConfirmacao(novoUsuario.Email, novoUsuario.Nome, novoUsuario.TokenAV)
	if errEmail != nil {
		fmt.Println("❌ ERRO NO ENVIO DE EMAIL:", errEmail)
	}

	novoUsuario.Senha = "****"
	c.JSON(http.StatusCreated, gin.H{
		"status":   "sucesso",
		"mensagem": "Cadastro realizado com sucesso. Verifique seu email para ativar sua conta.",
		"usuario":  novoUsuario,
	})
}

func ConfirmarCadastro(c *gin.Context) {
	tokenFiltro := c.Query("token")

	if tokenFiltro == "" {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte("<h1>Erro: Token inválido ou ausente!</h1>"))
		return
	}

	var usuario models.Usuario
	result := database.DB.Where("token_av = ?", tokenFiltro).First(&usuario)
	if result.Error != nil {
		c.Data(http.StatusNotFound, "text/html; charset=utf-8", []byte("<h1>Erro: Cadastro não encontrado ou token expirado!</h1>"))
		return
	}

	usuario.Ativo = true
	usuario.TokenAV = "" 
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

func LoginUsuario(c *gin.Context) {
	var dadosLogin struct {
		Email string `json:"email" binding:"required"`
		Senha string `json:"senha" binding:"required"`
	}

	if err := c.ShouldBindJSON(&dadosLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "E-mail e senha são obrigatórios"})
		return
	}

	var usuario models.Usuario
	if err := database.DB.Where("email = ?", dadosLogin.Email).First(&usuario).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-mail ou senha incorretos"})
		return
	}

	if !usuario.Ativo {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sua conta ainda não foi ativada. Verifique seu e-mail!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(dadosLogin.Senha)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-mail ou senha incorretos"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usuario_id": usuario.ID,
		"exp":        expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar o token de acesso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"usuario": gin.H{
			"id":      usuario.ID,
			"nome":    usuario.Nome,
			"veiculo": usuario.Veiculo,
		},
	})
}