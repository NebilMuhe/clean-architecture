package routing

import "github.com/gin-gonic/gin"

type Router struct {
	Path        string
	Method      string
	Handler     gin.HandlerFunc
	Middlewares []gin.HandlerFunc
}

func RegisterRoutes(group *gin.RouterGroup, routes []Router) {
	for _, route := range routes {
		var handler []gin.HandlerFunc
		handler = append(handler, route.Middlewares...)
		handler = append(handler, route.Handler)
		group.Handle(route.Method, route.Path, handler...)
	}
}
