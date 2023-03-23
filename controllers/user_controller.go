package controllers

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetUser(c echo.Context) error {
	// mendapatkan koneksi database dari konteks request
	db := c.Get("db").(*sql.DB)

	// melakukan query untuk mengambil data dari tabel 'users' dalam database
	rows, err := db.Query("SELECT id, name, age, address, type FROM users")
	if err != nil {
		// jika terjadi error saat melakukan query, kembalikan response dengan error bawaan Echo
		return echo.NewHTTPError(http.StatusInternalServerError, "Error querying users: "+err.Error())
	}
	defer rows.Close()

	// membuat slice kosong untuk menampung data pengguna
	var users []User

	// melakukan iterasi pada setiap baris data yang ditemukan dari hasil query
	for rows.Next() {
		var user User

		// mengisi data pengguna dari setiap kolom dalam baris saat ini
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Type); err != nil {
			// jika terjadi error saat memindai kolom dalam baris saat ini, kembalikan response dengan error bawaan Echo
			return echo.NewHTTPError(http.StatusInternalServerError, "Error scanning user row: "+err.Error())
		}

		// menambahkan pengguna saat ini ke dalam slice users
		users = append(users, user)
	}

	// kembalikan data pengguna dalam bentuk JSON dengan status code 200
	return c.JSON(http.StatusOK, users)
}

func InsertUser(c echo.Context) error {
	// mendapatkan koneksi database dari konteks request
	db := c.Get("db").(*sql.DB)

	// Mendapatkan data input dari request form.
	name := c.FormValue("name")
	ageStr := c.FormValue("age")
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		// jika input data usia tidak valid, kembalikan response dengan status code 400 dan pesan error
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid age"})
	}
	address := c.FormValue("address")
	typeStr := c.FormValue("type")
	userType, err := strconv.Atoi(typeStr)
	if err != nil {
		// jika input data tipe pengguna tidak valid, kembalikan response dengan status code 400 dan pesan error
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user type"})
	}

	// Validasi input data, jika ada input data yang tidak valid maka kembalikan error.
	if name == "" || age <= 0 || address == "" || userType <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid input data",
		})
	}

	// Mempersiapkan statement untuk memasukkan data user ke dalam database.
	stmt, err := db.Prepare("INSERT INTO users(name, age, address, type) VALUES (?, ?, ?, ?)")
	if err != nil {
		// jika terjadi error saat mempersiapkan statement, kembalikan response dengan status code 500 dan pesan error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer stmt.Close()

	// Menjalankan statement untuk memasukkan data user ke dalam database.
	result, err := stmt.Exec(name, age, address, userType)
	if err != nil {
		// jika terjadi error saat menjalankan statement, kembalikan response dengan status code 500 dan pesan error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Mendapatkan ID user yang baru saja dimasukkan ke dalam database.
	id, err := result.LastInsertId()
	if err != nil {
		// jika terjadi error saat mendapatkan ID user yang baru saja dimasukkan, kembalikan response dengan status code 500 dan pesan error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Membuat object User baru dengan data yang sudah dimasukkan ke dalam database.
	user := User{
		ID:      int(id),
		Name:    name,
		Age:     age,
		Address: address,
		Type:    userType,
	}

	// Mengembalikan response dengan status 201 Created dan data user yang baru saja dimasukkan ke dalam database.
	return c.JSON(http.StatusCreated, user)
}

func UpdateUser(c echo.Context) error {
	// mendapatkan koneksi database dari konteks request
	db := c.Get("db").(*sql.DB)

	// Mendapatkan ID user yang ingin diupdate dari parameter URL.
	id := c.Param("id")

	// Mendapatkan data input dari request form.
	name := c.FormValue("name")
	age := c.FormValue("age")
	address := c.FormValue("address")
	userType := c.FormValue("type")

	// Validasi input data, jika ada input data yang tidak valid maka kembalikan error.
	if name == "" || age == "" || address == "" || userType == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request data",
		})
	}

	// Menjalankan statement untuk mengupdate data user ke dalam database berdasarkan ID.
	result, err := db.Exec("UPDATE users SET name=?, age=?, address=?, type=? WHERE id=?", name, age, address, userType, id)
	if err != nil {
		// jika terjadi error saat menjalankan query, kembalikan response dengan status code 500 dan pesan error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// memeriksa apakah data pengguna yang ingin diupdate ada di dalam database
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// jika terjadi error saat memeriksa apakah data pengguna yang ingin diupdate ada di dalam database, kembalikan response dengan status code 500 dan pesan error
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if rowsAffected == 0 {
		// jika data pengguna yang ingin diupdate tidak ditemukan di dalam database, kembalikan response dengan status code 404 dan pesan error
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Mengembalikan response dengan status 200 OK dan pesan sukses.
	return c.JSON(http.StatusOK, map[string]string{
		"message": "User updated successfully",
	})
}

func DeleteUser(c echo.Context) error {
	// mendapatkan koneksi database dari konteks request
	db := c.Get("db").(*sql.DB)

	// Mendapatkan user ID dari parameter query
	userID := c.Param("id")

	// Menghapus user dari database
	result, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		// Jika terjadi error saat menghapus user, kembalikan response dengan status code 500 dan pesan error
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	// Mendapatkan jumlah baris yang terkena penghapusan
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Jika terjadi error saat mendapatkan jumlah baris yang terkena penghapusan, kembalikan response dengan status code 500 dan pesan error
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
	}

	if rowsAffected == 0 {
		// Jika tidak ada baris yang terkena penghapusan, kembalikan response dengan status code 404 dan pesan user tidak ditemukan
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "User not found",
		})
	}

	// Mengembalikan pesan JSON berhasil jika penghapusan berhasil
	return c.JSON(http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

//func DeleteUser(c echo.Context) error {
//	// mendapatkan koneksi database dari konteks request
//	db := c.Get("db").(*sql.DB)
//
//	// Mendapatkan user ID dari parameter query
//	userID := c.Param("id")
//
//	// Menghapus user dari database
//	result, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
//	if err != nil {
//		// Jika terjadi error saat menghapus user, kembalikan response dengan status code 500 dan pesan error
//		return c.JSON(http.StatusInternalServerError, map[string]string{
//			"message": "Failed to delete user",
//			"error":   err.Error(),
//		})
//	}
//
//	// Mendapatkan jumlah baris yang terkena penghapusan
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		// Jika terjadi error saat mendapatkan jumlah baris yang terkena penghapusan, kembalikan response dengan status code 500 dan pesan error
//		return c.JSON(http.StatusInternalServerError, map[string]string{
//			"message": "Failed to delete user",
//			"error":   err.Error(),
//		})
//	}
//
//	if rowsAffected == 0 {
//		// Jika tidak ada baris yang terkena penghapusan, kembalikan response dengan status code 404 dan pesan user tidak ditemukan
//		return c.JSON(http.StatusNotFound, map[string]string{
//			"message": "User not found",
//		})
//	}
//
//	// Mengembalikan pesan JSON berhasil jika penghapusan berhasil
//	return c.JSON(http.StatusOK, map[string]string{
//		"message": "User deleted successfully",
//	})
//}
