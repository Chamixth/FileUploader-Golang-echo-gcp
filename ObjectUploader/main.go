package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
)

const (
	// Change these values to match your GCP project and bucket details
	projectID       = ""
	bucketName      = ""
	credentialsFile = ""
)

func main() {
	e := echo.New()

	// Route to handle file uploads
	e.POST("/upload", handleUpload)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}

func handleUpload(c echo.Context) error {
	// Get the file from the request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println("Failed to get file from request:", err)
		return c.String(http.StatusBadRequest, "Failed to get file from request.")
	}

	// Open the file
	src, err := fileHeader.Open()
	if err != nil {
		log.Println("Failed to open file:", err)
		return c.String(http.StatusInternalServerError, "Failed to open file.")
	}
	defer src.Close()

	// Create a GCS client
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Println("Failed to create GCS client:", err)
		return c.String(http.StatusInternalServerError, "Failed to create GCS client.")
	}
	defer client.Close()

	// Create a writer to store the uploaded file in the GCS bucket using the original filename
	wc := client.Bucket(bucketName).Object(fileHeader.Filename).NewWriter(ctx)
	if _, err := io.Copy(wc, src); err != nil {
		log.Println("Failed to store file in GCS:", err)
		return c.String(http.StatusInternalServerError, "Failed to store file in GCS.")
	}
	if err := wc.Close(); err != nil {
		log.Println("Failed to close GCS writer:", err)
		return c.String(http.StatusInternalServerError, "Failed to close GCS writer.")
	}

	return c.String(http.StatusOK, "File uploaded successfully!")
}
