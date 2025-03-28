# nacon
 
## Description

used to upload, download and delete files from mongoDB using GridFS

## Usage

```bash
  go run . [command] [filename]
  example: go run . up test.txt
  example: go run . down test.txt
  example: go run . list 
  example: go run . delete test.txt
  
  ## For Images

  Supported image formats:
  - JPEG/JPG
  - PNG
  - GIF
  - BMP
  - TIFF

  Example usage with images:
  go run . up photo.jpg
  go run . down image.png
  go run . delete picture.gif
  go run . up test.jpeg
  go run . down test.jpeg
  go run . delete test.jpeg
  
  # List all files in the database
  go run . list
```

## API Reference

| Parameter | Type | Description |
|-----------|------|-------------|
| `down or up or delete or list` | `string` | down for download, up for upload, delete for delete file, list to list all files |
| `filename` | `string` | filename |


## Warning! Add .env file with MongoDB connection string 
```
MONGO_URI=your_mongodb_connection_string
DATABASE_NAME=your_database_name
COLLECTION_NAME=your_collection_name
```
