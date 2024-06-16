package models

type AttendanceType string

const (
	AttendanceStatus AttendanceType = "ATTENDANCE"
	TardyStatus      AttendanceType = "TARDY"
	AbsenceStatus    AttendanceType = "ABSENCE"
)

type Attendance struct {
	ID            uint           `gorm:"primaryKey;size:255;autoIncrement;"`
	CID           uint           `gorm:"column:cid;not null"`                                                    // Class ID
	UID           uint           `gorm:"column:uid;not null"`                                                    // User ID
	CSID          uint           `gorm:"column:csid;not null"`                                                   // Class Schedule ID
	IsAttendance  AttendanceType `gorm:"type:enum('ATTENDANCE', 'TARDY', 'ABSENCE');default:'ABSENCE';not null"` // 出席, 遅刻, 欠席
	ClassUser     ClassUser      `gorm:"foreignKey:CID,UID"`
	ClassSchedule ClassSchedule  `gorm:"foreignKey:CSID"`
}
