package a

import "log"

func connectToDatabase() error {
	return nil
}

func main() {
	// connect to the database // want "comment parrots the code"
	err := connectToDatabase()

	// if error return // want "comment parrots the code"
	if err != nil {
		log.Fatal(err)
	}

	// This explains WHY we retry: network can be flaky in k8s
	err = connectToDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// fetch user from database // want "comment parrots the code"
	user := fetchUserFromDatabase()
	_ = user
}

func fetchUserFromDatabase() string {
	return "alice"
}

// processOrder handles order processing // want "comment parrots the code"
func processOrder() {
}

// validateInput ensures business rules are met before processing
// This is a policy decision: we reject orders under $10 to avoid payment fees
func validateInput() {
}
