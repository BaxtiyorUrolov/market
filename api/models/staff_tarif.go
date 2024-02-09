package models

import "time"

type StaffTarif struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	TarifType     string `json:"tarif_type"`
	AmountForCash int    `json:"amount_for_cash"`
	AmountForCard int    `json:"amount_for_carsd"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type CreateStaffTarif struct {
	Name          string `json:"name"`
	TarifType     string `json:"tarif_type"`
	AmountForCash int    `json:"amount_for_cash"`
	AmountForCard int    `json:"amount_for_card"`
}

type UpdateStaffTarif struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	TarifType     string `json:"tarif_type"`
	AmountForCash int    `json:"amount_for_cash"`
	AmountForCard int    `json:"amount_for_carsd"`
}

type StaffTarifResponse struct {
	StaffTarifs   []StaffTarif `json:"staff_tarifs"`
	Count         int          `json:"count"`
}