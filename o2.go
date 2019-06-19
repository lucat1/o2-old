package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/git/routes"
	"github.com/lucat1/git/routes/git"
	"github.com/lucat1/git/shared"
	"go.uber.org/zap"
	"runtime"
)

func main() {
	// Early setup
	runtime.GOMAXPROCS(runtime.NumCPU())
	shared.InitializeLogger()
	shared.OpenDatabase()

	// Router creation
	gin.DisableConsoleColor()
	router := gin.New()
	router.LoadHTMLGlob("views/*.tmpl")

	// Logging setup
	log := shared.GetLogger()
	//router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	//router.Use(ginzap.RecoveryWithZap(log, true))

	// Routes
	router.Use(routes.LogMiddleware)
	router.Use(routes.AuthMiddleware)
	router.Use(static.Serve("/static", static.LocalFile("static", false)))
	router.GET("/", routes.Index)
	router.POST("/:user", routes.Logout, routes.Login, routes.Register, routes.Create)
	router.GET("/:user", routes.Logout, routes.Login, routes.Register, routes.Create, routes.User)
	router.GET("/:user/:repo", routes.Repo)
	router.GET("/:user/:repo/tree", routes.Tree)
	router.GET("/:user/:repo/tree/:ref", routes.Tree)
	router.GET("/:user/:repo/tree/:ref/*path", routes.Tree)
	router.GET("/:user/:repo/blob/:ref/*path", routes.Blob)
	router.GET("/:user/:repo/log", routes.Log)

	// Git smart http protocol
	router.GET("/:user/:repo/info/refs", routes.ExistsRepo(false), git.GetInfoRefs)
	router.GET("/:user/:repo/info/refs/*path", routes.ExistsRepo(false), git.GetInfoRefs)
	router.POST("/:user/:repo/git-upload-pack", routes.ExistsRepo(false), git.ServiceRPC)
	router.POST("/:user/:repo/git-upload-pack/*path", routes.ExistsRepo(false), git.ServiceRPC)
	router.POST("/:user/:repo/git-receive-pack", routes.ExistsRepo(false), git.ServiceRPC)
	router.POST("/:user/:repo/git-receive-pack/*path", routes.ExistsRepo(false), git.ServiceRPC)
	router.GET("/:user/:repo/HEAD", routes.ExistsRepo(false), git.GetTextFile)
	router.GET("/:user/:repo/HEAD/*path", routes.ExistsRepo(false), git.GetTextFile)
	router.GET("/:user/:repo/objects/info/alternates", routes.ExistsRepo(false), git.GetTextFile)
	router.GET("/:user/:repo/objects/info/alternates/*path", routes.ExistsRepo(false), git.GetTextFile)
	router.GET("/:user/:repo/objects/info/http-alternates", routes.ExistsRepo(false), git.GetTextFile)
	router.GET("/:user/:repo/objects/info/http-alternates/*path", routes.ExistsRepo(false), git.GetTextFile)
	router.GET("/:user/:repo/objects/info/packs", routes.ExistsRepo(false), git.GetInfoPacks)
	router.GET("/:user/:repo/objects/info/packs/*path", routes.ExistsRepo(false), git.GetInfoPacks)
	//router.GET("/:user/:repo/objects/info/*path", routes.ExistsRepo(false), git.GetTextFile)
	//router.GET("/:user/:repo/objects/*path", routes.ExistsRepo(false), git.GetLooseObject)
	router.GET("/:user/:repo/objects/pack/pack-:pack", routes.ExistsRepo(false), git.GetPackOrIdx)

	router.NoRoute(routes.NotFound)
	router.NoMethod(routes.NotFound)

	log.Fatal("Error while serving HTTP", zap.Error(router.Run(":80")))
}
