package models

import "gorm.io/gorm"

type Gasto struct {
	gorm.Model
	UsuarioID uint `json:"usuario_id"`
	Valor float64 `json:"valor" binding:"required"`
	Categoria string `json:"categoria" binding:"required"`
	Descricao string `json:"descricao"`
	DataCustomizada string `json:"data_customizada" gorm:"-"`
}