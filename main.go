package main

import (
	"net/http"
	"todo/app"
)

func main() {
	m := app.MakeHandler("./test.db") // appHandler 반환
	defer m.Close()

	http.ListenAndServe(":3000", m)

}
