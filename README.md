# Go Text Modifier

Simple Go program to complete/edit/auto-correct text.

## Used Technologies

- Go (libraries below):
    - fmt
	- os
	- regexp
	- strconv
	- strings

## Installation & How to Use

1. Clone repository
2. Make sure two txt files (original txt file and the other txt file where the result should be written) are placed in the project directory. You can use `sample.txt` and `result.txt`, included already in the project.
3. Run `textModifier.go` with the original file name and the result file name.

    Ex:
    ```
    go run textModifier.go sample.txt result.txt    
    ```

## Spec Detail

- `(hex)` should replace the hexadecimal number before with the decimal version of the word (Ex: `1E (hex) files were added` -> `30 files were added`)

- `(bin)` should replace the binary number before with the decimal version of the word. (Ex: `It has been 10 (bin) years` -> `It has been 2 years`)

- `(up)` converts the word before with the Uppercase version of it. (Ex: `Ready, set, go (up) !` -> `Ready, set, GO!`)

- `(low)` converts the word before with the Lowercase version of it. (Ex: `I should stop SHOUTING (low)` -> `I should stop shouting`)

- `(cap)` converts the word before with the capitalized version of it. (Ex: `Welcome to the Brooklyn bridge (cap)` -> `Welcome To The Brooklyn Bridge`)

- For `(low)`, `(up)` and `(cap)`, if a number appears next to it (like `(low, 4)`) it turns the previously specified number of words in lowercase, uppercase or capitalized accordingly. (Ex: `This is so exciting (up, 2)` -> `This is SO EXCITING`)

- punctuations such as `.`, `,`, `!`, `?`, `:` and `;`, should be close to the previous word and with space apart from the next one. (Ex: `I was sitting over there ,and then BAMM !!` -> `I was sitting over there, and then BAMM!!`). However, if there are groups of punctuation like `...` and `!?`, the program should format the text as in the following example: `I was thinking ... You were right` -> `I was thinking... You were right`.
  
- The punctuation mark `'` must always be  with another `'`, and they should be placed to the right and left of the sentence in the middle of them, without any spaces. (Ex: `I am exactly how they describe me: ' awesome '` -> `I am exactly how they describe me: 'awesome'`, `As Elton John said: ' I am the most well-known homosexual in the world '` -> `As Elton John said: 'I am the most well-known homosexual in the world'`)

- Every instance of `a` should be turned into `an` if the next word begins with a vowel (`a`, `e`, `i`, `o`, `u`), or a `h`. (Ex: `There it was. A amazing rock!` -> `There it was. An amazing rock!`).