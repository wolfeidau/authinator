package api

func errorMsg(msg string) map[string]string {
	return map[string]string{"msg": msg}
}

func validationErrors(msg, allErrs interface{}) map[string]interface{} {
	return map[string]interface{}{"msg": msg, "errors": allErrs}
}
