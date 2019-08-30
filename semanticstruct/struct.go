package semanticstruct

import (
	"kdatapack/connection"
)

type SContext struct {
	File					*SConfig
	DescList				[]string
	NormalizedDescList		[]string
	SplittedDescList		[][]string
	DescLen					int
	Rareword				SRareWord
	Dictionnary				map[string]SDicoword
}

type SConfig struct {
	Connect			connection.Connection	`json:"connection"`
	TableName		string					`json:"tableName"`
	DescriptionVar	string					`json:"descriptionVar"`
	Limit			int						`json:"limit"`
	Resultnb		int						`json:"resultnb"`
	Fuzzywuzzy		SFuzzywuzzyconfig		`json:"fuzzywuzzy"`
	Rareword		SRareWordconfig			`json:"rareword"`
	Dictionnary		SDictionnaryconfig		`json:"dictionnary"`
	Fuzzydico		SFuzzydicoconfig		`json:"fuzzydico"`
}

type SFuzzywuzzyconfig struct {
	Useit			bool	`json:"useIt"`
	Minscore		int		`json:"minscore"`
	Maxscore		int		`json:"maxscore"`
}

type SRareWordconfig struct {
	Useit			bool	`json:"useIt"`
	Minlength		int		`json:"minLength"`
	Percent			float64	`json:"percent"`
}

type SDictionnaryconfig struct {
	Useit			bool	`json:"useIt"`
	Databasename	string	`json:"databasename"`
	Tablename		string	`json:"tablename"`
}

type SFuzzydicoconfig struct {
	Useit			bool	`json:"useIt"`
}

type SRareWord struct {
	Descriptions	[]Sdesc
	Rarewords		map[string]float64
}

type Sdesc struct {
	Desc		string
	Words		map[string]struct{}
}

type SDicoword struct {
	OtherAppear	map[string]uint32
	Affinities	map[string]float32
	Appear		uint32
}