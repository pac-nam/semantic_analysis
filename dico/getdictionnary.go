package dico

import (
	"fmt"
	"kdatapack/utils/log"
	"semanticContent/semanticstruct"
)

func getDictionnary(ctx *semanticstruct.SContext) error {
	ctx.File.Connect.DatabaseName = ctx.File.Dictionnary.Databasename
	if err := ctx.File.Connect.NewFlag(); err != nil {
		return err
	}
	defer ctx.File.Connect.Close()
	if _, err := ctx.File.Connect.Db.Exec("DESCRIBE " + ctx.File.Dictionnary.Tablename); err != nil {
		// MySQL error 1146 is "table does not exist"
		if fmt.Sprint(err) !=  "Error 1146: Table '" + ctx.File.Dictionnary.Databasename + "." + ctx.File.Dictionnary.Tablename + "' doesn't exist" {
			return err
		}
		log.Warning.Println("Creating", ctx.File.Dictionnary.Tablename ,"table")
		if err := createDictionnary(ctx); err != nil {
			return err
		}
		if _, err := ctx.File.Connect.Db.Exec("DESCRIBE " + ctx.File.Dictionnary.Tablename); err != nil {
			return err
		}
	} else {
		rows, err := ctx.File.Connect.Db.Query("SELECT word, word2, affinity FROM " + ctx.File.Dictionnary.Tablename + ";")
		if err != nil {
			return err
		}
		var word, word2 string
		var affinity float32
		for rows.Next() {
			err := rows.Scan(&word, &word2, &affinity)
			if err != nil {
				return err
			}
			tmp, exist := ctx.Dictionnary[word]
			if !exist {
				tmp.Affinities = make(map[string]float32)
			}
			tmp.Affinities[word2] = affinity
			ctx.Dictionnary[word] = tmp
		}
		if err := rows.Err(); err != nil {
			return err
		}
	}
	return nil
}