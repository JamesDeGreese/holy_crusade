package core

type NewUser struct {
	ChatID int64
}

type CityInfoReq struct {
	ChatID int64
}

type AddWorkerReq struct {
	ChatID int64
	Count  int
}

type AddSolderReq struct {
	ChatID int64
	Count  int
}

type Response struct {
	ChatID  int64
	Payload interface{}
}

type CityInfoRes struct {
	Name       string
	Rating     int
	Gold       int
	Population int
	Workers    int
	Solders    int
	Heroes     int
}
