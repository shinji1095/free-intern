package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	_ "github.com/go-sql-driver/mysql"
)

type Recipe struct {
	Id int `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Body
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Body struct {
	Title       string `json:"title" gorm:"column:title"`
	MakingTime  string `json:"making_time"`
	Serves      string `json:"serves" gorm:"column:serves"`
	Ingredients string `json:"ingredients" gorm:"column:ingredients"`
	Cost        int    `json:"cost" gorm:"column:cost"`
}

func main() {
	// db := sqlConnect()
	// defer db.Close()
	e := echo.New()
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Fprintf(os.Stderr, "Request: %v\n", string(reqBody))
	}))
	Router(e)
	e.Logger.Fatal(e.Start(":8080"))
}

func Router(e *echo.Echo) {
	e.POST("/recipes", postRecipes)
	e.GET("/recipes", getAllRecipes)
	e.GET("/recipes/:id", getRecipe)
	e.PATCH("recipes/:id", patchRecipe)
	e.DELETE("/recipes/:id", deleteRecipe)
}

type FieldsToReplace struct {
	Replace1 string
}

func postRecipes(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()
	recipe := new(Recipe)

	// カラムがすべて存在するか確認
	if err := c.Bind(&recipe); err != nil {
		message := `{
			"message": "Recipe creation failed!"
			"required": "title, making_time, serves, ingredients, cost"
		   }`
		fmt.Print(err)
		return c.JSON(http.StatusOK, message)
	}

	fmt.Print("post Body: ", recipe, "\n")
	db.Create(&recipe)
	db.Last(&recipe)
	fmt.Print(recipe)
	jsonData, _ := json.Marshal(recipe)
	message := `{
		"message": "Recipe successfully created!",
		"recipe": [` + string(jsonData) + `]}`
	return c.JSON(http.StatusOK, message)
}

func getAllRecipes(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	recipes := []Recipe{}
	db.Find(&recipes)
	fmt.Println(recipes)
	jsonData, _ := json.Marshal(recipes)
	message := `{
		"recipes": 
			` + string(jsonData) + `
		}`
	return c.JSON(http.StatusOK, message)
}

func getRecipe(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()

	id := c.Param("id")
	fmt.Println("id is ", id)

	recipe := []Recipe{}
	db.Find(&recipe, "id=?", id)
	jsonData, _ := json.Marshal(recipe)
	message := `{
		"message": "Recipe details by id",
		"recipe": ` + string(jsonData) + `
	  }`
	return c.JSON(http.StatusOK, message)
}

func checkEmpty(recipe Recipe) map[string]interface{} {
	columns := make(map[string]interface{})
	if recipe.Cost != 0 {
		columns["cost"] = recipe.Cost
	}

	if recipe.Ingredients != "" {
		columns["ingredients"] = recipe.Ingredients
	}

	if recipe.MakingTime != "" {
		columns["making_time"] = recipe.MakingTime
	}

	if recipe.Serves != "" {
		columns["serves"] = recipe.Serves
	}

	if recipe.Title != "" {
		columns["title"] = recipe.Title
	}

	return columns
}

func patchRecipe(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()
	var id = c.Param("id")
	fmt.Print("id: ", id, "\n")

	// レシピidが取得できていない場合の判定
	if id == "" {
		id = "1"
	}
	var recipe Recipe

	c.Bind(&recipe)
	int_id, _ := strconv.Atoi(id)
	recipe.Id = int_id
	var columns = checkEmpty(recipe)
	fmt.Print(columns)
	db.Model(&recipe).Update(columns)
	// db.Save(&recipe)
	// db.Select("").Where("id=?", id).First(&patchResult)
	jsonData, _ := json.Marshal(recipe)
	fmt.Print("recipe: ", recipe, "\n")

	message := `{
		"message": "Recipe successfully updated!",
		"recipe": 	` + string(jsonData) + `
	  }`
	return c.JSON(http.StatusOK, message)
}

func deleteRecipe(c echo.Context) error {
	db := sqlConnect()
	defer db.Close()
	id := c.Param(("id"))
	fmt.Print("id", id, "\n")
	var recipe Recipe
	if err := db.Where("id=?", id).First(&recipe).Error; err != nil {
		message := `{ "message":"No Recipe found" }`
		fmt.Print(err)
		return c.JSON(http.StatusOK, message)
	}
	db.Delete(recipe, id)
	message := ` {  "message": "Recipe successfully removed!" }`
	return c.JSON(http.StatusOK, message)
}

func sqlConnect() (database *gorm.DB) {
	var DBMS string
	var USER string
	var PASS string
	var PROTOCOL string
	var DBNAME string
	var URL string

	switch os.Getenv("env") {
	case "production":
		log.Print("access as production")
		// DBMS = "postgres"
		// USER = "bnlpapzoyefidn"
		// PASS = "9cd9e4ff62abb18c514b75c532cafb621564c316c6d56a278561e5438f73d1ca"
		// PROTOCOL = "ec2-52-86-25-51.compute-1.amazonaws.com:5432"
		// DBNAME = "dirmm48brfasp"
		URL = os.Getenv("DATABASE_URL")
	default:
		log.Print("access as development")
		DBMS = "postgres"
		USER = "go"
		PASS = "pass"
		PROTOCOL = "tcp(db:3306)"
		DBNAME = "godb"
		URL = USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	}
	URL = os.Getenv("DATABASE_URL")

	count := 0
	db, err := gorm.Open(DBMS, URL)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				fmt.Println("DB接続失敗")
				panic(err)
			}
			db, err = gorm.Open(DBMS, URL)
		}
	}
	fmt.Println("DB接続成功")

	return db
}
