package controllers

func restReturn(errcode int, errmsg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"errcode": errcode,
		"errmsg":  errmsg,
		"data":    data,
	}
}
