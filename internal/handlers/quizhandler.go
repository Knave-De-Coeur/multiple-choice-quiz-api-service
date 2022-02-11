package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"quiz-api-service/internal/api"
	"quiz-api-service/internal/config"
	"quiz-api-service/internal/pkg"
)

// PostToEndpoint Grouped logic that posts to endpoint and receives message to be outputted to the console
func PostToEndpoint(data interface{}, endpoint string) bool {
	requestJSON, _ := json.Marshal(data)

	res, err := http.Post(config.CurrentConfigs.Host+endpoint, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return false
	}

	defer res.Body.Close()

	var responseMessage api.Response

	_ = json.NewDecoder(res.Body).Decode(&responseMessage)

	fmt.Println(responseMessage.Message)

	return !responseMessage.Error
}

// This is the gameplay section, loops through each question, outputting the questions and possible answers
// Will wait for user input and go to the next question, once they're all answered it posts to get the result
func play(reader *bufio.Reader) bool {

	key, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	if len(key) > 0 {
		var listOfSubmittedRunes []rune
		// for _, q := range ListOfQuestions {
		// 	fmt.Printf("Question: %d %s", q.ID, q.Description)
		// 	questionKeyStrings := make([]string, 0, len(q.AnswerSelection))
		// 	for k := range q.AnswerSelection {
		// 		questionKeyStrings = append(questionKeyStrings, string(k))
		// 	}
		// 	sort.Strings(questionKeyStrings)
		// 	for {
		// 		for _, option := range questionKeyStrings {
		// 			runeKey := []rune(option)
		// 			fmt.Printf("%v: %v \n", option, q.AnswerSelection[runeKey[0]])
		// 		}
		// 		submittedAnswer, err := reader.ReadString('\n')
		// 		check(err)
		// 		if !isAnswerValid(questionKeyStrings, strings.TrimSpace(submittedAnswer)) {
		// 			fmt.Println("Please enter only one of the following options: ")
		// 			continue
		// 		}
		// 		answerRune := []rune(submittedAnswer)[0]
		// 		fmt.Println("answer submitted: " + string(answerRune))
		// 		listOfSubmittedRunes = append(listOfSubmittedRunes, answerRune)
		// 		break
		// 	}
		// }
		fmt.Println("Evaluating answers...")
		time.Sleep(2 * time.Second) // for dramatic suspense

		submitAnswerRequest := api.SubmitAnswersRequest{
			// UserID:           currentUser.ID,
			SubmittedAnswers: listOfSubmittedRunes,
		}

		PostToEndpoint(submitAnswerRequest, "submit-answer")

	}

	return true
}

// Console input and logic to set user to post to endpoint
func createUser(reader *bufio.Reader) bool {
	newUser := pkg.User{}

	fmt.Println("Start creating your profile.")

	for {
		fmt.Println("Enter Full Name: ")
		newUser.Name, _ = reader.ReadString('\n')

		fmt.Println("Enter your age: ")
		rawAge, _ := reader.ReadString('\n')
		rawAge = strings.TrimSpace(rawAge)
		intAge, _ := strconv.Atoi(rawAge)
		newUser.Age = int8(intAge)

		fmt.Println("Enter a Username: ")
		newUser.Username, _ = reader.ReadString('\n')
		fmt.Println("Enter a password: ")
		newUser.Password, _ = reader.ReadString('\n')

		newUser.Name = strings.TrimSpace(newUser.Name)
		newUser.Username = strings.TrimSpace(newUser.Username)
		newUser.Password = strings.TrimSpace(newUser.Password)

		if PostToEndpoint(newUser, "new-player") {
			break
		}

	}

	return true
}

// This will check the users input and set the CurrentUser
func loginPrompt(reader *bufio.Reader) bool {
	fmt.Println("Please enter Username:")
	inputtedUsername, _ := reader.ReadString('\n')

	fmt.Println("Please enter Password:")
	inputtedPassword, _ := reader.ReadString('\n')

	loginRequest := api.LoginRequest{
		Username: inputtedUsername,
		Password: inputtedPassword,
	}

	return PostToEndpoint(loginRequest, "login")
}

// Console function to post to log-out endpoint by taking the CurrentUserID
func logoutPrompt(uID uint) bool {
	logoutRequest := api.LogoutRequest{
		UserID: int(uID),
	}

	return PostToEndpoint(logoutRequest, "logout")
}

// This is one of the endpoint functions that stores the users submitted answers, then
// Sets the users score and response with a Response struct parsed into json
func submitAnswersAndGetResults(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := ioutil.ReadAll(req.Body)

	var submitAnswerRequest api.SubmitAnswersRequest
	_ = json.Unmarshal(reqBody, &submitAnswerRequest)

	// currentUser := searchUsersByID(submitAnswerRequest.UserID)
	// currentUser.SubmittedAnswers = submitAnswerRequest.SubmittedAnswers
	// currentUser.Score = 0

	// for i := range ListOfQuestions {
	// 	if ListOfQuestions[i].CorrectAnswer == currentUser.SubmittedAnswers[i] {
	// 		currentUser.Score++
	// 	}
	// }

	// updateUser(currentUser)

	// message := "You have answered " + strconv.Itoa(currentUser.Score) + " out of " + strconv.Itoa(len(ListOfQuestions)) + " questions correctly!"

	responseMessage := api.Response{
		// Message: message,
	}

	_ = json.NewEncoder(res).Encode(responseMessage)
}
