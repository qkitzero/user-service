package user

type UserTable struct {
	ID    string
	Name  string
	Email string
}

func (UserTable) TableName() string {
	return "user"
}
