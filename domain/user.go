package domain

type User struct {
	Name     string
	Login    string
	NickName string
	Points   int
	Weekend  Weekend
}

func (u *User) AddPoint(point int) {
	u.Points += point
}

func (u *User) DeletePoint(point int) {
	u.Points -= point
}
