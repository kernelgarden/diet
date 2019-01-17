package constant

const (
	Unknown					= 10
	Wrong		  			= 30
	Invalid    				= 40
	InExist					= 50
	NeedPermission 			= 60
	InvalidRequestFormat 	= 70
)

type FailResp struct {
	FailCode	int		`json:"failCode"`
	Message 	string	`json:"message"`
}

