package cache

import "time"

type AppCache interface {
	// สำหรับดึงข้อมูลจาก cache
	Get(key string) ([]byte, error)

	// สำหรับเพิ่มข้อมูลใน cache
	Set(key string, val []byte, exp time.Duration) error

	// สำหรับลบข้อมูลใน cache
	Delete(key string) error

	// สำหรับเรียกใช้งาน callback หาก cache หมดอายุ
	ExpiredEvent(callback func(key string) error) error
}
