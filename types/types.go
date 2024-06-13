package types

type UserLogin struct {
	Tenant   string `json:"tenant_id" dynamodbav:"tenant_id"`
	Userid   string `json:"user_id" dynamodbav:"user_id"`
	Password string `json:"userkey" dynamodbav:"user_key"`
}

type User struct {
	Userid       string `json:"user_id" dynamodbav:"user_id"`
	Useremail    string `json:"useremail" dynamodbav:"tenant_id"` // TODO: change
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

type GetAllTableRes struct {
	User_id string `json:"user_id"`
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

func NewUserTables(user_id string, info *[]TableInfo) *UserTables {
	return &UserTables{
		Userid:    user_id,
		Tabledata: *info,
	}
}

func (u *UserTables) Format() *ResErr {
	return &ResErr{
		Err:  "",
		Code: 0,
		Data: u,
	}
}

type TableInfo struct {
	TableID     string            `json:"tableid" dynamodbav:"table_id"`
	TableName   string            `json:"tablename" dynamodbav:"tablename"`
	Discription string            `json:"tablediscription" dynamodbav:"discription"`
	LastUpdate  int64             `json:"lastupdate" dynamodbav:"lastdate"`
	Color       string            `json:"color" dynamodbav:"color"`
	Columns     map[string]string `json:"Column" dynamodbav:"Column"`
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
