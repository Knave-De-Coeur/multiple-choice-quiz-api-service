# Quiz-API-Server
A simple cli quiz game that posts to rest API

---

## Overview 

This is a simple console application and rest api server, where the user has the following options;

- Sign up 
- Log in 
- Exit
- Play 
- Compare scores

The user is prompted with these to being playing.
Upon successful interaction, the data is posted to rest endpoints which will at as simple CRUD using in-memory 
data to store retrieve or update data from users(players).

The game itself consists of only 5 questions and the user must answer them all in turn

**REST enpoints include;**

**GET:**
- **/** - Just a blank page with a message.
- **/players** - displays list of players 

**POST:** 
- **/new-player** - adds new user to list
- **/login** - logs in existing player allowing them access to the game
- **/logout** - Signs player out taking them back to first menu
- **/subimt-answers** - Sets or re-sets the players answers to the questions and returns score 
- **/compare-your-score** - Shows how well or how bad the player did in comparison to others

--- 

### Set up 

Once this has been placed in your `$GOROOT` dir simply run 

```
go install quiz-api-service
```

Which will create the binary on your machine, at which point you can follow the Run application section below.

---

### Running application

Run the application with:

```
quiz-api-service
```

or for a custom port:

```
quiz-api-service 9991
``` 

The latter might be in case, for some reason, the default port `(9990)` is unavailable at the time of launch.

### Gameplay

**Note: If you would like to skip the sign up you may use and existing user and log in with username: david54 password: pass765**

1. Press enter to start application (at this point server is running)
2. You must either use an existing account to log in or create a new account and log in with that.
3. Once logged you will have the option to play
4. Once the user clicks "play" they will be greeted with some instructions and they have to press enter to start
5. They are greeted with a question with three possible answers and has to enter the answer key a, b, or c 
6. There are 5 questions and once the player has answered them all, at which point the program posts to the endpoint 
7. The enpoint will return the users evaluation on the quiz saying `You answered x out of 5 questions correctly!`
8. User is taken back the logged in menu where they are able to play again, compare their results or log out.

**Note: Exiting the application will also cause the server to shut down**

---

## Technical Description

The project was built from the ground up using the [Cobra Generator](https://github.com/spf13/cobra/blob/master/cobra/README.md) to create the boilerplate for the cli. Standards for the cobra application were abided by where the `Main.go` file simply initialzies the cobra command and the command file it's has the logic for it's purpose, in this case running the server and game.

### Assumptions and considerations 

- I did not add anything special in relation to cobra's feature as i felt it was beyond the scope, simply used is as a boilerplate for the cli aspect of this app.
- Memory allocation was taken into account when setting up local and global variables, as well as scalability and maintainability.
- Error handling was implemented where appropriate as well as some custom validation for the requests, however not alot went into this as it might have gone beyond the scope of this task.
- No encryption was implemented on purpose since it was out of scope.
- An idea was to read the questions from an online csv, but as that might go a little far out of the cope of this task, I decide to simply hard code them 
- While there is a `Sign up` option in the main menu of the game, 3 users were added so that one could log in straight away with one of them and they all have their individual scores that are evalutaed on the `compare` option 
- Comments were added for all funcs and structs, structure was based on examples I saw on go projects
