package crates

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Grider handles the command-line interface for GridFS operations.
func Grider(args ...string) []string {
	if len(args) == 0 {
		log.Fatalf("Uso: go run . <comando> [filename]\nComandos: up, down, list, delete")
	}

	command := args[0]
	var fileName string
	if command != "list" {
		if len(args) != 2 {
			log.Fatalf("Uso: go run . %s <filename>", command)
		}
		fileName = args[1]
	} else {
		if len(args) != 1 {
			log.Fatalf("Uso: go run . list")
		}
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI não está definida no arquivo .env")
	}
	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Fatal("DATABASE_NAME não está definida no arquivo .env")
	}
	collection := os.Getenv("COLLECTION_NAME")

	client, ctx, err := connectMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB Atlas: %v", err)
	}
	defer client.Disconnect(ctx)

	bucketOptions := options.GridFSBucket().SetName(collection)
	bucket, err := gridfs.NewBucket(client.Database(databaseName), bucketOptions)
	if err != nil {
		log.Fatalf("Erro ao criar bucket do GridFS: %v", err)
	}

	switch command {
	case "up":
		if err := uploadFile(bucket, fileName); err != nil {
			log.Fatalf("Erro ao enviar arquivo: %v", err)
		}
	case "down":
		outputPath := "./uploads/"+fileName 
		if err := downloadFile(bucket, fileName, outputPath); err != nil {
			log.Fatalf("Erro ao baixar arquivo: %v", err)
		}
	case "list":
		if w,err := listFiles(bucket); err != nil {
			log.Fatalf("Erro ao listar arquivos: %v", err)
		}else{return w}

	case "delete":
		if err := deleteFile(bucket, fileName); err != nil {
			log.Fatalf("Erro ao deletar arquivo: %v", err)
		}
	default:
		log.Fatalf("Comando desconhecido: %v. Use 'up' para upload, 'down' para download, 'list' para listar ou 'delete' para excluir.", command)
	}
	return make([]string, 0)
}

/*

------------------------------------------------------------

                        usage

func main() {
    Grider(os.Args[1:]...)
}

// Comandos:
//   go run . up <filename>        : Faz upload do arquivo
//   go run . down <filename>      : Baixa o arquivo
//   go run . list                 : Lista todos os arquivos
//   go run . delete <filename>    : Exclui o arquivo

------------------------------------------------------------

*/
