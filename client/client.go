package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	pb "newapp/proto"

	"google.golang.org/grpc"
)

type User struct {
	Username, Password, Firstname, Lastname string
}
type AuthInfo struct {
	Username, Password string
}

func create(w http.ResponseWriter, r *http.Request) {
	jsonbody, err := ioutil.ReadAll(r.Body)
	errHandler("read the response body", err)
	user := new(User)
	err = json.Unmarshal(jsonbody, &user)
	errHandler("unmarshall the response body", err)
	log.Print(user)
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	errHandler("establish connection with server", err)
	client, ctx := pb.NewCrudClient(conn), context.Background()
	response, err := client.Create(ctx, &pb.UserInfo{
		Info:      &pb.PrivateUserInfo{Username: user.Username, Password: user.Password},
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	})
	errHandler("connect with the grpc server", err)
	webresponse, err := json.Marshal(response)
	errHandler("marshaling the response form grpc server", err)
	w.Write(webresponse)
}

func read(w http.ResponseWriter, r *http.Request) {
	jsonbody, err := ioutil.ReadAll(r.Body)
	errHandler("read the response body", err)
	user := new(AuthInfo)
	err = json.Unmarshal(jsonbody, &user)
	errHandler("unmarshall the response body", err)
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	errHandler("establish connection with server", err)
	client, ctx := pb.NewCrudClient(conn), context.Background()
	response, err := client.Get(ctx, &pb.PrivateUserInfo{Username: user.Username, Password: user.Password})
	errHandler("connect with the grpc server", err)
	webresponse, err := json.Marshal(response)
	errHandler("marshal the response form grpc server", err)
	w.Write(webresponse)
}

func update(w http.ResponseWriter, r *http.Request) {
	jsonbody, err := ioutil.ReadAll(r.Body)
	errHandler("read the response body", err)
	user := new(User)
	err = json.Unmarshal(jsonbody, &user)
	errHandler("unmarshall the response body", err)
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	errHandler("establish connection with server", err)
	client, ctx := pb.NewCrudClient(conn), context.Background()
	response, err := client.Update(ctx, &pb.UserInfo{
		Info:      &pb.PrivateUserInfo{Username: user.Username, Password: user.Password},
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	})
	errHandler("connect with the grpc server", err)
	webresponse, err := json.Marshal(response)
	errHandler("marshaling the response form grpc server", err)
	w.Write(webresponse)
}

func delete(w http.ResponseWriter, r *http.Request) {
	jsonbody, err := ioutil.ReadAll(r.Body)
	errHandler("read the response body", err)
	user := new(AuthInfo)
	err = json.Unmarshal(jsonbody, &user)
	errHandler("unmarshall the response body", err)
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	errHandler("establish connection with server", err)
	client, ctx := pb.NewCrudClient(conn), context.Background()
	response, err := client.Delete(ctx, &pb.PrivateUserInfo{Username: user.Username, Password: user.Password})
	errHandler("connect with the grpc server", err)
	webresponse, err := json.Marshal(response)
	errHandler("marshal the response form grpc server", err)
	w.Write(webresponse)
}

func errHandler(statement string, err error) {
	if err != nil {
		log.Fatalf("failed to %s, %v", statement, err)
	}
}

func main() {
	http.HandleFunc("/create", create)
	http.HandleFunc("/read", read)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	err := http.ListenAndServe(":9000", nil)
	errHandler("listen and serve)", err)
}
