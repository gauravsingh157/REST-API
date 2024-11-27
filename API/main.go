package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	name string
	age  int
}
type errorResp struct {
	Statuscode int
	massage    string
}

var (
	user = map[string]User{}
)

func main() {

	http.HandleFunc("/", AddUser)

	fmt.Println("Server started")
	log.Fatalf("server start nhi huaa hai ,err : %v\n", http.ListenAndServe(":8000", nil))

}

func AddUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	User := User{}

	err := json.NewDecoder(r.Body).Decode(&User)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errorResp{
			Statuscode: http.StatusBadRequest,
			massage:    "JO PAYLOAD AYA HAI WO SAHI NHI HAI,DEYAN DO US PAR ERR =" + err.Error(),
		}
		json.NewEncoder(w).Encode(err)

		return

		user[User.name] = User
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

		fmt.Println("users are : ",)

		return
	}
}
