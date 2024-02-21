package models

type Attendance struct {
	Key          string    `gorm:"primaryKey;size:255"`
	CID          uint      `gorm:"not null"` // Class ID
	UID          uint      `gorm:"not null"` // User ID
	AttendanceID string    `gorm:"size:255;not null"`
	IsAttendance string    `gorm:"type:enum('Attendance', 'Tardy', 'Absence');default:'Absence';not null"` // 出席, 遅刻, 欠席
	ClassUser    ClassUser `gorm:"foreignKey:CID,UID"`
}
