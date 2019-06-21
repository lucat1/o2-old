package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/lucat1/o2/routes"
	"github.com/lucat1/o2/routes/git"
	"github.com/lucat1/o2/shared"
	"go.uber.org/zap"
)

func main() {
	// Early setup
	runtime.GOMAXPROCS(runtime.NumCPU())
	shared.InitializeLogger()
	shared.OpenDatabase()

	if os.Getenv("O2") != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Router creation
	gin.DisableConsoleColor()
	router := gin.New()
	if os.Getenv("O2") == "dev" {
		router.LoadHTMLGlob("views/*.tmpl")
	} else {
		t, err := loadTemplate()
		if err != nil {
			shared.GetLogger().Fatal("Could not load views templates, quitting", zap.Error(err))
		}
		router.SetHTMLTemplate(t)

		router.Use(ginzap.Ginzap(shared.GetLogger(), time.RFC3339, true))
		router.Use(ginzap.RecoveryWithZap(shared.GetLogger(), true))
	}

	// Routes
	router.Use(routes.AuthMiddleware)
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(routes.Static(Assets.Files))

	router.GET("/", routes.Index)
	router.POST("/:user", routes.Logout, routes.Login, routes.Register, routes.Create)
	router.GET("/:user", routes.Logout, routes.Login, routes.Register, routes.Create, routes.User)
	router.GET("/:user/:repo", routes.ExistsRepo(true), routes.Repo)
	router.GET("/:user/:repo/", routes.ExistsRepo(true), routes.Repo)
	router.GET("/:user/:repo/tree", routes.ExistsRepo(true), routes.Tree)
	router.GET("/:user/:repo/tree/:ref", routes.ExistsRepo(true), routes.Tree)
	router.GET("/:user/:repo/tree/:ref/*path", routes.ExistsRepo(true), routes.Tree)
	router.GET("/:user/:repo/blob/:ref/*path", routes.ExistsRepo(true), routes.Blob)
	router.GET("/:user/:repo/log", routes.ExistsRepo(true), routes.Log)
	router.GET("/:user/:repo/log/:page", routes.ExistsRepo(true), routes.Log)
	router.GET("/:user/:repo/diff/:sha", routes.ExistsRepo(true), routes.Diff)
	router.GET("/:user/:repo/settings", routes.ExistsRepo(true), routes.HasAccess([]string{"repo:settings"}), routes.Settings)

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

	shared.GetLogger().Fatal("Error while serving HTTP", zap.Error(router.Run(":80")))
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		name = strings.Replace(name, "/views/", "", 1)
		t, err = t.New(name).Parse(string(h))
		shared.GetLogger().Info("Registering template", zap.String("name", name))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
