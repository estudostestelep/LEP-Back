package main

import (
	"lep/resource"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	resource.Inject()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/login", resource.ServersControllers.SourceAuth.ServiceLogin)
	r.Use(func(c *gin.Context) {

		if c.Request.Method == "POST" && c.Request.URL.Path == "/login" {
			c.Next()
			return
		}
		if c.Request.Method == "POST" && c.Request.URL.Path == "/user" {
			c.Next()
			return
		}

		token := resource.ServersControllers.SourceAuth.ServiceValidateTokenIn(c)

		if token {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
		}

	})

	r.POST("/logout", resource.ServersControllers.SourceAuth.ServiceLogout)
	r.POST("/checkToken", resource.ServersControllers.SourceAuth.ServiceValidateToken)

	// Rotas para Usuários
	r.GET("/user/:id", resource.ServersControllers.SourceUsers.ServiceGetUser)
	r.GET("/user/group/:id", resource.ServersControllers.SourceUsers.ServiceGetUserByGroup)
	r.POST("/user", resource.ServersControllers.SourceUsers.ServiceCreateUser)
	r.PUT("/user", resource.ServersControllers.SourceUsers.ServiceUpdateUser)
	r.DELETE("/user/:id", resource.ServersControllers.SourceUsers.ServiceDeleteUser)

	// Rotas para Produtos
	r.GET("/product/:id", resource.ServersControllers.SourceProducts.ServiceGetProduct)
	r.GET("/product/purchase/:id", resource.ServersControllers.SourceProducts.ServiceGetProductByPurchase)
	r.POST("/product", resource.ServersControllers.SourceProducts.ServiceCreateProduct)
	r.PUT("/product", resource.ServersControllers.SourceProducts.ServiceUpdateProduct)
	r.DELETE("/product/:id", resource.ServersControllers.SourceProducts.ServiceDeleteProduct)

	// Rotas para Compras
	r.GET("/purchase/:id", resource.ServersControllers.SourcePurchases.ServiceGetPurchases)
	r.GET("/purchase/group/:id", resource.ServersControllers.SourcePurchases.ServiceGetPurchasesByGroup)
	r.POST("/purchase", resource.ServersControllers.SourcePurchases.ServiceCreatePurchase)
	r.PUT("/purchase", resource.ServersControllers.SourcePurchases.ServiceUpdatePurchase)
	r.DELETE("/purchase/:id", resource.ServersControllers.SourcePurchases.ServiceDeletePurchase)

	r.Run(":8080")
}
