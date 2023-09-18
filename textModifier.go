package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getFirstLetter(s string) (rune, string) {
	firstLetter := ' '
	rest := ""

	for i, rune := range s {
		if i == 0 {
			firstLetter = rune
		}
		rest += string(rune)
	}

	return firstLetter, rest
}

func toTitle(s string) string {
	result := ""

	for i, letter := range s {
		if i == 0 {
			result = string(letter - 32)
		} else {
			result += string(letter)
		}
	}
	fmt.Println("capitalized word", result)
	return result
}

func main() {
	args := os.Args
	argsLength := len(args)

	// if having correct arguments (number of files)
	if argsLength == 3 {
		// check files exist
		for i := 1; i < argsLength; i++ {
			_, err := os.Open(os.Args[i])

			// if file does not exist
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				os.Exit(1)
			}
		}

		original := os.Args[1]
		bytes, err := os.ReadFile(original)
		// in case unable to read the content
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}

		// making new []string to write in the result
		modified := []string{}

		// check each words[i+1] and treat words[i]
		words := strings.Split(string(bytes), " ")
		i := 0
		shouldSkip := 0 // count the words should not be recoreded (e.g: (hex), (low, 2))
		for i < len(words) {
			// add last words anyway
			if i == len(words)-1 {
				fmt.Printf("last word added: %v\n", words[i])
				modified = append(modified, words[i])
				i++

			} else {

				fmt.Printf("check: %v, skipcount: %v\n", words[i+1], shouldSkip)

				// for the case to check punctuation (case found:)
				puncts, _ := regexp.Compile(`^[.,!?;:]`)
				found := puncts.MatchString(words[i+1])

				switch true {

				case strings.Contains(words[i+1], "(hex)"):
					if dec, err := strconv.ParseInt(words[i], 16, 64); err == nil {
						modified = append(modified, strconv.Itoa(int(dec)))
						shouldSkip = 1
						i += 1
					}

				case strings.Contains(words[i+1], "(bin)"):
					if dec, err := strconv.ParseInt(words[i], 2, 64); err == nil {
						modified = append(modified, strconv.Itoa(int(dec)))
						shouldSkip = 1
						i += 1
					}

				case strings.Contains(words[i+1], "(low") || strings.Contains(words[i+1], "(up") || strings.Contains(words[i+1], "(cap"):
					count := 1 // how many words to modify
					shouldSkip = 1

					// check if there's number after the sign
					if i+1 != len(words)-1 { // avoid access to out of range
						word := regexp.MustCompile(`\d+(\.\d+)?`)
						numberFound := word.FindAllString(words[i+2], -1)
						if numberFound != nil {
							stringNumbers := ""
							for _, stringNumber := range numberFound {
								stringNumbers += stringNumber
							}
							count, _ = strconv.Atoi(stringNumbers)
							shouldSkip = 2

							// remove the element to be modified
							modified = modified[:len(modified)-count+1]
						}
					}

					for count > 0 {
						switch true {
						case strings.Contains(words[i+1], "(low"):
							modified = append(modified, strings.ToLower(words[i-count+1]))
						case strings.Contains(words[i+1], "(up"):
							modified = append(modified, strings.ToUpper(words[i-count+1]))
						case strings.Contains(words[i+1], "(cap"):
							modified = append(modified, toTitle(words[i-count+1])) // they say strings.Title is obsolete function :(
						}
						count--
					}
					i++
					break

				// check puctuations
				case found:
					// deal with ... or ?! differently
					punctsEx, _ := regexp.Compile(`^(\.{3}|\?!)`)
					punctsExFound := punctsEx.FindString(words[i+1])
					if punctsExFound != "" {
						modified = append(modified, words[i]+punctsExFound)
						fmt.Printf("word added: %v\n", words[i])
						shouldSkip = 1
					} else {
						punctsFound := puncts.FindString(words[i+1])
						modified = append(modified, words[i]+punctsFound)
						fmt.Printf("word added: %v\n", words[i])
						shouldSkip = 1
					}
					i++

				default:
					if shouldSkip == 0 {
						modified = append(modified, words[i])
						fmt.Printf("word added: %v\n", words[i])
					}
					if shouldSkip > 0 {
						shouldSkip--
					}
					i++
				}

			}
		}

		fmt.Println(strings.Join(modified, " "))

		// write the result to the file
		err = os.WriteFile(os.Args[2], []byte(strings.Join(modified, " ")), 0666)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}
		fmt.Println("The result is successfully written in the file.")

	} else { // not enough or too many arguments
		fmt.Println("Please specify files to read and write.")
	}
}
