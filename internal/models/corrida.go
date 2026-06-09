package models

import "gorm.io/gorm"

type Corrida struct {
	gorm.Model
	UsuarioID           uint    `json:"usuario_id"`
	VeiculoAtivo        string  `json:"veiculo_ativo"` 
	ValorRecebido       float64 `json:"valor_recebido"`
	TipoPagamento       string  `json:"tipo_pagamento"` 
	DistanciaPercorrida float64 `json:"distancia_km"`
	DuracaoMinutos       int     `json:"duracao_min"`
	HorarioCorrida      string  `json:"horario_corrida"` 
	Plataforma          string  `json:"plataforma"`  
	DataCustomizada     string  `json:"data_customizada" gorm:"-"`    
}

type Gasto struct {
	gorm.Model
	UsuarioID uint    `json:"usuario_id"`
	Descricao string  `json:"descricao"` 
	Valor     float64 `json:"valor"`
}