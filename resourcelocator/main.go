package main

import "resourcelocator"


func main() {
	resourcelocator.CreateEmbeddedLocator("embeddedFiles.go", "main", "resources")
}