package types

type Student struct {
	Id    int
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,gt=0,lt=150"`
}