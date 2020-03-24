package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// Interactive updating.
var Interactive bool

// Direct packages only.
var Direct bool

func init() {
	rootCmd.AddCommand(updatesCmd)
	updatesCmd.Flags().BoolVarP(&Interactive, "interactive", "u", false, "Interactive update")
	updatesCmd.Flags().BoolVarP(&Direct, "direct", "d", false, "Direct packages")

}

// Update struct for the Outputs that have available updates.
type Update struct {
	Path    string
	Version string
	Time    string
}

// Output from go list.
type Output struct {
	Path     string
	Version  string
	Time     string
	Indirect bool
	Main     bool
	Update   Update
}

func decodeFormatOutput(out []byte) ([]Output, []string) {
	dec := json.NewDecoder(bytes.NewReader(out))

	var updates []Output
	var questions []string

	for {
		var o Output
		if err := dec.Decode(&o); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if o.Main {
			continue
		}

		if Direct && o.Indirect {
			continue
		}

		if o.Update.Time == "" {
			continue
		}

		updates = append(updates, o)
		questions = append(questions, fmt.Sprintf("%s %s -> %s", o.Path, o.Version, o.Update.Version))
	}
	if Verbose {
		fmt.Printf("Found %d packages\n", len(updates))
		fmt.Printf("%v\n", questions)
	}
	return updates, questions
}

func interactivity(updates []Output, questions []string) {
	var multiQs = []*survey.Question{
		{
			Prompt: &survey.MultiSelect{
				Message: "Select the packages you want to update",
				Options: questions,
			},
		},
	}

	selectedUpdates := []string{}

	if err := survey.Ask(multiQs, &selectedUpdates); err != nil {
		log.Fatal(err)
		return
	}

	for _, text := range selectedUpdates {
		packagePath := strings.Split(text, " v")[0]
		for _, update := range updates {
			if update.Path != packagePath {
				continue
			}

			fmt.Printf("Updating %s to %s\n", packagePath, update.Update.Version)

			execOutput, err := exec.Command("go", "get", "-v", fmt.Sprintf("%s@%s", update.Path, update.Update.Version)).Output()
			if err != nil {
				log.Fatal(err)
			}
			if Verbose && len(execOutput) > 0 {
				fmt.Println("OUTPUT:")
				fmt.Println(execOutput)
				fmt.Println("")
			}
		}
	}
}

func listUpdates(questions []string) {
	if len(questions) == 0 {
		fmt.Println("\nNo updates available!")
		return
	}

	fmt.Printf("\n")

	for _, update := range questions {
		fmt.Printf("* %s\n", update)
	}
}

var updatesCmd = &cobra.Command{
	Use:   "updates",
	Short: "Check for updates",
	Long:  `With or without interactivitiy`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Searching for package updates...")
		out, err := exec.Command("go", "list", "-u", "-m", "-json", "all").Output()
		if err != nil {
			log.Fatal(err)
		}

		updates, questions := decodeFormatOutput(out)

		if len(questions) == 0 {
			fmt.Println("\nNo updates available!")
			return
		}

		if Interactive {
			interactivity(updates, questions)
		} else {
			listUpdates(questions)
		}
	},
}
