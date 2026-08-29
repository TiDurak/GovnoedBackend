package handler

type generateResponse struct {
	CurrentCard string   `json:"current_card"`
	Cards       []string `json:"cards"`
}
