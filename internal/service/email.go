package service

import (
	"fmt"
	"net/smtp"
	"os" // 🔥 ADICIONADO: Pacote nativo do Go para ler variáveis do sistema (.env)
)

func EnviarEmailConfirmacao(emailDestino, nomeUsuario, token string) error {
	from := "no-reply@enailedrive.com.br"
	host := "sandbox.smtp.mailtrap.io"
	port := "2525"

	// 🔥 ALTERADO: Agora o Go não usa texto chumbado, ele pega direto do seu arquivo .env protegido
	usuarioMailtrap := os.Getenv("MAILTRAP_USER")
	senhaMailtrap := os.Getenv("MAILTRAP_PASSWORD")

	// 💡 Dica extra: Corrigido o "htpp" para "http" e a porta 800 para 8080 para o link do botão funcionar depois!
	linkAtivacao := fmt.Sprintf("http://localhost:8080/api/v1/usuarios/confirmar?token=%s", token)

	assunto := "Subject: Enaile Drive - Confirme seu Cadastro! 🚗\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	
	corpoHtml := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; color: #333;">
				<h2>Olá, %s!</h2>
				<p>Seja muito bem-vindo ao <strong>Enaile Drive</strong>, o seu gerenciador de backend automotivo.</p>
				<p>Para ativar sua conta e começar a registrar suas corridas, clique no botão abaixo:</p>
				<p style="margin: 25px 0;">
					<a href="%s" style="background-color: #25d366; color: white; padding: 12px 20px; text-decoration: none; border-radius: 5px; font-weight: bold;">
						ATIVAR MINHA CONTA
					</a>
				</p>
				<p style="font-size: 12px; color: #777;">Se o botão não funcionar, copie e cole este link no seu navegador: %s</p>
			</body>
		</html>
	`, nomeUsuario, linkAtivacao, linkAtivacao)

	mensagem := []byte(assunto + mime + corpoHtml)

	// 4. Faz a autenticação nos servidores do Mailtrap
	auth := smtp.PlainAuth("", usuarioMailtrap, senhaMailtrap, host)

	// 5. Dispara o e-mail de fato
	err := smtp.SendMail(host+":"+port, auth, from, []string{emailDestino}, mensagem)
	if err != nil {
		return err
	}

	return nil
}