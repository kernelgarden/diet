package factory

import "github.com/kernelgarden/diet/constant"

func NewFailResp(failCode int) constant.FailResp {
	var msg string

	switch failCode {
	case constant.Wrong:
		msg = "틀린 값입니다.."
		break
	case constant.Invalid:
		msg = "유효하지 않은 값입니다.."
		break
	case constant.InExist:
		msg = "존재하지 않습니다."
		break
	case constant.NeedPermission:
		msg = "권한이 없습니다."
		break
	case constant.InvalidRequestFormat:
		msg = "잘못된 형식의 요청입니다."
		break
	default:
		msg = "알 수 없는 이유"
		break
	}

	return constant.FailResp{FailCode: failCode, Message: msg}
}