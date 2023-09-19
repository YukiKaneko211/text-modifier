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
		punctsOnly, _ := regexp.Compile(`^[.,!?:;']$|^((?:\.{3}|[\?!]{2}))$`)
		endWithPuncts, _ := regexp.Compile(`[\W]$|((?:\.{3}|[\?!]{2}))$`)
		startWithPuncts, _ := regexp.Compile(`^[.,!?:;]([\S]+)`)
		// validNum, _ := regexp.Compile(`^[\d]+[)]$`)
		anyNumber, _ := regexp.Compile(`\d+(\.\d+)?`)
		vowels, _ := regexp.Compile(`(?i)\b[aeiouh]\w*\b`)
		startWithMark, _ := regexp.Compile(`^[\']([\S]+)`)
		nonWords, _ := regexp.Compile(`^[^a-zA-Z]+$`)

		aChecked := []string{}
		for i := 0; i < len(words)-1; i++ {
			if strings.Compare(words[i], "a") == 0 && vowels.MatchString(words[i+1]) {
				aChecked = append(aChecked, "an")
			} else if strings.Compare(words[i], "A") == 0 && vowels.MatchString(words[i+1]) {
				aChecked = append(aChecked, "An")
			} else {
				aChecked = append(aChecked, words[i])
			}
		}
		aChecked = append(aChecked, words[len(words)-1])

		fmt.Println("A checked: ", strings.Join(aChecked, " "))

		for i < len(aChecked) {
			fmt.Printf("loop num: %v, check: %v, len(aChecked): %v, count %v\n", i, aChecked[i], len(aChecked), count)
			fmt.Println(strings.Join(modified, " "), " ,modifiedLength:", modifiedLength)
			switch true {
			case strings.Contains(aChecked[i], "(hex)"):
				if dec, err := strconv.ParseInt(aChecked[i-1], 16, 64); err == nil {
					modified = append(modified[:modifiedLength-1], strconv.Itoa(int(dec)))
					i++
				} else {
					fmt.Printf("Error: invalid word given to (hex) in the sample text: %v\n", err.Error())
					os.Exit(1)
				}
			case strings.Contains(aChecked[i], "(bin)"):
				if dec, err := strconv.ParseInt(aChecked[i-1], 2, 64); err == nil {
					modified = append(modified[:modifiedLength-1], strconv.Itoa(int(dec)))
					i++
				} else {
					fmt.Printf("Error: invalid word given to (bin) in the sample text: %v\n", err.Error())
					os.Exit(1)
				}

			case strings.Contains(aChecked[i], "(low") || strings.Contains(aChecked[i], "(up") || strings.Contains(aChecked[i], "(cap"):
				// check if there's valid number after the sign
				count := 0
				numberFound := ""
				if aChecked[i][len(aChecked[i])-1] == ')' {
					count = 1
				} else if i != len(aChecked)-1 && aChecked[i+1][len(aChecked[i+1])-1] == ')' { // avoid access to out of range
					numberFound = anyNumber.FindString(aChecked[i+1])
					if numberFound == "" {
						count = 1
					} else {
						count, _ = strconv.Atoi(numberFound)
					}
				} else {
					fmt.Println("Error: invalid syntax to convert cases in the sample text.")
					os.Exit(1)
				}

				// gives an error if there is not enough words to converted for given number
				if len(modified) < count {
					fmt.Println("Error: invalid numbers to convert cases in the sample text.")
					os.Exit(1)
				}

				// gives an error if there is invalid words to converted for given number
				for j := i - count; j < i; j++ {
					fmt.Println("BUG ", j, i, count, aChecked[j])
					isInvalid := nonWords.MatchString(aChecked[j])
					if isInvalid {
						fmt.Println("Error: invalid words to convert cases in the sample text.")
						os.Exit(1)
					}
				}

				modifiedLength -= count
				modified = modified[:modifiedLength]
				for j := count; j > 0; j-- {
					fmt.Printf("capitalized word %v\n", aChecked[i-j])
					fmt.Println(strings.Join(modified, " "), modifiedLength)
					fmt.Println("removed modified:", strings.Join(modified, " "))
					switch true {
					case strings.Contains(aChecked[i], "(low"):
						modified = append(modified, strings.ToLower(aChecked[i-j]))
					case strings.Contains(aChecked[i], "(up"):
						modified = append(modified, strings.ToUpper(aChecked[i-j]))
					case strings.Contains(aChecked[i], "(cap"):
						modified = append(modified, toTitle(aChecked[i-j])) // they say strings.Title is obsolete function :(
					}
					modifiedLength++
				}
				// skip if there was number
				if numberFound != "" {
					i++
				}
				i++

			case startWithMark.MatchString(aChecked[i]):
				openMark = true
				modified = append(modified, aChecked[i])
				modifiedLength++
				i++

			// ... & !?s and single punctuations
			case punctsOnly.MatchString(aChecked[i]):
				punctsOnlyFound := punctsOnly.FindString(aChecked[i])
				if punctsOnlyFound == "'" && !openMark { // only ' stick to the right letter
					fmt.Printf("StartMark Found")
					modified = append(modified, punctsOnlyFound+aChecked[i+1])
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
			case startWithPuncts.MatchString(aChecked[i]):
				startWithPunctsFound := startWithPuncts.FindString(aChecked[i])
				if !endWithPuncts.MatchString(modified[modifiedLength-1]) {
					fmt.Printf("Testing: add %v\n", modified[modifiedLength-1]+string(startWithPunctsFound[0]))
					modified = append(modified[:modifiedLength-1], modified[modifiedLength-1]+string(startWithPunctsFound[0]))
				}
				modified = append(modified, string(startWithPunctsFound[1:]))
				modifiedLength++
				i++
			default:
				fmt.Println("Default")
				modified = append(modified, aChecked[i])
				modifiedLength++
				i++
			}
			fmt.Println("aChecked checked: ", strings.Join(aChecked[:i], " "))
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
