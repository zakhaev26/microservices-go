package model

type AdminData struct {
	Scorer     string `json:"scorer" bson:"scorer,omitempty"`
	ScorerTeam string `json:"scorer_team" bson:"scorer_team,omitempty"`
	RivalTeam  string `json:"rival_team" bson:"rival_team,omitempty"`
	TeamA      string `json:"teama" bson:"teama,omitempty"`
	TeamB      string `json:"teamb" bson:"teamb,omitempty"`
	
}
