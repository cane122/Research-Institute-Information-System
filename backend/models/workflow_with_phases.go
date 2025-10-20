package models

type CreateWorkflowWithPhasesRequest struct {
	Naziv   string   `json:"naziv"`
	TipToka string   `json:"tip_toka"`
	Faze    []string `json:"faze"`
}
