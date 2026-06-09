package database

import (
	"github.com/elton-santos/enaile-drive/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("enaile_drive.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados: ", err)
	}

	err = DB.AutoMigrate(&models.Corrida{}, &models.Gasto{}, &models.Usuario{})
	if err != nil {
		log.Fatal("Falha ao executar migration: ", err)
	}
}