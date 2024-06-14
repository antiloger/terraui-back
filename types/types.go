package types

type DataFormat interface {
	Format() *ResErr
}

// ------------this is the response template ----------------
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

//------------END: this is the response template ----------------

type UserLogin struct {
	Tenant   string `json:"user_email" dynamodbav:"tenant_id"`
	Password string `json:"userkey" dynamodbav:"user_key"`
}

type UserRegister struct {
	Useremail string `json:"user_email" dynamodbav:"tenant_id"`
	Userid    string `json:"user_id" dynamodbav:"user_id"`
	Username  string `json:"username" dynamodbav:"username"`
	Role      string `json:"role" dynamodbav:"role"`
	Userkey   string `json:"userkey" dynamodbav:"userkey"`
}

type User struct {
	Userid       string `json:"user_id" dynamodbav:"user_id"`
	Useremail    string `json:"useremail" dynamodbav:"tenant_id"` // TODO: change
	Username     string `json:"username" dynamodbav:"username"`
	Userrole     string `json:"role" dynamodbav:"role"`
	Subscription string `json:"subscription" dynamodbav:"subscription"`
	Userkey      string `json:"userkey" dynamodbav:"userkey"`
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

type GetAllItemRes struct {
	Table_id string `json:"table_id"`
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

type TableItems struct {
	Tableinfo *TableInfo        `json:"tabeleinfo"`
	Tabledata *[]map[string]any `json:"tabledata"`
}

func NewTableItems(i *TableInfo, d *[]map[string]any) *TableItems {
	return &TableItems{
		Tableinfo: i,
		Tabledata: d,
	}
}

func (d *TableItems) Format(e string, c int8) *ResErr {
	return &ResErr{
		Err:  e,
		Code: c,
		Data: d,
	}
}
