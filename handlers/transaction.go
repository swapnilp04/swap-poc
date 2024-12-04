package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
)

func GetTransactions(c echo.Context) error {
	t := &models.Transaction{}
	s := &models.Student{}
	page := c.QueryParam("page")
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}

	search := c.QueryParam("search")

	err, ids := s.SearchIds(search)
	if err != nil {
		fmt.Println("s.ALL(SearchStudents)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	transactions, err := t.AllStudents(int(newPage), ids)
	if err != nil {
		fmt.Println("s.ALL(GetTransactions)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	count, err := t.Count(ids)
	if err != nil {
		fmt.Println("s.ALL(GetTransactions)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"transactions": transactions, "total": count})
}

func GetStudentTransactions(c echo.Context) error {
	// Get student
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	transaction := &models.Transaction{}
	transactions, err := transaction.All(newStudentId)
	if err != nil {
		fmt.Println("s.ALL(GetTransactions)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, transactions)
}

func GetStudentBalance(c echo.Context) error {
	// Get student
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	debits, credits, discounts := student.GetBalance()
	return c.JSON(http.StatusOK, map[string]interface{}{"debits": debits, "credits": credits, "discounts": discounts})
}

func GetStudentTransaction(c echo.Context) error {
	// Get a single user by ID
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	Id := c.Param("id")
	newId, err := strconv.Atoi(Id)
	transaction , err := student.GetTransaction(uint(newId))
	if err != nil {
		fmt.Println("s.Find(GetTransaction)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}	
	transaction.AddWordPayment()
	return c.JSON(http.StatusOK, transaction)
}

func PayStudentFee(c echo.Context) error {
	// Get a single user by ID
	cc := c.(CustomContext)
	session := cc.session
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	transactionData := make(map[string]interface{})
	if err := c.Bind(&transactionData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	transaction := models.NewTransaction(transactionData, *student)
	transaction.TransactionType = "cridit"
	transaction.Name = "Pay Fee"
	transaction.UserID = uint(session.UserID)

	if err := transaction.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}
	err = transaction.Create()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	if transaction.PaymentMode == "Cheque" {
		cheque := models.NewCheque(transactionData["Cheque"].(map[string]interface{}))
		cheque.TransactionId = transaction.ID
		cheque.Amount = transaction.Amount
		if err := cheque.Validate(); err != nil {
			transaction.Delete()
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		}
		err = cheque.Create()
		if err != nil {
			transaction.Delete()
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
		}
	}

	err = student.SaveBalance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Transaction created", "transaction": transaction})
}

func AddStudentDues(c echo.Context) error {
	// Get a single user by ID
	cc := c.(CustomContext)
	session := cc.session
	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	transactionData := make(map[string]interface{})
	if err := c.Bind(&transactionData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	transaction := models.NewTransaction(transactionData, *student)
	transaction.TransactionType = "debit"
	transaction.Name = "Dues"
	transaction.PaymentMode = "-"
	transaction.PaidBy = "-"
	transaction.UserID = uint(session.UserID)
	if err := transaction.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}
	err = transaction.Create()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	err = student.SaveBalance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Transaction created", "transaction": transaction})
}

func AddDiscount(c echo.Context) error {
	// Get a single user by ID
	cc := c.(CustomContext)
	session := cc.session

	studentId := c.Param("student_id")
	newStudentId, err := strconv.Atoi(studentId)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}
	student := &models.Student{ID: uint(newStudentId)}
	err = student.Find()
	if err != nil {
		fmt.Println("s.Find(GetStudent)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	transactionData := make(map[string]interface{})
	if err := c.Bind(&transactionData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}
	
	transaction_category, err := models.GetDiscountTransactionCategory()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	transaction := models.NewTransaction(transactionData, *student)
	transaction.TransactionType = "cridit"
	transaction.Name = "Discount"
	transaction.TransactionCategoryId = transaction_category.ID
	transaction.UserID = uint(session.UserID)

	if err := transaction.Validate(); err != nil {
		formErr := MarshalFormError(err)	
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": formErr})
	}
	err = transaction.Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	err = student.SaveBalance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Transaction created", "transaction": transaction})
}