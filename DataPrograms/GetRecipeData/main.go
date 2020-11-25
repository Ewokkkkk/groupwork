package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type CategoryList struct {
	CategoryID   string
	OriginID     int
	CategoryName string
}
type RecipeList struct {
	Result []struct {
		FoodImageURL string `json:"foodImageUrl"`
		// RecipeDescription string   `json:"recipeDescription"`
		// RecipePublishday  string   `json:"recipePublishday"`
		// Shop              int      `json:"shop"`
		// Pickup            int      `json:"pickup"`
		RecipeID int `json:"recipeId"`
		// Nickname          string   `json:"nickname"`
		// SmallImageURL     string   `json:"smallImageUrl"`
		RecipeMaterial   []string `json:"recipeMaterial"`
		RecipeIndication string   `json:"recipeIndication"`
		RecipeCost       string   `json:"recipeCost"`
		// Rank              string   `json:"rank"`
		RecipeURL string `json:"recipeUrl"`
		// MediumImageURL    string   `json:"mediumImageUrl"`
		RecipeTitle string `json:"recipeTitle"`
		OriginID    int
		CategoryID  string
	} `json:"result"`
}

func selectCategoryList() []CategoryList {
	db, err := sql.Open("mysql", "admin:"+os.Getenv("RDS_PASS")+"@tcp(database-1.cop2pvzm3623.ap-northeast-1.rds.amazonaws.com)/groupwork_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	rows, err := db.Query("select category_id, category_name, parent_category_id from category_list")
	if err != nil {
		log.Fatal(err)
	}
	var (
		id   int
		name string
		pid  interface{}
	)
	var category []CategoryList
	c := CategoryList{}

	for rows.Next() {
		if err := rows.Scan(&id, &name, &pid); err != nil {
			log.Fatal(err)
		}
		if pid != nil {
			// parendIDがあれば（中カテゴリ）
			x := pid.([]uint8)
			c.CategoryID = string(x) + "-" + strconv.Itoa(id)
			c.OriginID = id
			c.CategoryName = name
		} else {
			// parentIDがなければ（大カテゴリ）
			c.CategoryID = strconv.Itoa(id)
			c.OriginID = id
			c.CategoryName = name
		}
		category = append(category, c)
	}
	return category
}
func getRecipiData(cl []CategoryList) []RecipeList {
	url := "https://app.rakuten.co.jp/services/api/Recipe/CategoryRanking/20170426"
	var rl []RecipeList

	for _, val := range cl {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		params := request.URL.Query()
		params.Add("applicationId", "1086382364385531386")
		params.Add("categoryId", val.CategoryID)
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
		recipe := RecipeList{}
		if err := json.Unmarshal(body, &recipe); err != nil {
			log.Fatal(err)
		}

		for i := range recipe.Result {
			recipe.Result[i].OriginID = val.OriginID
			recipe.Result[i].CategoryID = val.CategoryID
		}

		rl = append(rl, recipe)
		time.Sleep(time.Second * 1)
		// break
	}
	return rl
}

func insertRecipeData(rl []RecipeList) {
	db, err := sql.Open("mysql", "admin:"+os.Getenv("RDS_PASS")+"@tcp(database-1.cop2pvzm3623.ap-northeast-1.rds.amazonaws.com)/groupwork_db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT IGNORE INTO recipe(recipe_id, image, indication, cost, url, title) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	for _, val := range rl {
		for _, recipe := range val.Result {
			id := recipe.RecipeID
			img := recipe.FoodImageURL
			indication := recipe.RecipeIndication
			cost := recipe.RecipeCost
			url := recipe.RecipeURL
			title := recipe.RecipeTitle

			_, err := stmtInsert.Exec(id, img, indication, cost, url, title)
			if err != nil {
				panic(err.Error())
			}
			// fmt.Println(id, img, indication, cost, url, title)
		}
	}

	stmtInsertMaterial, err := db.Prepare("INSERT INTO material(material_name) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsertMaterial.Close()

	stmtInsertMR, err := db.Prepare("INSERT INTO material_recipe(material_id, recipe_id) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsertMR.Close()

	stmtInsertCR, err := db.Prepare("INSERT INTO category_recipe(category_id, recipe_id) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsertCR.Close()

	for _, val := range rl {
		for _, recipe := range val.Result {
			rid := recipe.RecipeID

			_, err = stmtInsertCR.Exec(recipe.OriginID, rid)
			if err != nil {
				panic(err.Error())
			}

			for _, m := range recipe.RecipeMaterial {

				result, err := stmtInsertMaterial.Exec(m)
				if err != nil {
					panic(err.Error())
				}
				mid, _ := result.LastInsertId()
				_, err = stmtInsertMR.Exec(mid, rid)
				if err != nil {
					panic(err.Error())
				}

				// fmt.Println(m, mid, rid)
			}
		}
	}

}

func main() {
	categoryList := selectCategoryList()
	// for _, cal := range categoryList {
	// 	fmt.Println(cal)
	// }
	recipeList := getRecipiData(categoryList)
	// for _, recipe := range recipeList {
	// 	fmt.Println(recipe.Result[0].RecipeCost)
	// }
	insertRecipeData(recipeList)
}
