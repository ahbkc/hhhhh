package core

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"structs"
	"utils"
)

//index page
func IndexGet(w http.ResponseWriter, r *http.Request) {
	//PRODUCT  从打包的静态文件中获取文件
	if utils.ConfigureMap["BASE"]["ENVIRONMENT"] == "PRODUCT" {
		t, err = template.New("index").Parse(utils.ReadHTMLFileToString(utils.HtmlPath + "index.html"))
	} else {
		//读取.html文件  DEVELOP
		path := filepath.Join(utils.Dir, "/src/resource", utils.HtmlPath + "index.html")
		t, err = template.ParseFiles(path)
	}
	utils.CheckErr(err)
	data := utils.GetCommonParamMap()

	var page string
	paramVal, err := url.ParseQuery(r.URL.RawQuery)
	utils.CheckErr(err)

	if len(paramVal) <= 0 || len(paramVal["page"]) < 0 {
		page = "1"
	}

	if page == "" {
		page = paramVal["page"][0]
	}

	//Regular expression judgment
	rege, err := regexp.Compile("[0-9]+")
	utils.CheckErr(err)
	if !rege.Match([]byte(page)) {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: utils.LanguageMap["GENERAL"]["WrongInputValue"]})
		return
	}

	//获取文章列表集合
	db, err := gorm.Open(utils.ConfigureMap["DATABASE"]["dialect"], utils.Dir+utils.ConfigureMap["DATABASE"]["db_path"])
	utils.CheckErr(err)
	db.SingularTable(true)
	defer db.Close()

	//get the default paging parameters
	limit := utils.ConfigureMap["BASE"]["LIMIT"]

	pageIntVal, err := strconv.Atoi(page)
	utils.CheckErr(err)

	limitIntVal, err := strconv.Atoi(limit)
	utils.CheckErr(err)

	var articles []structs.Article
	err = db.Table("article").Where("state = ?", 0).Limit(limit).Offset((pageIntVal - 1) * limitIntVal).Find(&articles).Error
	utils.CheckErr(err)

	//count number
	var count int
	db.Table("article").Count(&count)
	//utils.CheckErr(err)
	var next int = 0
	if count%limitIntVal > 0 {
		next = 1
	}

	data["PAGE_Count"] = strconv.Itoa((count/limitIntVal)+next) + "0"
	data["PAGE_Curr"] = page

	data["List"] = articles //数据
	data["Title"] = "blob 首页"
	t.Execute(w, data)
}

//detail page
func DetailPageGet(w http.ResponseWriter, r *http.Request) {
	if utils.ConfigureMap["BASE"]["ENVIRONMENT"] == "PRODUCT" {
		t, err = template.New("detail").Parse(utils.ReadHTMLFileToString(utils.HtmlPath + "detail.html"))
	} else {
		//读取.html文件  DEVELOP
		path := filepath.Join(utils.Dir, "/src/resource", utils.HtmlPath + "detail.html")
		t, err = template.ParseFiles(path)
	}
	utils.CheckErr(err)
	data := utils.GetCommonParamMap()
	vars := mux.Vars(r)
	if len(vars) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: utils.LanguageMap["API_MESSAGE"]["PARAMETERS_CANNOT_BE_EMPTY"]})
		return
	}

	var page string

	paramVal, err := url.ParseQuery(r.URL.RawQuery)
	utils.CheckErr(err)

	if len(paramVal) <= 0 || len(paramVal["page"]) < 0 {
		page = "1"
	}

	if page == "" {
		page = paramVal["page"][0]
	}

	//Regular expression judgment
	rege, err := regexp.Compile("[0-9]+")
	utils.CheckErr(err)
	if !rege.Match([]byte(page)) {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: utils.LanguageMap["GENERAL"]["WrongInputValue"]})
		return
	}

	id := vars["id"]
	db, err := gorm.Open(utils.ConfigureMap["DATABASE"]["dialect"], utils.Dir+utils.ConfigureMap["DATABASE"]["db_path"])
	utils.CheckErr(err)
	defer db.Close()
	db.SingularTable(true)
	var article structs.Article
	err = db.Table("article").Where("id = ?", id).Find(&article).Error
	utils.CheckErr(err)
	data["Detail"] = article

	//get the default paging parameters
	limit := utils.ConfigureMap["BASE"]["LIMIT"]

	pageIntVal, err := strconv.Atoi(page)
	utils.CheckErr(err)

	limitIntVal, err := strconv.Atoi(limit)
	utils.CheckErr(err)

	//查询评论列表
	var comments []structs.Comment
	err = db.Table("comment").Where("relevancy_id = ?", id).Limit(limit).Offset((pageIntVal - 1) * limitIntVal).Find(&comments).Error
	//utils.CheckErr(err)
	data["Comments"] = comments

	//count number
	var count int
	db.Table("comment").Where("relevancy_id = ?", id).Count(&count)
	//utils.CheckErr(err)
	var next int = 0
	if count%limitIntVal > 0 {
		next = 1
	}

	data["PAGE_Count"] = strconv.Itoa((count/limitIntVal)+next) + "0"
	data["PAGE_Curr"] = page

	err = t.Execute(w, data)
	utils.CheckErr(err)
}

//about
func GetAboutPage(w http.ResponseWriter, r *http.Request) {
	if utils.ConfigureMap["BASE"]["ENVIRONMENT"] == "PRODUCT" {
		t, err = template.New("about").Parse(utils.ReadHTMLFileToString(utils.HtmlPath + "about.html"))
	} else {
		//读取.html文件  DEVELOP
		path := filepath.Join(utils.Dir, "/src/resource", utils.HtmlPath + "about.html")
		t, err = template.ParseFiles(path)
	}
	utils.CheckErr(err)

	var page string       //default 1
	id := "adfasgasdfasd" //comment key value

	paramVal, err := url.ParseQuery(r.URL.RawQuery)
	utils.CheckErr(err)

	if len(paramVal) <= 0 || len(paramVal["page"]) < 0 {
		page = "1"
	}

	if page == "" {
		page = paramVal["page"][0]
	}

	//Regular expression judgment
	rege, err := regexp.Compile("[0-9]+")
	utils.CheckErr(err)
	if !rege.Match([]byte(page)) {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: utils.LanguageMap["GENERAL"]["WrongInputValue"]})
		return
	}

	db, err := gorm.Open(utils.ConfigureMap["DATABASE"]["dialect"], utils.Dir+utils.ConfigureMap["DATABASE"]["db_path"])
	utils.CheckErr(err)
	defer db.Close()
	db.SingularTable(true)
	data := utils.GetCommonParamMap()

	//get the default paging parameters
	limit := utils.ConfigureMap["BASE"]["LIMIT"]

	pageIntVal, err := strconv.Atoi(page)
	utils.CheckErr(err)

	limitIntVal, err := strconv.Atoi(limit)
	utils.CheckErr(err)

	//查询评论列表
	var comments []structs.Comment
	err = db.Table("comment").Where("relevancy_id = ?", id).Limit(limit).Offset((pageIntVal - 1) * limitIntVal).Find(&comments).Error

	//count number
	var count int
	db.Table("comment").Where("relevancy_id = ?", id).Count(&count)
	//utils.CheckErr(err)
	var next int = 0
	if count%limitIntVal > 0 {
		next = 1
	}
	data["Comments"] = comments
	data["RelevancyId"] = id
	data["PAGE_Count"] = strconv.Itoa((count/limitIntVal)+next) + "0"
	data["PAGE_Curr"] = page

	t.Execute(w, data)
}
