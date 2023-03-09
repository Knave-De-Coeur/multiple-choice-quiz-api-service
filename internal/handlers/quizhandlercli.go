package handlers

// Console input and logic to set user to post to endpoint
// func createUser(reader *bufio.Reader) bool {
// 	newUser := pkg.User{}
//
// 	fmt.Println("Start creating your profile.")
//
// 	for {
// 		fmt.Println("Enter Full Name: ")
// 		newUser.Name, _ = reader.ReadString('\n')
//
// 		fmt.Println("Enter your age: ")
// 		rawAge, _ := reader.ReadString('\n')
// 		rawAge = strings.TrimSpace(rawAge)
// 		intAge, _ := strconv.Atoi(rawAge)
// 		newUser.Age = int8(intAge)
//
// 		fmt.Println("Enter a Username: ")
// 		newUser.Username, _ = reader.ReadString('\n')
// 		fmt.Println("Enter a password: ")
// 		newUser.Password, _ = reader.ReadString('\n')
//
// 		newUser.Name = strings.TrimSpace(newUser.Name)
// 		newUser.Username = strings.TrimSpace(newUser.Username)
// 		newUser.Password = strings.TrimSpace(newUser.Password)
//
// 		if PostToEndpoint(newUser, "new-player") {
// 			break
// 		}
//
// 	}
//
// 	return true
// }

// This will check the users input and set the CurrentUser
// func loginPrompt(reader *bufio.Reader) bool {
// 	fmt.Println("Please enter Username:")
// 	inputtedUsername, _ := reader.ReadString('\n')
//
// 	fmt.Println("Please enter Password:")
// 	inputtedPassword, _ := reader.ReadString('\n')
//
// 	loginRequest := api.LoginRequest{
// 		Username: inputtedUsername,
// 		Password: inputtedPassword,
// 	}
//
// 	return PostToEndpoint(loginRequest, "login")
// }

// This is the gameplay section, loops through each question, outputting the questions and possible answers
// Will wait for user input and go to the next question, once they're all answered it posts to get the result
// func play(reader *bufio.Reader) bool {
//
// 	key, err := reader.ReadString('\n')
// 	if err != nil {
// 		return false
// 	}
//
// 	if len(key) > 0 {
// 		var listOfSubmittedRunes []rune
// 		// for _, q := range ListOfQuestions {
// 		// 	fmt.Printf("Question: %d %s", q.ID, q.Description)
// 		// 	questionKeyStrings := make([]string, 0, len(q.AnswerSelection))
// 		// 	for k := range q.AnswerSelection {
// 		// 		questionKeyStrings = append(questionKeyStrings, string(k))
// 		// 	}
// 		// 	sort.Strings(questionKeyStrings)
// 		// 	for {
// 		// 		for _, option := range questionKeyStrings {
// 		// 			runeKey := []rune(option)
// 		// 			fmt.Printf("%v: %v \n", option, q.AnswerSelection[runeKey[0]])
// 		// 		}
// 		// 		submittedAnswer, err := reader.ReadString('\n')
// 		// 		check(err)
// 		// 		if !isAnswerValid(questionKeyStrings, strings.TrimSpace(submittedAnswer)) {
// 		// 			fmt.Println("Please enter only one of the following options: ")
// 		// 			continue
// 		// 		}
// 		// 		answerRune := []rune(submittedAnswer)[0]
// 		// 		fmt.Println("answer submitted: " + string(answerRune))
// 		// 		listOfSubmittedRunes = append(listOfSubmittedRunes, answerRune)
// 		// 		break
// 		// 	}
// 		// }
// 		fmt.Println("Evaluating answers...")
// 		time.Sleep(2 * time.Second) // for dramatic suspense
//
// 		submitAnswerRequest := api.SubmitAnswersRequest{
// 			// UserID:           currentUser.ID,
// 			SubmittedAnswers: listOfSubmittedRunes,
// 		}
//
// 		PostToEndpoint(submitAnswerRequest, "submit-answer")
//
// 	}
//
// 	return true
// }
