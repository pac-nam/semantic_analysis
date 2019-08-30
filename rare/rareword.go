package rare

import (
	"fmt"
	"kdatapack/utils/log"
	"strings"
	"semanticContent/semanticstruct"
)

type descscore struct {
	index		int
	score		float64
}

func Rareword(ctx *semanticstruct.SContext) {
	formatDescAndRarewords(ctx)
	var tmpdescscore descscore
	for i1, desc := range ctx.Rareword.Descriptions {
		if len(desc.Words) != 0 {
			leaderboard := make([]descscore, ctx.File.Resultnb)
			for i2 := range ctx.Rareword.Descriptions {
				if i1 == i2 {
					continue
				}
				score := calculScore(ctx, i1, i2)
				tmpdescscore.index, tmpdescscore.score = i2, score
				index := 0
				for index < ctx.File.Resultnb && leaderboard[index].score <= tmpdescscore.score {
					if strings.Compare(ctx.DescList[tmpdescscore.index], ctx.DescList[leaderboard[index].index]) == 0 {
						index = 0
						break
					}
					index++
				}
				for index > 0 {
					index--
					tmpdescscore, leaderboard[index] = leaderboard[index], tmpdescscore
				}
			}
			fmt.Println(ctx.DescList[i1])
			for _, semantic := range leaderboard {
				if semantic.score > 0.001 {
					fmt.Printf("\t%.2f %s\n", semantic.score, ctx.DescList[semantic.index])
				}
			}
			fmt.Println()
		}
	}
}

func formatDescAndRarewords(ctx *semanticstruct.SContext) {
	ctx.Rareword.Descriptions = make([]semanticstruct.Sdesc, ctx.DescLen)
	ctx.Rareword.Rarewords = make(map[string]float64)
	for index, desc := range ctx.NormalizedDescList {
		ctx.Rareword.Descriptions[index].Desc = desc
		tmpwordtab := strings.Split(desc, " ")
		ctx.Rareword.Descriptions[index].Words = make(map[string]struct{})
		for _, word := range tmpwordtab {
			if len(word) >= ctx.File.Rareword.Minlength {
				ctx.Rareword.Descriptions[index].Words[word] = struct{}{}
				ctx.Rareword.Rarewords[word]++
			}
		}
	}
	floatDescLen := float64(ctx.DescLen)
	for word, appear := range ctx.Rareword.Rarewords {
		if appear > floatDescLen * ctx.File.Rareword.Percent || appear < 2 {
			delete(ctx.Rareword.Rarewords, word)
		}
	}
	for index := range ctx.Rareword.Descriptions {
		for word := range ctx.Rareword.Descriptions[index].Words {
			_, exist := ctx.Rareword.Rarewords[word]
			if !exist {
				delete(ctx.Rareword.Descriptions[index].Words, word)
			}
		}
		if len(ctx.Rareword.Descriptions[index].Words) < 1 && ctx.Rareword.Descriptions[index].Desc != ""{
			log.Warning.Println("no rareword:", ctx.Rareword.Descriptions[index].Desc)
		}
	}
}

func calculScore(ctx *semanticstruct.SContext, index1, index2 int) float64 {
	var score float64
	for word := range ctx.Rareword.Descriptions[index1].Words {
		_, exist := ctx.Rareword.Descriptions[index2].Words[word]
		if exist {
			score += (1 / ctx.Rareword.Rarewords[word] ) * float64(len(word))
		}
	}
	return score
}