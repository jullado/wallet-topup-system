package logs

type AppLog interface {
	// Info logging สำหรับแสดงการทำงานต่างๆ
	Info(msg string)

	// Debug logging สำหรับทดสอบการทำงาน
	Debug(msg string)

	// Warning logging สำหรับแสดงความผิดผลาด
	Warning(msg string)

	// Error logging เมื่อเกิดความผิดพลาดร้ายแรง
	Error(msg interface{})
}
