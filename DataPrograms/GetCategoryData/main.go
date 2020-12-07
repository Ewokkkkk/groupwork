package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// CategoryList は楽天レシピAPIからカテゴリ情報を取得し、格納する構造体
type CategoryList struct {
	Result struct {
		// Smallいらないかも
		// Small []struct {
		// 	CategoryName     string `json:"categoryName"`
		// 	ParentCategoryID string `json:"parentCategoryId"`
		// 	CategoryID       int    `json:"categoryId"`
		// } `json:"small"`
		Medium []struct {
			CategoryName     string `json:"categoryName"`
			ParentCategoryID string `json:"parentCategoryId"`
			CategoryID       int    `json:"categoryId"`
		} `json:"medium"`
		Large []struct {
			CategoryName string `json:"categoryName"`
			CategoryID   string `json:"categoryId"`
		} `json:"large"`
	} `json:"result"`
}

func getCategoryData() CategoryList {
	url := "https://app.rakuten.co.jp/services/api/Recipe/CategoryList/20170426"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	params := request.URL.Query()
	params.Add("applicationId", "1086382364385531386")
	request.URL.RawQuery = params.Encode()
	fmt.Println(request.URL.String())
	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var cl CategoryList
	if err := json.Unmarshal(body, &cl); err != nil {
		log.Fatal(err)
	}

	return cl
}

func insertCategoryData(cl CategoryList) {
	db, err := sql.Open("mysql", "admin:"+os.Getenv("RDS_PASS")+"@tcp(database-1.cop2pvzm3623.ap-northeast-1.rds.amazonaws.com)/groupwork_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO category_list(category_id, category_name, parent_category_id) VALUES(?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	// for _, val := range cl.Result.Small {
	// 	id := val.CategoryID
	// 	name := val.CategoryName
	// 	pid := val.ParentCategoryID

	// 	_, err := stmtInsert.Exec(id, name, pid)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}

	// 	fmt.Println(id, name, pid)
	// }
	for _, val := range cl.Result.Medium {
		id := val.CategoryID
		name := val.CategoryName
		pid := val.ParentCategoryID

		_, err := stmtInsert.Exec(id, name, pid)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(id, name, pid)
	}
	for _, val := range cl.Result.Large {
		id := val.CategoryID
		name := val.CategoryName

		_, err := stmtInsert.Exec(id, name, nil)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(id, name)
	}
}
func main() {
	categoryList := getCategoryData()
	insertCategoryData(categoryList)
}
