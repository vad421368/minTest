package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/check_health/", healthCheckHandler)
	mux.HandleFunc("/matrix_operations", matrixHandler)
	mux.HandleFunc("/heximal_operations", heximalHandler)
	http.ListenAndServe(":8000", mux)
}

type healthChecker struct {
	Status string `json:"status"`
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	var health healthChecker
	if r.Method == "GET" {
		health.Status = "OK"
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(&health)
		w.Write(response)
	} else {
		w.Write([]byte("Incorrect http method\n"))
	}
}

type MatrixRequest struct {
	Matrix1       [][]int `json:"matrix_1"`
	Matrix2       [][]int `json:"matrix_2"`
	OperationType string  `json:"operation_type"`
}

type MatrixResponse struct {
	Matrix        [][]int `json:"matrix"`
	OperationType string  `json:"operation_type"`
}
type MatrixResponseError struct {
	Error         string `json:"error"`
	OperationType string `json:"operation_type"`
}

func matrixHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		request := &MatrixRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			fmt.Println("Error decoding request body")
		}
		responseError := MatrixResponseError{}
		if len(request.Matrix1) < 2 || len(request.Matrix2) < 2 || len(request.Matrix1[0]) != len(request.Matrix2) {
			responseError.Error = "Matrix dimension mismatch"
			responseError.OperationType = request.OperationType
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			reply, _ := json.Marshal(&responseError)
			w.Write(reply)
		} else {
			if request.OperationType == "multiply" {
				response := MatrixResponse{}
				response.Matrix = multiplyMatrix(request.Matrix1, request.Matrix2)

				response.OperationType = request.OperationType
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				reply, _ := json.Marshal(&response)
				w.Write(reply)

			} else {
				w.Write([]byte("Incorrect http method\n"))
			}
		}
	}

}

func multiplyMatrix(matrix1, matrix2 [][]int) [][]int {
	rows1, cols1, cols2 := len(matrix1), len(matrix1[0]), len(matrix2[0])

	C := make([][]int, rows1)
	for i := range C {
		C[i] = make([]int, cols2)
	}
	for i := 0; i < rows1; i++ {
		for j := 0; j < cols2; j++ {
			for k := 0; k < cols1; k++ {
				C[i][j] += matrix1[i][k] * matrix2[k][j]
			}
		}
	}
	return C
}

type HexRequest struct {
	Heximal1      string `json:"heximal_1"`
	Heximal2      string `json:"heximal_2"`
	OperationType string `json:"operation_type"`
}

type HexResponse struct {
	Heximal       string `json:"heximal"`
	OperationType string `json:"operation_type"`
}

type HexResponseError struct {
	Error         string `json:"error"`
	OperationType string `json:"operation_type"`
}

func heximalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		request := &HexRequest{}
		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(request)
		if err != nil {
			fmt.Println("Error decoding request body")
		}
		responseError := &HexResponseError{}
		if !hexChecker(request.Heximal1) || !hexChecker(request.Heximal2) {
			responseError.Error = "Not a heximal number"
			responseError.OperationType = request.OperationType
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			reply, _ := json.Marshal(&responseError)
			w.Write(reply)
		} else {
			if request.OperationType == "multiply" {
				response := &HexResponse{}
				response.Heximal = multiplyHex(request.Heximal1, request.Heximal2)

				response.OperationType = request.OperationType
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				reply, _ := json.Marshal(&response)
				w.Write(reply)

			}
		}
	} else {
		w.Write([]byte("Incorrect http method\n"))
	}
}

func hexChecker(s string) bool {
	t := []byte(s)
	re := regexp.MustCompile(`^[0-9A-Fa-f]+$`)
	return re.Match(t)
}

func multiplyHex(hex1, hex2 string) string {
	hex1Parsed, _ := strconv.ParseInt(hex1, 16, 64)
	d1 := strconv.FormatInt(hex1Parsed, 10)
	hex2Parsed, _ := strconv.ParseInt(hex2, 16, 64)
	d2 := strconv.FormatInt(hex2Parsed, 10)
	d1int, _ := strconv.Atoi(d1)
	d2int, _ := strconv.Atoi(d2)
	totalInt := d1int * d2int
	totalString := strconv.Itoa(totalInt)
	totalParsed, _ := strconv.ParseInt(totalString, 10, 64)
	totalHex := strconv.FormatInt(totalParsed, 16)

	return totalHex
}
