package http

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Module interface {
	RegisterRoutes(*gin.Engine)
}

type Server struct {
	router  *gin.Engine
	modules []Module
}

func NewServer() *Server {
	return &Server{
		router:  gin.New(),
		modules: make([]Module, 0),
	}
}

func (s *Server) AddModule(module Module) {
	s.modules = append(s.modules, module)
}

func (s *Server) SetupRoutes() {
	s.router.HandleMethodNotAllowed = true

	// Swagger
	// available only in development mode
	// http://localhost:8080/swagger/index.html
	if gin.Mode() != gin.ReleaseMode {
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	for _, module := range s.modules {
		module.RegisterRoutes(s.router)
	}
}

func (s *Server) AddMiddlewares(middlewares ...gin.HandlerFunc) {
	s.router.Use(middlewares...)
}

func (s *Server) AddModules(modules ...Module) {
	s.modules = append(s.modules, modules...)
}

func (s *Server) Start(addr string) error {
	s.SetupRoutes()
	return s.router.Run(addr)
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
