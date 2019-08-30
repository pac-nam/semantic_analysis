package fuzzydico

// int fortytwo()
// {
//    return 42;
// }
import "C"
import (
	"fmt"
	"github.com/charlesvdv/fuzmatch"

	//kdata "kdatapack/utils"
	"kdatapack/utils/log"
	"semanticContent/semanticstruct"
	"strings"
)

type Sscore struct {
	index int
	score float32
}

func Fuzzydico(ctx *semanticstruct.SContext) {
	ctx.Dictionnary = make(map[string]semanticstruct.SDicoword)
	if err := getDictionnary(ctx); err != nil {
		log.Error.Println(err)
		return
	}
	log.Info.Println("start looking for semantic similar on", ctx.DescLen, "descriptions")
	var tmp Sscore
	//bar := kdata.InitBar(ctx.DescLen)
	for index := range ctx.DescList {
		if ctx.DescList[index] == "" {
			continue
		}
		leaderboard := make([]Sscore, ctx.File.Resultnb)
		for index2 := range ctx.DescList {
			if strings.Compare(ctx.DescList[index], ctx.DescList[index2]) == 0 || ctx.DescList[index2] == "" {
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
		fmt.Println(index, ctx.DescList[index])
		for i := ctx.File.Resultnb; i > 0; {
			i--
			fmt.Printf("\t%.2f %s\n", leaderboard[i].score, ctx.DescList[leaderboard[i].index])
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
	result = result / float32(len(ctx.SplittedDescList[index2])*len(ctx.SplittedDescList[index1]))
	result *= float32(fuzmatch.TokenSortRatio(ctx.NormalizedDescList[index1], ctx.NormalizedDescList[index2])+
		fuzmatch.TokenSortRatio(ctx.NormalizedDescList[index2], ctx.NormalizedDescList[index1])) / 2
	return result
}
