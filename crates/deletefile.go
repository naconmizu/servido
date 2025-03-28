package crates

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// deleteFile deletes a file from GridFS by name.
func deleteFile(bucket *gridfs.Bucket, filename string) error {
	fileStream, err := bucket.OpenDownloadStreamByName(filename)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo no GridFS: %w", err)
	}
	defer fileStream.Close()

	fileID := fileStream.GetFile().ID
	if err := bucket.Delete(fileID); err != nil {
		return fmt.Errorf("erro ao excluir arquivo no GridFS: %w", err)
	}

	fmt.Printf("Arquivo '%s' excluído com sucesso!\n", filename)
	return nil
}
