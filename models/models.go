package models

func init() {
	migrateUser()
	migrateSession()
	migrateStudent()
	migrateBatch()
	migrateHostel()
	migrateHostelRoom()
	migrateHostelStudent()
}
