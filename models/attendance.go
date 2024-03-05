package models

type Attendance struct {
	ID            uint          `gorm:"primaryKey;size:255;autoIncrement;"`
	CID           uint          `gorm:"column:cid;not null"`                                                    // Class ID
	UID           uint          `gorm:"column:uid;not null"`                                                    // User ID
	CSID          uint          `gorm:"column:csid;not null"`                                                   // Class Schedule ID
	IsAttendance  string        `gorm:"type:enum('Attendance', 'Tardy', 'Absence');default:'Absence';not null"` // 出席, 遅刻, 欠席
	ClassUser     ClassUser     `gorm:"foreignKey:CID,UID"`
	ClassSchedule ClassSchedule `gorm:"foreignKey:CSID"`
}
