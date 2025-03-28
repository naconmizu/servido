package crates

import (
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// downloadFile downloads a file from GridFS and saves it locally.
func downloadFile(bucket *gridfs.Bucket, filename, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo local: %w", err)
	}
	defer outFile.Close()

	downloadStream, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo no GridFS: %w", err)
	}
	defer downloadStream.Close()

	if _, err := io.Copy(outFile, downloadStream); err != nil {
		return fmt.Errorf("erro ao copiar arquivo do GridFS: %w", err)
	}

	fmt.Printf("Arquivo '%s' baixado com sucesso para '%s'!\n", filename, outputPath)
	return nil
}

