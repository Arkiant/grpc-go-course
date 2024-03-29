package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/arkiant/grpc-go-course/blog/blogpb"

	"github.com/arkiant/grpc-go-course/blog/database"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := req.GetBlog()

	data := database.CreateBlog(blog.GetAuthorId(), blog.GetTitle(), blog.GetContent())

	oid, err := database.InsertOne(data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid,
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil

}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogID := req.GetBlogId()
	data, err := database.FindOneByID(blogID)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil

}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()

	data := database.CreateBlog(blog.GetAuthorId(), blog.GetTitle(), blog.GetContent())

	if updateError := database.ReplaceOneByID(data, blog.GetId()); updateError != nil {
		return nil, status.Errorf(
			codes.Internal,
			updateError.Error(),
		)
	}

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blog.GetId(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Content:  data.Content,
		},
	}, nil

}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	blogID := req.GetBlogId()

	res, err := database.DeleteByID(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog in MongoDB: %v\n", err),
		)
	}

	return &blogpb.DeleteBlogResponse{BlogId: blogID}, nil
}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {

	collection := database.MongoCollection()
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v\n", err),
		)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := database.CreateBlog("", "", "")
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v\n", err),
			)
		}

		stream.Send(&blogpb.ListBlogResponse{
			Blog: &blogpb.Blog{
				Id:       data.ID.Hex(),
				AuthorId: data.AuthorID,
				Title:    data.Title,
				Content:  data.Content,
			},
		})
	}

	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v\n", err),
		)
	}

	return nil
}

func main() {

	// If we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Conecting to MongoDB")
	collection := database.MongoCollection()
	fmt.Printf("Conected to %s collection\n", collection.Name())

	fmt.Println("Blog Service Started")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	// go get github.com/ktr0731/evans
	// cli inspection
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Stopping the listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	database.CloseConnection()
	fmt.Println("End of Program")
}
