package models

type MessageResponse struct {
	Message string `json:"message"`
}

type KeyResponse struct {
	Key string `json:"key"`
}

type BooleanResponse struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

type BooleansResponse struct {
	Names []string `json:"names"`
}
