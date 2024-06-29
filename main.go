package main

import "net/http"

func main() {
	sql, err := NewSqlStorage()
	if err != nil {
		panic(err)
	}
	server := &Server{
		ChoreStorage: sql,
	}
	http.ListenAndServe(":8081", server.GetHandler())
	
}
