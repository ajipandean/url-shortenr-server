package main

import (
  "fmt"
  "time"
  "net/http"
  "math/rand"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  b62 "github.com/catinello/base62"
)

type url struct {
  ID int `json:"id" gorm:"primaryKey"`
  Clicks uint `json:"clicks"`
  URL string `json:"url"`
  Base string `json:"base"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  DeletedAt *time.Time `json:"deleted_at"`
}
func (u *url) BeforeSave(tx *gorm.DB) error {
  u.ID = rand.Intn(9999)

  if u.Base == "" {
    base := b62.Encode(u.ID)
    u.Base = base
  }
  return nil
}

const DSN string = "root:@tcp(127.0.0.1:3306)/shortenrbase?charset=utf8mb4&parseTime=true&loc=Local"

var db *gorm.DB

func main() {
  var err error
  db, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
  if err != nil {
    fmt.Println(err.Error())
  } else {
    db.AutoMigrate(&url{})
    fmt.Println("Connected to shortenrbase")
  }

  app := echo.New()
  app.Use(middleware.CORS())
  app.GET("/:base", func(context echo.Context) error {
    schema := new(url)
    base := context.Param("base")
    decoded, err := b62.Decode(base)
    if err != nil {
      return err
    }
    result := db.Where("id = ?", decoded).First(&schema)
    if result.Error != nil {
      return result.Error
    }
    return context.Redirect(http.StatusMovedPermanently, schema.URL)
  })
  app.POST("/shorten", func(context echo.Context) error {
    schema := new(url)
    if err := context.Bind(schema); err != nil {
      return err
    }
    result := db.Create(&schema)
    if result.Error != nil {
      return result.Error
    }
    return context.JSON(http.StatusOK, schema)
  })
  app.Logger.Fatal(app.Start("localhost:8000"))
}
