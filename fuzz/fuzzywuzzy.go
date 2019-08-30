package fuzz

import (
	"fmt"
	"github.com/charlesvdv/fuzmatch"
	"semanticContent/semanticstruct"
	"strings"
)

type	Ssemantic	struct {
	score		int
	index		int
}

func Fuzzywuzzy(ctx *semanticstruct.SContext) {
	tmpsemantic := Ssemantic{
		score:0,
		index:0}
	for i1, desc := range ctx.NormalizedDescList {
		semantics := make([]Ssemantic, ctx.File.Resultnb)
		for i2, desc2 := range ctx.NormalizedDescList {
			score := fuzmatch.TokenSortRatio(desc, desc2)
			if score >= ctx.File.Fuzzywuzzy.Minscore && score <= ctx.File.Fuzzywuzzy.Maxscore && i1 != i2{
				tmpsemantic.index, tmpsemantic.score = i2, score
				index := 0
				for index < ctx.File.Resultnb && semantics[index].score <= tmpsemantic.score {
					if strings.Compare(ctx.DescList[tmpsemantic.index], ctx.DescList[semantics[index].index]) == 0 {
						index = 0
						break
					}
					index++
				}
				for index > 0 {
					index--
					tmpsemantic, semantics[index] = semantics[index], tmpsemantic
				}
			}
		}
		fmt.Println(desc)
		for _, semantic := range semantics {
			if semantic.score >= ctx.File.Fuzzywuzzy.Minscore {
				fmt.Println("\t", semantic.score, ctx.DescList[semantic.index])
			}
		}
		fmt.Println()
	}
}