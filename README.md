# ft-quiz
A simple cli quiz game that posts to rest API

**This was initially written with go version 1.15.2 at the time.**

## Description 

This is a simple console application, where the user has the following options;

- Sign up 
- Log in 
- Exit
- Play 
- Compare scores

Each of these actions is interacted with via a console application where the user is prompted to enter text.
Upon successful interaction, the data is posted to rest endpoints which will at as simple CRUD using in-memory 
data to store retrieve or update data from users(players)

**REST enpoints include;**

**GET:**
- / - Just a blank page with a message.
- /players - displays list of players 

**POST:** 
- /new-player - adds new user to list
- /login - logs in existing player allowing them access to the game
- /logout - Signs player out taking them back to first menu
- /subimt-answers - Sets or re-sets the players answers to the questions and returns score 
- /compare-your-score - Shows how well or how bad the player did in comparison to others

 
## Gameplay Instructions
1. Start the game by going into a console and typing in `ft-quiz`
2. Press enter to start application (at this point server is running)
3. You must either use an existing account to log in or create a new account and log in with that.
4. Once logged you will have the option to play
5. Once the user clicks "play" they will be greeted with some instructions and they haev to press enter to start
6. They are greeted with a question with three possible answers and has to enter the answer key a, b, or c 
7. There are 5 questions and once the player has answered them all, at which poiint the program posts to the endpoint 
8. The enpoint will return the users evaluation on the quiz saying `You answered x out of 5 questions correctly!`
9. User is taken back the logged in menu where they are able to play again, compare their results or log out.

**Note: Exiting the application will also cause the server to shut down**

# Technical Description

The project started off as a boilerplate of the [spf13 Cobra project](https://github.com/spf13/cobra), 
quite specifically it started out using [Cobra Generator](https://github.com/spf13/cobra/blob/master/cobra/README.md) and abided by the general standards found online where the `Main.go` file simply initialzies the cobra command and the command file it's has the logic for it's purpose, in this case running the server and game.
