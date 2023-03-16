package controllers

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Type    int    `json:"type"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Transaction struct {
	ID        int      `json:"id"`
	UserID    int      `json:"userId"`
	ProductID int      `json:"productId"`
	Quantity  int      `json:"quantity"`
	User      *User    `json:"user"`
	Product   *Product `json:"product"`
}
