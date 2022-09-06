// GRPC server implementation of TodoServiceServer
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	pb "github.com/yuhuishi-convect/grpc-todo/gen/proto"
	grpc "google.golang.org/grpc"
)

// CLI flag for grpc server port
var (
	port = flag.Int("port", 8080, "The server port")
)

// initiate a database connection to local sqlite database
func initDB() (*sql.DB, error) {
	// create a database connection
	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	// create a table if not exists
	log.Printf("Creating table if not exists")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todo (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, done BOOLEAN)")
	return db, err
}

type TodoServiceServer struct {
	pb.UnimplementedTodoServiceServer
}

func (s *TodoServiceServer) List(ctx context.Context, in *pb.ListRequest) (*pb.ListResponse, error) {
	// fetch all todo items from the database
	rows, err := db.Query("SELECT id, title, done FROM todo")
	if err != nil {
		log.Fatalf("failed to query database: %v", err)
	}
	defer rows.Close()

	todoItem := []*pb.TodoItem{}

	for rows.Next() {
		var id int64
		var title string
		var done bool
		if err := rows.Scan(&id, &title, &done); err != nil {
			log.Fatalf("failed to scan row: %v", err)
		}
		// append the row to the todo item list
		todoItem = append(todoItem, &pb.TodoItem{
			Id:    fmt.Sprintf("%d", id),
			Title: title,
			Done:  done,
		})
	}

	return &pb.ListResponse{Items: todoItem}, nil
}

// Create a todo item handler function
func (s *TodoServiceServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {

	// prepare statement to insert a todo item
	stmt, err := db.Prepare("INSERT INTO todo(title, done) VALUES(?, ?)")
	if err != nil {
		log.Fatalf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// insert the todo item using the prepared statement
	res, err := stmt.Exec(in.GetTitle(), false)
	if err != nil {
		log.Fatalf("failed to insert todo item: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("failed to get last insert id: %v", err)
	}

	// write the todo item to the database
	return &pb.CreateResponse{
		Item: &pb.TodoItem{
			Id:    fmt.Sprintf("%d", id),
			Title: in.GetTitle(),
			Done:  false,
		},
	}, nil
}

var db *sql.DB

// start the grpc server
func main() {
	flag.Parse()

	// initiate the database connection
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("failed to initiate database: %v", err)
	}
	defer db.Close()

	server := grpc.NewServer()
	log.Printf("Starting grpc server on port %d", *port)
	pb.RegisterTodoServiceServer(server, &TodoServiceServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
