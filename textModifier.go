package main

import (
	"fmt"
	"os"
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

	if argsLength == 3 {
		// check files exist
		for i := 1; i < argsLength; i++ {
			_, err := os.Open(os.Args[i])
			if err != nil {
				fmt.Printf("%v\n", err.Error())
				os.Exit(1)
			}
		}

		// if files exist
		original := os.Args[1]
		bytes, err := os.ReadFile(original)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}

		words := strings.Split(string(bytes), " ")

		// making new string to write in the result
		modified := []string{}
		i := 0
		for i < len(words) {
			fmt.Println(i)

			// avoid the error ( i+1 == len(words) )
			if i == len(words)-1 {
				modified = append(modified, words[i])
				i++

			} else {

				nextFirstLetter, nextRestLetter := getFirstLetter(words[i+1])
				fmt.Printf("check %v\n", words[i+1])

				switch true {

				case nextFirstLetter == '.' || nextFirstLetter == ',' || nextFirstLetter == '!' || nextFirstLetter == '?' || nextFirstLetter == ':' || nextFirstLetter == ';':
					fmt.Printf("First %v, rest %v\n", nextFirstLetter, nextRestLetter)
					modified = append(modified, words[i]+string(nextFirstLetter))
					modified = append(modified, nextRestLetter)
					i += 2

				// if "(hex)" found
				case strings.Contains(words[i+1], "(hex)"):
					if dec, err := strconv.ParseInt(words[i], 16, 64); err == nil {
						modified = append(modified, strconv.Itoa(int(dec)))
						i += 1
					}

				// if "(bin)" found
				case strings.Contains(words[i+1], "(bin)"):
					if dec, err := strconv.ParseInt(words[i], 10, 64); err == nil {
						modified = append(modified, strconv.Itoa(int(dec)))
						i += 1
					}

				case strings.Contains(words[i+1], "(low") || strings.Contains(words[i+1], "(up") || strings.Contains(words[i+1], "(cap"):
					count := 1
					extraSkip := 0 // use it to skip the number element to record

					// check if there's number after the sign
					if i+1 != len(words)-1 { // avoid access to out of range
						fmt.Printf("not a last word %v\n", words[i+1])
						if words[i+2][0] >= '1' && words[i+2][0] <= '9' {
							fmt.Printf("count found %v\n", words[i+2][0:len(words[i+2])-1])
							count, _ = strconv.Atoi(words[i+2][0 : len(words[i+2])-1])
							extraSkip = 1

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
							modified = append(modified, toTitle(words[i-count+1]))
						}
						count--
					}
					i = i + 1 + extraSkip

				default:
					modified = append(modified, words[i])
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

	} else { // not enough or too many arguments
		fmt.Println("Please specify files to read and write.")
	}
}
