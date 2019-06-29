package main

import (
	"context"
	"fmt"
	"log"

	"github.com/arkiant/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial(":50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create the blog
	fmt.Println("Creating the blog")

	blog := &blogpb.Blog{
		AuthorId: "Samuel",
		Title:    "My first Blog",
		Content:  "Content of the first blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	fmt.Printf("Blog has been created: %v\n", createBlogRes)

}
