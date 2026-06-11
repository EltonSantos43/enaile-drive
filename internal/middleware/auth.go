package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 🔥 IMPORTANTE: Precisa ser exatamente a mesma chave que você usou no usuario.go
var jwtKey = []byte("sua_chave_secreta_super_segura_enaile_2026")

func Autenticacao() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Pega o cabeçalho "Authorization" da requisição
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de acesso não fornecido"})
			c.Abort() // Barra a requisição aqui mesmo!
			return
		}

		// 2. O formato padrão do mercado é: "Bearer TRIPADO_TOKEN_JWT"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Descriptografa e valida o token usando a nossa chave secreta
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		// 4. Extrai o ID do usuário de dentro do token e joga no "Contexto" da requisição
		usuarioIDFl64, ok := claims["usuario_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido (falha ao ler ID)"})
			c.Abort()
			return
		}

		// O JWT salva números como float64, convertemos para uint (padrão do GORM)
		c.Set("usuario_id", uint(usuarioIDFl64))

		// Deixa a requisição seguir em frente para o Handler!
		c.Next()
	}
}