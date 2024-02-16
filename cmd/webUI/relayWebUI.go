package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RelayInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	UDP  bool   `json:"udp"`
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/ui") })
	r.StaticFS("/ui", http.Dir("static"))

	mockData := []RelayInfo{
		{
			Id:   uuid.New().String(),
			Name: "#1",
			Src:  "1.1.1.1:53",
			Dst:  "2.2.2.2:53",
			UDP:  true,
		},
		{
			Id:   uuid.New().String(),
			Name: "#2",
			Src:  "10.10.10.10:80",
			Dst:  "20.20.20.20:80",
			UDP:  false,
		},
	}
	r.GET("api/list.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, mockData)
	})
	r.Any("api/cp", func(c*gin.Context){
		fmt.Println("clip:", c.Request.Method, c.Request)
	})

	r.Run()
}
