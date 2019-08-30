package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"kdatapack/bin"
	"kdatapack/utils/log"
	"os"
	"semanticContent/dico"
	"semanticContent/fuzz"
	"semanticContent/fuzzydico"
	"semanticContent/rare"
	"semanticContent/semanticstruct"
)

func init() {
	log.InitAuto("kdata")
	fmt.Println(bin.MakeSoftwareHeader())
	if len(os.Args) != 2 {
		log.Error.Printf("Wrong number of arguments. Usage : %s <conf.json>\n", os.Args[0])
		os.Exit(1)
	}

}

func main() {
	var err error
	ctx := new(semanticstruct.SContext)
	ctx.File, err = loadConfigFile(os.Args[1])
	if err != nil {
		log.Error.Println(err)
		return
	}
	err = GetListFromSQL(ctx)
	if err != nil {
		log.Error.Println(err)
		return
	}
	switch {
	case ctx.File.Fuzzywuzzy.Useit:
		fuzz.Fuzzywuzzy(ctx)
	case ctx.File.Rareword.Useit:
		rare.Rareword(ctx)
	case ctx.File.Dictionnary.Useit:
		dico.Dictionnary(ctx)
	case ctx.File.Fuzzydico.Useit:
		fuzzydico.Fuzzydico(ctx)
	}
}