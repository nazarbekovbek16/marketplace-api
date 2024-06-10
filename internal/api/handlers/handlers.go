package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"strings"
)

type Handlers struct {
	AuthHandler        *AuthHandler
	DistributorHandler *DistributorHandler
	StoreHandler       *StoreHandler
	AdminHandler       *AdminHandler
}

func NewHandlers(authHandler *AuthHandler, distributorHandler *DistributorHandler, storeHandler *StoreHandler, adminHandler *AdminHandler) *Handlers {
	return &Handlers{AuthHandler: authHandler, DistributorHandler: distributorHandler, StoreHandler: storeHandler, AdminHandler: adminHandler}
}

func (h *Handlers) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	uniqueId := uuid.New()

	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	fileExt := filepath.Ext(file.Filename)
	image := fmt.Sprintf("%s%s", filename, fileExt)

	err = c.SaveUploadedFile(file, fmt.Sprintf("./images/%s", image))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	imageUrl := fmt.Sprintf("http://localhost:4000/api/images/%s", image)
	c.JSON(http.StatusOK, gin.H{"image_url": imageUrl})
}

func (h *Handlers) UploadImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["images[]"]
	if len(files) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	var ImageUrls []string
	for _, file := range files {
		uniqueId := uuid.New()

		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := filepath.Ext(file.Filename)
		image := fmt.Sprintf("%s%s", filename, fileExt)

		err := c.SaveUploadedFile(file, fmt.Sprintf("./images/%s", image))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}
		imageUrl := fmt.Sprintf("http://localhost:4000/api/images/%s", image)
		ImageUrls = append(ImageUrls, imageUrl)
	}

	c.JSON(http.StatusOK, gin.H{"image_urls": ImageUrls})
}
