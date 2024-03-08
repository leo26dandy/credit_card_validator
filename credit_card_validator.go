package main // This line declares that this is the main package

import ( // This line starts the import section
	"encoding/json" // This imports the json package for JSON encoding/decoding
	"fmt"           // This imports the fmt package for formatted I/O operations
	"net/http"      // This imports the http package for creating an HTTP server
	"strings"       // This imports the strings package for string manipulation
)

// Request struct represents the JSON request payload
type Request struct {
	CreditCardNumber string `json:"credit_card_number"` // This field holds the credit card number
}

// Response struct represents the JSON response payload
type Response struct {
	IsValid bool `json:"is_valid"` // This field holds the validation result
}

func main() { // This is the entry point of the program
	http.HandleFunc("/validate", handleValidateRequest) // This registers the handleValidateRequest function as the handler for the /validate endpoint
	fmt.Println("Server listening on port 8080")        // This prints a message to indicate that the server is running
	http.ListenAndServe(":8080", nil)                   // This starts the HTTP server and listens on port 8080
}

func handleValidateRequest(w http.ResponseWriter, r *http.Request) { // This function handles the /validate endpoint
	if r.Method != http.MethodGet { // This checks if the request method is GET
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed) // If not, it returns a 405 error
		return
	}

	var req Request                             // This declares a variable of type Request
	err := json.NewDecoder(r.Body).Decode(&req) // This decodes the JSON request body into the req variable
	if err != nil {                             // This checks if there was an error during decoding
		http.Error(w, "Invalid JSON request", http.StatusBadRequest) // If there was an error, it returns a 400 error
		return
	}

	cardNumber := strings.ReplaceAll(req.CreditCardNumber, " ", "") // This removes spaces from the credit card number
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")            // This removes hyphens from the credit card number

	isValid := isValidCreditCardNumber(cardNumber) // This calls the isValidCreditCardNumber function to validate the credit card number

	resp := Response{IsValid: isValid}                 // This creates a Response struct with the validation result
	w.Header().Set("Content-Type", "application/json") // This sets the response header to indicate that the response is JSON
	json.NewEncoder(w).Encode(resp)                    // This encodes the resp struct as JSON and writes it to the response writer
}

func isValidCreditCardNumber(cardNumber string) bool { // This function implements the Luhn algorithm for credit card validation
	numbers := make([]int, len(cardNumber)) // This creates a slice of integers with the length of the card number
	for i, r := range cardNumber {          // This loop iterates over the characters in the card number
		numbers[i] = int(r - '0') // This converts the character to an integer and stores it in the numbers slice
	}

	for i := len(numbers) - 2; i >= 0; i -= 2 { // This loop starts from the second-to-last digit and moves backwards by two digits
		numbers[i] *= 2     // This doubles the digit
		if numbers[i] > 9 { // This checks if the doubled digit is greater than 9
			numbers[i] -= 9 // If so, it subtracts 9 from the digit
		}
	}

	sum := 0                        // This variable will hold the sum of the digits
	for _, digit := range numbers { // This loop iterates over the digits in the numbers slice
		sum += digit // This adds each digit to the sum
	}

	return sum%10 == 0 // This returns true if the sum is divisible by 10 (indicating a valid credit card number), and false otherwise
}
