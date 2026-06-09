package models

import "gorm.io/gorm"

type Usuario struct {
	gorm.Model
	Nome     string `json:"nome" binding:"required"`
	Email    string `json:"email" binding:"required,email" gorm:"uniqueIndex"`
	Telefone string `json:"telefone" binding:"required,e164" gorm:"uniqueIndex"`
	Senha    string `json:"senha" binding:"required,min=6"`
	Ativo    bool   `json:"ativo" gorm:"default:false"`
	TokenAV  string `json:"token_av"` 
}