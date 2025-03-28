package crates

import (
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// uploadFile uploads a file to GridFS.
func uploadFile(bucket *gridfs.Bucket, filename string) error {
	// Check if the file exists

	filepath := "./uploads/" + filename

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("arquivo '%s' não existe", filename)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir stream de upload: %w", err)
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, file)
	if err != nil {
		return fmt.Errorf("erro ao copiar arquivo para o GridFS: %w", err)
	}

	fmt.Printf("Arquivo '%s' enviado com sucesso para o GridFS!\n", filename)
	return nil
}
