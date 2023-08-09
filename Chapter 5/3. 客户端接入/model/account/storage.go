package account

import "ProjectX/library/storage/core"

type Account struct {
	core.Model
	UserId       string `dynamo:",hash"`
	AccountId    string `index:"AccountId-Global,hash"`
	Platform     string
	RegisterTime int64
}
