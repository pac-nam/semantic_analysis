package dico

import (
	kdata "kdatapack/utils"
	"kdatapack/utils/log"
	//"os"
	"semanticContent/semanticstruct"
)

func createDictionnary(ctx *semanticstruct.SContext) error {
	if _, err := ctx.File.Connect.Db.Exec("CREATE TABLE IF NOT EXISTS `" + ctx.File.Dictionnary.Databasename + "`.`" + ctx.File.Dictionnary.Tablename + "` (`word` VARCHAR(45) NOT NULL,`word2` VARCHAR(45) NOT NULL, `affinity` FLOAT UNSIGNED NOT NULL, PRIMARY KEY (`word`, `word2`));"); err != nil {
		return err
	}
	bar := kdata.InitBar(ctx.DescLen)
	for _, desc := range ctx.SplittedDescList {
		for _, word := range desc {
			_, exist := ctx.Dictionnary[word]
			if !exist {
				ctx.Dictionnary[word] = semanticstruct.SDicoword{
					OtherAppear: getOtherWordAppear(ctx, word),
					Appear: 0,
				}
			}
			tmp := ctx.Dictionnary[word]
			tmp.Appear++
			ctx.Dictionnary[word] = tmp
		}
		bar.Add(1)
	}
	bar.Finish()
	if err := pushDataToDico(ctx); err != nil {
		return err
	}
	return nil
}

func getOtherWordAppear(ctx *semanticstruct.SContext, word string) map[string]uint32 {
	affinities := make(map[string]uint32)
	for _, desc := range ctx.SplittedDescList {
		for _, wordToVerify := range desc {
			if word == wordToVerify {
				for _, wordToAdd := range desc {
					affinities[wordToAdd]++
				}
				break
			}
		}
	}
	return affinities
}

func pushDataToDico(ctx *semanticstruct.SContext) error {
	sqlParams := make([]interface{}, 0)
	for word, dicopage := range ctx.Dictionnary {
		dicopage.Affinities = make(map[string]float32)
		for word2, appear := range dicopage.OtherAppear {
			dicopage.Affinities[word2] = float32(appear) / float32(ctx.Dictionnary[word2].Appear)
			sqlParams = append(sqlParams, word)
			sqlParams = append(sqlParams, word2)
			sqlParams = append(sqlParams, dicopage.Affinities[word2])
		}
		ctx.Dictionnary[word] = dicopage
	}
	log.Info.Println("pushing dictionnary in table:", ctx.File.Dictionnary.Tablename)
	err := kdata.MassInsertIgnore(ctx.File.Connect.Db,
		4,
		ctx.File.Dictionnary.Tablename,
		[]string{"word", "word2", "affinity"},
		"MassInsert in " + ctx.File.Dictionnary.Tablename,
		sqlParams...)
	if err != nil {
		return err
	}
	return nil
}