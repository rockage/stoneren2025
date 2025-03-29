package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
)

func Users(ctx iris.Context) { // 主题详情：显示所有回帖
	var query, where, sort, uid, username string
	var page, offset int

	page, _ = strconv.Atoi(ctx.URLParamDefault("page", "1"))
	limit := 25
	offset = (page - 1) * limit
	uid = ctx.URLParam("uid")
	sort = ctx.URLParam("sort")
	username = ctx.URLParam("username")

	if uid == "0" {
		where = ""
	} else {
		// 带Where的查询
		where = fmt.Sprintf("WHERE uid = %s", uid)
	}

	if username != "" {
		where = fmt.Sprintf("WHERE username LIKE '%%%s%%'", username)
	}

	switch sort {
	case "0":
	}

	switch sort {
	case "0":
		sort = ""
	case "1":
		sort = "ORDER BY regdate ASC,uid ASC"
	case "2":
		sort = "ORDER BY regdate DESC,uid ASC"
	case "3":
		where = " WHERE threads > 0 "
		sort = "ORDER BY threads DESC,uid ASC"
	case "4":
		where = " WHERE threads > 0 "
		sort = "ORDER BY threads ASC,uid ASC"
	case "5":
		where = " WHERE POSTS > 0 "
		sort = "ORDER BY posts DESC,uid ASC"
	case "6":
		where = " WHERE POSTS > 0 "
		sort = "ORDER BY posts ASC,uid ASC"
	case "7":
		sort = "ORDER BY uid ASC"
	default:
		sort = ""
	}

	query = `SELECT uid,email,username,regdate,threads,posts,credits,signature FROM pre_common_member  
				##WHERE## ##SORT## LIMIT ##L## OFFSET ##O##`

	// 定义正则表达式匹配 WHERE、SORT、L、O
	whereRegex := regexp.MustCompile(`##WHERE##`)
	sortRegex := regexp.MustCompile(`##SORT##`)
	limitRegex := regexp.MustCompile(`##L##`)
	offsetRegex := regexp.MustCompile(`##O##`)

	query = whereRegex.ReplaceAllString(query, where)
	query = sortRegex.ReplaceAllString(query, sort)
	query = limitRegex.ReplaceAllString(query, strconv.Itoa(limit))
	query = offsetRegex.ReplaceAllString(query, strconv.Itoa(offset))

	datas, err := QueryDB(query)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	for k, v := range datas {
		ts := fmt.Sprintf("%v", v["regdate"])
		datas[k]["regdate"] = timestamp(ts)
	}
	//SELECT uid,email,username,regdate,credits,threads FROM pre_common_member
	ctx.JSON(datas)
}

func Comments(ctx iris.Context) { // 主题详情：显示所有回帖
	page, _ := strconv.Atoi(ctx.URLParamDefault("page", "1"))
	tid := ctx.URLParam("tid")
	limit := 10
	offset := (page - 1) * limit

	query := fmt.Sprintf(`SELECT pid,author,authorid,dateline,message,position FROM pre_forum_post WHERE tid = %s LIMIT %d OFFSET %d;`, tid, limit, offset)

	comments, err := QueryDB(query)
	if err != nil {

		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	for k, v := range comments {
		ts := fmt.Sprintf("%v", v["dateline"])
		comments[k]["dateline"] = timestamp(ts)
	}

	ctx.JSON(comments)
}

func Threads(ctx iris.Context) { // 主题列表
	var query, where, sort, uid, keyword string

	page, _ := strconv.Atoi(ctx.URLParamDefault("page", "1"))
	limit := 10
	offset := (page - 1) * limit
	fid := ctx.URLParam("fid")
	sort = ctx.URLParam("sort")
	uid = ctx.URLParam("uid")
	keyword = ctx.URLParam("keyword")

	where = ""
	if fid != "0" && fid != "" {
		where = fmt.Sprintf("WHERE t.fid = %s", fid)
	}

	if uid != "0" && uid != "" {
		where = fmt.Sprintf("WHERE t.authorid = %s", uid)
	}

	if keyword != "0" && keyword != "" {
		where = "WHERE t.subject like '%" + keyword + "%'"
	}

	if sort == "0" {
		sort = "ASC"
	} else {
		// tid 排序
		sort = "DESC"
	}

	query = `SELECT t.tid, t.fid, t.authorid, t.author, t.subject, t.dateline, t.lastpost, t.lastposter, t.views, t.replies, 
				COALESCE(p.message, '') AS lastmessage 
				FROM pre_forum_thread AS t	
				LEFT JOIN pre_forum_post AS p ON t.tid = p.tid AND p.position = 1  
				##WHERE## ORDER BY t.tid ##SORT## LIMIT ##L## OFFSET ##O##`

	// 定义正则表达式匹配 WHERE、SORT、L、O
	whereRegex := regexp.MustCompile(`##WHERE##`)
	sortRegex := regexp.MustCompile(`##SORT##`)
	limitRegex := regexp.MustCompile(`##L##`)
	offsetRegex := regexp.MustCompile(`##O##`)

	query = whereRegex.ReplaceAllString(query, where)
	query = sortRegex.ReplaceAllString(query, sort)
	query = limitRegex.ReplaceAllString(query, strconv.Itoa(limit))
	query = offsetRegex.ReplaceAllString(query, strconv.Itoa(offset))

	threads, err := QueryDB(query)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	for k, v := range threads {
		ts := fmt.Sprintf("%v", v["dateline"])
		threads[k]["dateline"] = timestamp(ts)
		laststr := fmt.Sprintf("%s", v["lastmessage"])

		// 确保字符串长度足够
		runes := []rune(laststr)

		if len(runes) > 80 {
			runes = runes[:80] // 取前 80 个字符
			threads[k]["lastmessage"] = string(runes) + "……"
		} else {
			threads[k]["lastmessage"] = laststr
		}

	}

	ctx.JSON(threads)
}

func timestamp(timestamp string) string {

	num, _ := strconv.ParseInt(timestamp, 10, 64)

	t := time.Unix(num, 0)
	// 格式化为 "YYYY-MM-DD HH:MM:SS"（2006-01-02 15:04:05 是 Go 规定的时间格式）
	formattedTime := t.Format("2006-01-02")

	// 输出格式化后的时间
	return formattedTime

}
