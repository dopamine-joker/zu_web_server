package dao

import "time"

type User struct {
	Id         int32     `json:"Id" db:"id"`
	Email      string    `json:"Email" db:"email"`
	Name       string    `json:"Name" db:"name"`
	Password   string    `json:"Password" db:"password"`
	CreateTime time.Time `json:"Create_time" db:"create_time"`
}
