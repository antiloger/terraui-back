package types

type UserLogin struct {
	Tenant   string `json:"tenant_id" dynamodbav:"tenant_id"`
	Userid   string `json:"user_id" dynamodbav:"user_id"`
	Password string `json:"userkey" dynamodbav:"user_key"`
}

type User struct {
	Userid       string `json:"user_id" dynamodbav:"user_id"`
	Useremail    string `json:"useremail" dynamodbav:"useremail"`
	Username     string `json:"username" dynamodbav:"username"`
	Userrole     string `json:"role" dynamodbav:"role"`
	Subscription string `json:"subscription" dynamodbav:"subscription"`
}

func (u *User) Format(err string, code int8) *ResErr {
	return &ResErr{
		Err:  err,
		Code: code,
		Data: u,
	}
}

type AuthUser struct {
	Userid string
	Role   string
	Email  string
}

type UserTables struct {
	Userid    string      `json:"username"`
	Tabledata []TableInfo `json:"tabledata"`
}

func (u *UserTables) Format() *ResErr {
	return &ResErr{
		Err:  "",
		Code: 0,
		Data: u,
	}
}

type TableInfo struct {
	TableID     string `json:"tableid"`
	TableName   string `json:"tablename"`
	Discription string `json:"tablediscription"`
	LastUpdate  int64  `json:"lastupdate"`
}

type DataFormat interface {
	Format() *ResErr
}

type ResErr struct {
	Err  string      `json:"err"`
	Code int8        `json:"code"`
	Data interface{} `json:"data"`
}

func NewResErr(e string, c int8) ResErr {
	return ResErr{
		Err:  e,
		Code: c,
	}
}
