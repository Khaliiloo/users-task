package routes

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"users-task/configs"
	"users-task/controllers"
	"users-task/helpers"
)

var staticRoutes = map[string]int{}

func UserRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.Use(func(c *gin.Context) {
		// running logging inside goroutine to speed up the api response
		go func() {
			if configs.Logger == nil {
				configs.Logger = configs.NewLogger("log.txt")
			}
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ := io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)

			headerMap := map[string]string{}
			for k, v := range c.Request.Header {
				value := ""
				for _, val := range v {
					value += val + ", "
				}
				headerMap[k] = value
			}

			requestInfo := map[string]string{
				"method":    c.Request.Method,
				"url":       c.Request.URL.String(),
				"header":    helpers.ConvertMapToJson(c.Request.Header),
				"body":      string(body),
				"form":      helpers.ConvertMapToJson(c.Request.Form),
				"post_form": helpers.ConvertMapToJson(c.Request.PostForm),
			}
			configs.Logger.Infoln(requestInfo)

		}()
		c.Next()
		go func() {
			if configs.Logger == nil {
				configs.Logger = configs.NewLogger("log.txt")
			}
			blw := &configs.BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw

			configs.Logger.Infoln(blw.Body.String())
		}()
	})
	v1.POST("/user", controllers.CreateUser())
	v1.GET("/user/:id", controllers.GetUser())
	v1.PUT("/user/:id", controllers.EditAUser())
	v1.DELETE("/user/:id", controllers.DeleteAUser())
	v1.GET("/users", controllers.GetAllUsers())
	v1.POST("/user/:id/file", controllers.AddFile())
	v1.GET("/user/:id/files",
		func(c *gin.Context) {
			id := c.Param("id")
			if _, ok := staticRoutes["/files/"+id]; !ok {
				v1.StaticFS("/files/"+id, http.Dir("./files/"+id+"/"))
				staticRoutes["/files/"+id] = 1
			}
			files := helpers.ListFilesInDir("./files/" + id + "/")
			type file struct {
				Name string
				Link string
			}
			filesResponse := make([]file, len(files))
			for i, f := range files {
				filesResponse[i] = file{f, "localhost" + configs.ServerPort + "/api/v1/files/" + id + "/" + f}
			}
			c.JSON(200, filesResponse)
		})

}
