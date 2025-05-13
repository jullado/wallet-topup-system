package utils

// เป็น Extention Method โดยจะส่งจาก Business Logic ไปยัง Presentation Layer โดยจะ Conform ตาม Type Error
type ErrHandler struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"error message"`
}

func (e ErrHandler) Error() string {
	return e.Message
}
