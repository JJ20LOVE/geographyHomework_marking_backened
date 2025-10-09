package utils

var CodeMsg = map[int]string{
	200: "success",
	201: "bind failed",
	203: "upload file failed",
	204: "",
	205: "extract file failed",
	206: "file not found",
	207: "json error",
	//sign up
	300: "validation failed",
	301: "user already exists",
	//login
	310: "user not found",
	311: "password error",
	//token
	320: "token not found",
	321: "token invalid",
	322: "token generate failed",
	//db
	400: "db error",
	//student
	500: "student not found",
	501: "student already exists",
	502: "class not found",
	503: "class already exists",
	504: "you can't delete a class with students",
	505: "you can't delete a student with answer sheets",
	//answersheet&exam
	600: "answersheet number error",
	601: "answer sheet already exists",
	602: "evaluation interface error",
	603: "upload answer failed",
	604: "exam already exists",
	605: "exam not found",
	606: "answersheet not found",
	607: "you can't delete a exam with answer sheets",
	608: "ocr timeout",
	609: "delete answersheet failed",
}

func GetErrMsg(code int) string {
	if msg, ok := CodeMsg[code]; ok {
		return msg
	}
	return "未知错误"
}
