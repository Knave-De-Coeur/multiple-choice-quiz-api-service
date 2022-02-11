package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"quiz-api-service/internal/config"
)

func main() {
	Execute()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "multiple-choice-quiz-api-service",
	Short: "My test for Fast Track",
	Long: `This is simple quiz where the user is presses ted with a couple questions
			and they have to select one from three to get the right answer.`,
	Args: func(cmd *cobra.Command, args []string) error {
		port := ""
		if len(args) < 1 {
			port = config.CurrentConfigs.Port
		} else {
			port = args[0]
		}
		checkAndAssignPort(port)
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		// runs game
		// go runGame()
	},
}

// This simply checks if the port is available and assigns it to the global variables
func checkAndAssignPort(port string) {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Printf("Can't listen on port %q: %s \n", port, err)
		os.Exit(1)
	}

	_ = ln.Close()

	config.CurrentConfigs.Port = port
	config.CurrentConfigs.Host = config.CurrentConfigs.Host + ":" + port + "/"
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("something went wrong: %+v", err)
	}
}

func init() {
	cobra.OnInitialize()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&config.CurrentConfigs.CfgFile, "config", "", "config file (default is $HOME/.multiple-choice-quiz-api-service.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// // This is one of the main goroutines of the application that runs the user interface part
func runGame() {
	fmt.Println("Welcome to Alex's quiz! Press enter to begin.")

	reader := bufio.NewReader(os.Stdin)
	key, _ := reader.ReadString('\n')

	if len(key) > 0 {
		for {
			fmt.Println("Enter option letter and press enter")

			// TODO: re-do flow
			// if CurrentUserID == 0 {
			// 	fmt.Println("a: Sign Up")
			// 	fmt.Println("b: Login")
			// 	fmt.Println("c: Exit")
			//
			// 	optionStr, _ := reader.ReadString('\n')
			//
			// 	option := []rune(optionStr)
			//
			// 	switch option[0] {
			// 	case 'a':
			// 		createUser(reader)
			// 	case 'b':
			// 		loginPrompt(reader)
			// 	case 'c':
			// 		os.Exit(0)
			// 	default:
			// 		fmt.Println("Invalid try again")
			// 	}
			// } else {
			// 	fmt.Println("a: Play")
			// 	fmt.Println("b: Logout")
			// 	fmt.Println("c: Compare")
			// 	fmt.Println("d: Exit")
			//
			// 	optionStr, _ := reader.ReadString('\n')
			//
			// 	option := []rune(optionStr)
			//
			// 	switch option[0] {
			// 	case 'a':
			// 		play(reader)
			// 	case 'b':
			// 		logoutPrompt()
			// 	case 'c':
			// 		compare()
			// 	case 'd':
			// 		os.Exit(0)
			// 	default:
			// 		fmt.Println("Invalid try again")
			// 	}
			// }

		}
	}
}
