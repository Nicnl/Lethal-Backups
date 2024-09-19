package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lethal_company_save_manager/routes"
	"lethal_company_save_manager/save_historizer"
	"time"
)

func main() {
	// 1. Launch the save historizer loop
	go func() {
		for {
			err := save_historizer.Loop()
			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// 2. Start the web server
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	r.GET("/saves", routes.GetSaves)
	r.GET("/items", routes.GetItems)
	r.GET("/item_icon/:itemName", routes.GetItemIcon)
	r.GET("/planets", routes.GetPlanets)

	err := r.Run("127.0.0.1:51245")
	if err != nil {
		panic(err)
	}
}
