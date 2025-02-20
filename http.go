package loomcache

import "github.com/gin-gonic/gin"

const DefaultBasePath = "/loomcache/api"

type HttpServer struct {
	addr     string
	basePath string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		addr:     addr,
		basePath: DefaultBasePath,
	}
}

func (h *HttpServer) Run() {
	r := gin.New()
	baser := r.Group(h.basePath)
	{
		baser.GET("/get/:group/:key", func(c *gin.Context) {
			group := c.Param("group")
			key := c.Param("key")
			if group == "" || key == "" {
				c.JSON(400, gin.H{
					"error": "group or key is empty",
				})
				return
			}
			if value, err := GetGroup(group).Get(key); err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
			} else {
				c.JSON(200, gin.H{
					"value": string(value.ByteSlice()),
				})
			}

		})

	}
	r.Run(h.addr)
}
