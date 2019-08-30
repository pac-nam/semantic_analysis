package dico

import (
	"fmt"
	//kdata "kdatapack/utils"
	"kdatapack/utils/log"
	"semanticContent/semanticstruct"
	"strings"
)

type Sscore struct{
	index	int
	score	float32
}

func Dictionnary(ctx *semanticstruct.SContext) {
	ctx.Dictionnary = make(map[string]semanticstruct.SDicoword)
	if err := getDictionnary(ctx); err != nil {
		log.Error.Println(err)
		return
	}
	log.Info.Println("start looking for semantic similar")
	var tmp Sscore
	//bar := kdata.InitBar(ctx.DescLen)
	for index := range ctx.DescList {
		leaderboard := make([]Sscore, ctx.File.Resultnb)
		for index2 := range ctx.DescList {
			if strings.Compare(ctx.DescList[index], ctx.DescList[index2]) == 0 {
				continue
			}
			score := getScore(ctx, index, index2)
			tmp.index, tmp.score = index2, score
			indexleaderboard := 0
			for indexleaderboard < ctx.File.Resultnb && score >= leaderboard[indexleaderboard].score {
				if strings.Compare(ctx.DescList[tmp.index], ctx.DescList[leaderboard[indexleaderboard].index]) == 0 {
					indexleaderboard = 0
					break
				}
				indexleaderboard++
			}
			for indexleaderboard > 0 {
				indexleaderboard--
				tmp, leaderboard[indexleaderboard] = leaderboard[indexleaderboard], tmp
			}
		}
		fmt.Println(ctx.DescList[index])
		for _, leader := range leaderboard {
			fmt.Printf("\t%.2f %s\n", leader.score, ctx.DescList[leader.index])
		}
		fmt.Println()
		//bar.Add(1)
	}
	//bar.Finish()
}

func getScore(ctx *semanticstruct.SContext, index1, index2 int) float32 {
	var result float32
	if len(ctx.SplittedDescList[index2]) < 6 {
		return result
	}
	for _, word1 := range ctx.SplittedDescList[index1] {
		for _, word2 := range ctx.SplittedDescList[index2] {
			result += ctx.Dictionnary[word1].Affinities[word2]
		}
	}
	return result / float32(len(ctx.SplittedDescList[index2]) * len(ctx.SplittedDescList[index1]))
}