package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"semanticContent/semanticstruct"
)

func loadConfigFile(cfPath string) (*semanticstruct.SConfig, error) {
	config := new(semanticstruct.SConfig)

	file, err := ioutil.ReadFile(cfPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, config)
	if err != nil {
		return config, err
	}
	if err = ValidConfig(config); err != nil {
		return config, err
	}
	return config, nil
}

func ValidConfig(ConfigFile *semanticstruct.SConfig) error {
	algocount := 0
	missing := make([]string, 0)
	if ConfigFile.Connect.BadInit() {
		missing = append(missing, "wrong connection")
	}
	if ConfigFile.TableName == "" {
		missing = append(missing, "no tableName in config")
	}
	if ConfigFile.DescriptionVar == "" {
		missing = append(missing, "no descriptionVar in config")
	}
	if ConfigFile.Resultnb == 0 {
		missing = append(missing, "no Resultnb in config")
	}
	if ConfigFile.Fuzzywuzzy.Useit {
		algocount++
		if ConfigFile.Fuzzywuzzy.Maxscore < ConfigFile.Fuzzywuzzy.Minscore {
			missing = append(missing, "Fuzzywuzzy maxscore cannot be higher than Fuzzywuzzy minscore")
		}
		if ConfigFile.Fuzzywuzzy.Maxscore > 100 {
			missing = append(missing, "Fuzzywuzzy maxscore cannot be higher than 100")
		}
		if ConfigFile.Fuzzywuzzy.Minscore < 0 {
			missing = append(missing, "Fuzzywuzzy minscore cannot be lower than 0")
		}
	}
	if ConfigFile.Rareword.Useit {
		algocount++
		if ConfigFile.Rareword.Minlength < 0 {
			missing = append(missing, "Rareword minlength cannot be lower than 0")
		}
		if ConfigFile.Rareword.Percent <= 0 || ConfigFile.Rareword.Percent >= 1{
			missing = append(missing, "Rareword percent have to be between 0 and 1")
		}
	}
	if ConfigFile.Dictionnary.Useit {
		algocount++
		if ConfigFile.Dictionnary.Databasename == "" {
			missing = append(missing, "Dictionnary Databasename is not set")
		}
		if ConfigFile.Dictionnary.Tablename == "" {
			missing = append(missing, "Dictionnary Databasename is not set")
		}
	}
	if ConfigFile.Fuzzydico.Useit {
		algocount++
	}
	if algocount != 1 {
		missing = append(missing, "1 algorithm has to be used. " + strconv.Itoa(algocount) + " are used")
	}
	if len(missing) != 0 {
		strerror := strconv.Itoa(len(missing)) + " error(s) detected in config file:\n"
		for _, missunit := range missing {
			strerror += "\t" + missunit + "\n"
		}
		return errors.New(strerror)
	}
	return nil
}

func usedalgo(ConfigFile *semanticstruct.SConfig) int {
	count := 0
	if ConfigFile.Fuzzywuzzy.Useit {
		count++
	}
	if ConfigFile.Rareword.Useit {
		count++
	}
	if ConfigFile.Dictionnary.Useit {
		count++
	}
	return count
}