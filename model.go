package fazpass

import "time"

type Data struct {
	SessionId string     `json:"session_id"`
	TimeStamp *time.Time `json:"time_stamp"`
	Device    Device     `json:"device"`
}

type Geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Device struct {
	Score           float64     `json:"score"`
	Platform        string      `json:"platform"`
	FazpassId       string      `json:"fazpass_id"`
	IsRooted        bool        `json:"is_rooted"`
	IsEmulator      bool        `json:"is_emulator"`
	IsGpsSpoof      bool        `json:"is_gps_spoof"`
	IsAppTemper     bool        `json:"is_app_temper"`
	IsVpn           bool        `json:"is_vpn"`
	IsScreenSharing bool        `json:"is_share_screen"`
	IsDebuging      bool        `json:"is_debuging"`
	Name            string      `json:"name"`
	SimSerial       []string    `json:"sim_serial"`
	Geolocation     Geolocation `json:"geolocation"`
}

type Transmission struct {
	Message string `json:"message"`
}

type CheckRequest struct {
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required"`
	Data  string `json:"data" validate:"required"`
}

type EnrollRequest struct {
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required"`
	Data  string `json:"data" validate:"required"`
}

type ValidateRequest struct {
	FazpassId string `json:"fazpass_id" validate:"required"`
	Data      string `json:"data" validate:"required"`
}

type RemoveRequest struct {
	FazpassId string `json:"fazpass_id" validate:"required"`
	Data      string `json:"data" validate:"required"`
}

type CheckResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
