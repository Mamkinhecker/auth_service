package db

type User struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`
	Password    string `db:"password_hash"`
	PhotoObj    string `db:"photo_object"`
}
