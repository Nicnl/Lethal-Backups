package main

import (
	"fmt"
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

	// 2. Wait forever
	select {}
}
