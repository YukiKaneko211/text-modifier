package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func toTitle(s string) string {
	result := ""

	for i, letter := range s {
		if i == 0 {
			result = string(letter - 32)
		} else {
			result += string(letter)
		}
	}
	fmt.Println("capped word", result)
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
		modifiedLength := 0

		words := strings.Split(string(bytes), " ")
		i := 0
		count := 0
		openMark := false

		// to check the key letters in the loop
		punctsOnly, _ := regexp.Compile(`^[\W]$|^((?:\.{3}|[\?!]{2}))$`)
		endWithPuncts, _ := regexp.Compile(`[\W]$|((?:\.{3}|[\?!]{2}))$`)
		startWithPuncts, _ := regexp.Compile(`^\W+([\S]+)`)
		anyNumber := regexp.MustCompile(`\d+(\.\d+)?`)
		vowels, _ := regexp.Compile(`(?i)\b[aeiou]\w*\b`)

		for i < len(words) {
			fmt.Printf("loop num: %v, check: %v, len(words): %v, count %v\n", i, words[i], len(words), count)
			fmt.Println(strings.Join(modified, " "), " ,modifiedLength:", modifiedLength)

			// add last words

			switch true {

			case strings.Contains(words[i], "(hex)"):
				if dec, err := strconv.ParseInt(words[i-1], 16, 64); err == nil {
					modified = append(modified[:modifiedLength-1], strconv.Itoa(int(dec)))
					i++
				}

			case strings.Contains(words[i], "(bin)"):
				if dec, err := strconv.ParseInt(words[i-1], 2, 64); err == nil {
					modified = append(modified[:modifiedLength-1], strconv.Itoa(int(dec)))
					i++
				}

			case strings.Contains(words[i], "(low") || strings.Contains(words[i], "(up") || strings.Contains(words[i], "(cap"):

				// check if there's number after the sign
				if i != len(words)-1 { // avoid access to out of range
					numberFound := anyNumber.FindString(words[i+1])
					count := 1
					if numberFound == "" {
						count = 1
					} else {
						count, _ = strconv.Atoi(numberFound)
					}
					fmt.Printf("count: %v\n", count)
					modifiedLength -= count
					modified = modified[:modifiedLength]
					for j := count; j > 0; j-- {
						fmt.Printf("capitalized word %v\n", words[i-j])
						fmt.Println(strings.Join(modified, " "), modifiedLength)

						fmt.Println("removed modified:", strings.Join(modified, " "))
						switch true {
						case strings.Contains(words[i], "(low"):
							modified = append(modified, strings.ToLower(words[i-j]))
						case strings.Contains(words[i], "(up"):
							modified = append(modified, strings.ToUpper(words[i-j]))
						case strings.Contains(words[i], "(cap"):
							modified = append(modified, toTitle(words[i-j])) // they say strings.Title is obsolete function :(
						}
						modifiedLength++
					}
					// skip if there was number
					if numberFound != "" {
						i++
					}
				}
				i++

			// ... & !?s and single punctuations
			case punctsOnly.MatchString(words[i]):
				fmt.Printf("Singlepunct Found: check %v\n", modified[modifiedLength-1])
				punctsOnlyFound := punctsOnly.FindString(words[i])
				if punctsOnlyFound == "'" && !openMark { // only ' stick to the right letter
					fmt.Printf("StartMark Found")
					modified = append(modified, punctsOnlyFound+words[i+1])
					modifiedLength++
					openMark = true
					i++
				} else { // stick to the left letter
					if punctsOnlyFound != string(modified[modifiedLength-1][len(modified[modifiedLength-1])-1]) {
						fmt.Printf("Testing: add %v\n", modified[modifiedLength-1]+punctsOnlyFound)
						modified = append(modified[:modifiedLength-1], modified[modifiedLength-1]+punctsOnlyFound)
						if punctsOnlyFound == "'" && openMark {
							openMark = false
						}
					}
				}
				i++

			case startWithPuncts.MatchString(words[i]):
				startWithPunctsFound := startWithPuncts.FindString(words[i])
				if !endWithPuncts.MatchString(modified[modifiedLength-1]) {
					fmt.Printf("Testing: add %v\n", modified[modifiedLength-1]+string(startWithPunctsFound[0]))
					modified = append(modified[:modifiedLength-1], modified[modifiedLength-1]+string(startWithPunctsFound[0]))
				}
				modified = append(modified, string(startWithPunctsFound[1:]))
				modifiedLength++
				i++

			case strings.Compare(words[i], "a") == 0 && vowels.MatchString(words[i+1]):
				modified = append(modified, "an")
				modifiedLength++
				fmt.Printf("word added: %v\n", "an")
				i++

			default:
				fmt.Println("Default")
				modified = append(modified, words[i])
				modifiedLength++
				i++
			}
			fmt.Println("Words checked: ", strings.Join(words[:i], " "))
			fmt.Println("Modified RN: ", strings.Join(modified, " "))
		}

		fmt.Println("End of loop: ", strings.Join(modified, " "))

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
