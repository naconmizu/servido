package crates

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"log"
)

// listFiles lists all files stored in GridFS.
func listFiles(bucket *gridfs.Bucket) ([]string,error) {
	cursor, err := bucket.Find(bson.D{})
	if err != nil {
		return nil,fmt.Errorf("erro ao buscar arquivos: %w", err)
	}
	defer cursor.Close(context.Background())

	var files []string // Slice para armazenar os nomes dos arquivos
	found := false
	for cursor.Next(context.Background()) {
		found = true
		var file gridfs.File
		if err := cursor.Decode(&file); err != nil {
			return nil,fmt.Errorf("erro ao decodificar arquivo: %w", err)
		}
		files = append(files, file.Name) // Adiciona o nome do arquivo ao slice
	}

	if !found {
		fmt.Println("Nenhum arquivo encontrado no GridFS.")
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Tipo do erro: %T, Mensagem: %v", err, err)
		return nil,fmt.Errorf("erro ao buscar arquivos: %w", err)
	}

	return files,nil // Retorna o slice com os nomes dos arquivos
}

