package dto

type SecurityEncryptPayload struct {
	Data      string `json:"data"`
	PartnerID uint64 `json:"partner_id"`
}

type SecurityEncryptResponse struct {
	Data      string `json:"data"`
	Tag       string `json:"tag"`
	SecretKey string `json:"secret_key"`
}

type SecurityDecryptPayload struct {
	Data        string `json:"data"`
	PartnerID   uint64 `json:"partner_id"`
	Tag         string `json:"tag"`
	KeypairHash string `json:"keypair_hash"`
}

type SecurityDecryptResponse struct {
	Data string `json:"data"`
}
