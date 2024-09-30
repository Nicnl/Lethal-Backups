//go:generate goversioninfo -icon=icon.ico

package main

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/browser"
	"io/fs"
	"lethal_company_save_manager/routes"
	"lethal_company_save_manager/save_historizer"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

var (
	//go:embed vue/dist
	vueApp embed.FS
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

	// 2. Start the API web server
	go func() {
		r := gin.Default()

		r.Use(cors.New(cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		}))

		r.GET("/saves", routes.GetSaves)
		r.GET("/items", routes.GetItems)
		r.GET("/item_icon/:itemName", routes.GetItemIcon)
		r.GET("/planets", routes.GetPlanets)
		r.POST("/restore/:hash", routes.PostRestoreHash)

		err := r.Run("127.0.0.1:51245")
		if err != nil {
			panic(err)
		}
	}()

	// 3. Start the Vue.js web server
	go func() {
		r := gin.Default()

		vueSubFs, err := fs.Sub(vueApp, "vue/dist")
		if err != nil {
			panic(err)
		}
		r.StaticFS("/", http.FS(vueSubFs))

		err = r.Run("127.0.0.1:51246")
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(750 * time.Millisecond)
	openbrowser("http://127.0.0.1:51246")

	// 4. Wait forever
	select {}
}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = browser.OpenURL(url)
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
