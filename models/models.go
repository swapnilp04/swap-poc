package models

func init() {
	migrateUser()
	migrateSession()
	migrateStudent()
	migrateBatch()
	migrateBatchStandard()
	migrateHostel()
	migrateHostelRoom()
	migrateHostelStudent()
	migrateStandard()
	migrateBatch()
	migrateHostelStudent()
	migrateBatchStandardStudent()
	migrateTransactionCategory()
	migrateTransactionDiscountCategory()
	migrateTransaction()
	migrateCheque()
	migrateStudentAccount()
	migrateCommentCategory()
	migrateComment()
	migrateCommentCategoryData() //Data load
	migrateExam()
	migrateExamStudent()
	migrateSubject()
	migrateTeacher()
	migrateLogCategory()
	//migrateLogCategoryData() // data load
	migrateTeacherLog()
	migrateChapter()
}
