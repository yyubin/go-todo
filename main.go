package main

import (
	"net/http"
	"todo/app"

	"github.com/urfave/negroni"
)

func main() {
	m := app.MakeHandler("./test.db") // appHandler 반환
	defer m.Close()
	n := negroni.Classic()
	n.UseHandler(m)

	http.ListenAndServe(":3000", n)

}
