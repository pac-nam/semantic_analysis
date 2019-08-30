package main

import (
	"database/sql"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"kdatapack/utils/log"
	"regexp"
	"semanticContent/semanticstruct"
	"strconv"
	"strings"
	"unicode"
)

func GetListFromSQL(ctx *semanticstruct.SContext) error {

	if err := ctx.File.Connect.NewFlag(); err != nil {
		log.Error.Println("Connection error")
		return err
	}
	defer ctx.File.Connect.Close()
	limit := ""
	if ctx.File.Limit == 0 {
		row := ctx.File.Connect.Db.QueryRow("SELECT COUNT(*) FROM " + ctx.File.TableName + ";")
		err := row.Scan(&ctx.DescLen)
		if err != nil {
			log.Error.Println("Count error")
			return err
		}
	} else {
		ctx.DescLen = ctx.File.Limit
		limit = " LIMIT " + strconv.Itoa(ctx.File.Limit)
	}
	ctx.DescList = make([]string, ctx.DescLen)
	ctx.NormalizedDescList = make([]string, ctx.DescLen)
	ctx.SplittedDescList = make([][]string, ctx.DescLen)
	rows, err := ctx.File.Connect.Db.Query("SELECT " + ctx.File.DescriptionVar + " from " + ctx.File.TableName + limit + ";" )
	if err != nil {
		log.Error.Println("Query error")
		return err
	}

	var description sql.NullString
	index := 0
	for rows.Next() {
		err := rows.Scan(&description)
		if err != nil {
			log.Error.Println("row scan error")
			return err
		}
		if description.Valid && description.String != ""{
			ctx.DescList[index] = description.String
			ctx.NormalizedDescList[index] = normalize(description.String)
			ctx.SplittedDescList[index] = strings.Split(ctx.NormalizedDescList[index], " ")
		}
		index++
	}
	if err := rows.Err(); err != nil || index != ctx.DescLen {
		log.Error.Println("row Next error", err)
		return err
	}
	return nil
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func normalize(unicodestr string) string {
	normalbytes := make([]byte, len(unicodestr))
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	_, _, e := t.Transform(normalbytes, []byte(unicodestr), true)
	if e != nil {
		panic(e)
	}
	normalstr := strings.ToLower(string(normalbytes))
	reg := regexp.MustCompile("[^a-z ]+")
	normalstr = reg.ReplaceAllString(normalstr, " ")
	reg = regexp.MustCompile(" +")
	normalstr = reg.ReplaceAllString(normalstr, " ")
	return strings.Trim(normalstr, " ")
}