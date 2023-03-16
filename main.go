// Nama : David Kharis Elio M
// NIM : 1121028
package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/latihan_framework/controllers"
	_ "github.com/latihan_framework/controllers"
	//"log"
	//"net/http"

	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"
)

func main() {
	// membuat instance baru dari Echo web framework
	e := echo.New()

	// menambahkan middleware pada instance Echo yang akan menambahkan koneksi database pada setiap request
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// membuat koneksi database baru dan menyimpannya ke dalam context Echo
			c.Set("db", controllers.Connect())
			return next(c)
		}
	})

	// menambahkan routing untuk masing-masing HTTP method
	// untuk mengakses endpoint '/users' dengan memanggil fungsi controller yang sesuai
	e.GET("/users", controllers.GetUser)
	e.POST("/users", controllers.InsertUser)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	// mencetak pesan ke konsol bahwa server telah dimulai pada port 8080
	fmt.Println("Server started at :8080")

	// memulai server dan menunggu permintaan dari klien
	e.Start(":8080")
}
