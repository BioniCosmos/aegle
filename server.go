package main

import "github.com/gin-gonic/gin"

func serveSubscriptions(subLinks map[string]string, serverListen, serverPath string) {
	r := gin.Default()

	r.GET(serverPath+":id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(200, "%s", subLinks[id])
	})

	r.Run(serverListen)
}
