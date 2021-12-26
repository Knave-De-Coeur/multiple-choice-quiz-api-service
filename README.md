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

To build project, start api and run migrations : 

``` 
cd multiple-choice-quiz-api-service

docker-compose up
```
---

### Running application

Build console:

```
go build console
```

Then to simply run a new interface to play the game:

```
cd dir-to-binary/

console
```

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

Docker was used to create the mysql database and build an image of the api inside another container

Viper was used to manage configurations including fallbacks.

Whilst the api sets up the endpoints to communicate with the database it also creates the schema and handles migrations.

The cli is a simple interface that performs REST API requests to grab data. This is a minimalistic interface built from the ground up using the
[Cobra Generator]
(https://github.
com/spf13/cobra/blob/master/cobra/README.md)
to create the boilerplate for the cli.

