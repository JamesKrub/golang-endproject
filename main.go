package main

import (
	"KBTGCourse/project/pq/post"
	"log"
)

func main() {
	post.ConnectDB()
	log.Fatal(startServer())
}
