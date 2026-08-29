package handler

type generateRequest struct {
	DiscordID   int64   `json:"discord_id"`
	AddCard     *string `json:"add_card,omitempty"`
	CurrentCard *string `json:"current_card,omitempty"`
}
