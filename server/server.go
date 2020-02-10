package server

import "github.com/gin-gonic/gin"
import "net/http"
import "fmt"
import "time"
import "github.com/andrewarrow/feedback/util"
import "github.com/andrewarrow/feedback/persist"
import "github.com/andrewarrow/feedback/controllers"

func Serve(port string) {
	prefix := util.AllConfig.Path.Prefix

	controllers.Db = persist.Connection()
	router := gin.Default()

	router.Static("/assets", prefix+"assets")
	router.GET("/", controllers.WelcomeIndex)
	inboxes := router.Group("/inboxes")
	inboxes.GET("/", controllers.InboxesIndex)
	users := router.Group("/users")
	users.GET("/", controllers.UsersIndex)
	user := router.Group("/user")
	user.GET("/:id", controllers.UsersShow)
	sessions := router.Group("/sessions")
	sessions.GET("/new", controllers.SessionsNew)
	sessions.POST("/", controllers.SessionsCreate)
	sessions.POST("/destroy", controllers.SessionsDestroy)

	admin := router.Group("/admin")
	users = admin.Group("/users")
	users.GET("/", controllers.AdminUsersIndex)
	user = admin.Group("/user")
	user.GET("/:id", controllers.AdminUsersShow)

	for _, path := range util.AllConfig.Directories.Active {
		router.GET("/"+path+"/", controllers.DirectoriesIndex)
		router.GET("/"+path+"/:name", controllers.DirectoriesDownload)
		router.GET("/"+path+"/:name/", controllers.DirectoriesNameIndex)
		router.GET("/"+path+"/:name/:extra", controllers.DirectoriesDownloadExtra)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	AddTemplates(router, prefix)
	go router.Run(fmt.Sprintf(":%s", port))

	for {
		time.Sleep(time.Second)
	}

}
