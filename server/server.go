package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	pb "newapp/proto"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCrudServer
	database *sql.DB
}

// error handler to make code less and clear)
func errHandler(statement string, err error) {
	if err != nil {
		log.Fatalf("failed to %s, %v", statement, err)
	}
}

// create function to handle the request from client
func (server *server) Create(ctx context.Context, request *pb.UserInfo) (*pb.UserInfo, error) {
	err := server.database.Ping()
	errHandler("connect to the database", err)
	_, err = server.database.Exec("insert into users (username, password, firstname, lastname) values ($1, $2, $3, $4)",
		request.Info.GetUsername(),
		request.Info.GetPassword(),
		request.GetFirstname(),
		request.GetLastname(),
	)
	errHandler("insert user to the database", err)
	return request, nil
}

// collecting the full info about user from database
func (server *server) Get(ctx context.Context, request *pb.PrivateUserInfo) (*pb.UserInfo, error) {
	err := server.database.Ping()
	errHandler("connect to the database", err)
	row := server.database.QueryRow("select firstname, lastname from users where username=$1 and password=$2",
		request.GetUsername(),
		request.GetPassword(),
	)
	userinfo := new(pb.UserInfo)
	userinfo.Info = request
	err = row.Scan(&userinfo.Firstname, &userinfo.Lastname)
	errHandler("get information about user from database", err)
	return userinfo, nil
}

func (server *server) Update(ctx context.Context, request *pb.UserInfo) (*pb.UserInfo, error) {
	err := server.database.Ping()
	errHandler("connect to the database", err)
	_, err = server.database.Exec("update users set username=$1, password=$2, firstname=$3, lastname=$4 where username=$1 and password=$2",
		request.Info.GetUsername(),
		request.Info.GetPassword(),
		request.GetFirstname(),
		request.GetLastname(),
	)
	errHandler("update user info", err)
	return request, nil
}

func (server *server) Delete(ctx context.Context, request *pb.PrivateUserInfo) (*pb.UserInfo, error) {
	err := server.database.Ping()
	errHandler("connect to the database", err)
	row := server.database.QueryRow("select firstname, lastname from users where username=$1 and password=$2",
		request.GetUsername(),
		request.GetPassword(),
	)
	userinfo := new(pb.UserInfo)
	userinfo.Info = request
	err = row.Scan(&userinfo.Firstname, &userinfo.Lastname)
	errHandler("get information about user from database", err)
	_, err = server.database.Exec("delete from users where username=$1 and password=$2", request.GetUsername(), request.GetPassword())
	errHandler("delete user from db", err)
	return userinfo, nil
}

// establishing connect to the database
func connDB(host, port, username, password, dbname string) *sql.DB {
	conn, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, username, password, dbname,
		),
	)
	errHandler("connect to database", err)
	return conn
}

// declaring the constants for the database
const (
	host     = "127.0.0.1"
	port     = "5432"
	username = "otabek"
	password = "otabek123"
	dbname   = "test"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	errHandler("listen", err)
	serv := grpc.NewServer()
	pb.RegisterCrudServer(serv, &server{database: connDB(host, port, username, password, dbname)})
	err = serv.Serve(lis)
	errHandler("serve", err)
}
