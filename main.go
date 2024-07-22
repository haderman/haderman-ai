package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/brunotm/uulid"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const expectedHost string = "0.0.0.0:8080"

func main() {
	loadEnvVariables()

	println("**** USER_NAME ****")
	println(os.Getenv("USER_NAME"))
	println("**** *** ****")

	// generateUULID()

	router := gin.Default()

	authMiddleware := AuthMiddleware()
	router.Use((HandlerAuthMiddleWare(authMiddleware)))

	router.Use(secureHeaders)

	router.POST("/login", authMiddleware.LoginHandler)

	v1 := router.Group("/v1", authMiddleware.MiddlewareFunc())
	{
		v1.GET("/albums", getAlbums)
		v1.GET("/albums/:id", getAlbumByID)
		v1.POST("/albums", postAlbums)
		v1.POST("/images", generateImage)
	}

	router.Run(expectedHost)
}

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	test()
	c.IndentedJSON(http.StatusOK, albums)
}

// getAlbumByID responds with the album with the given ID as JSON.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"error": "album not found"})
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		fmt.Println(err)
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func test() {
	primes := [6]int{2, 3, 5, 7, 11, 13}

	for _, v := range primes {
		fmt.Println(v)
	}
}

func secureHeaders(c *gin.Context) {
	if c.Request.Host != expectedHost {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
		return
	}
	c.Header("X-Frame-Options", "DENY")
	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Referrer-Policy", "strict-origin")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	c.Next()
}

func loadEnvVariables() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func generateUULID() string {
	id, err := uulid.New()
	if err != nil {
		// handle err
	}
	fmt.Println(id.String())
	return id.String()
}
