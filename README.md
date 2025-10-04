# Coral-CTC
[![GitHub](https://img.shields.io/badge/-GitHub-181717?style=for-the-badge&logo=GitHub&logoColor=white)](https://github.com/Lich-Corals/coral-ctc-terminal-calculator)
[![Coffee Logo](https://img.shields.io/badge/-Buy%20me%20a%20coffee-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=black)](https://www.coff.ee/lichcorals)

## CTC terminal calculator
CTC is a minimal and easy-to-use calculator application for your terminal.<br/>
This project is meant for people who are quick with terminal applications, who may want a replacement for a GUI calculator.

## Features
Supported operations are the following:
- addition
- subtraction
- multiplication
- division
- factorials
- powers
- roots
- modulo
- grouping with parentheses

## Installation
1. Download the binary `ctc` from the [latest release](https://github.com/Lich-Corals/coral-ctc-terminal-calculator/releases/latest).
2. Place the file in a useful location (e.g. `~/.local/bin/ctc`)
3. Make the file executable (e.g. `chmod +x ~/.local/bin/ctc`)
4. Add an alias to your shells configuration file:
Bash users can add the following line to their `.bashrc`:
```bash
alias ctc="~/.local/bin/ctc"
```

Fish users can add the following to their `~/.config/fish/fish.conf` file:
```fish
function ctc
    ~/.local/bin/ctc $argv
end
```
<br/>

The `ctc` command should be available in every newly launched terminal now.

> [!TIP]   
> You can use any path and any alias.
> `~/.local/bin/ctc` and `ctc` are just the recommended options.

> [!NOTE]   
> To compile the code yourself, clone the repository and run `go build` in the directory.

## Usage
### Basic syntax
The application runs with a single argument in quotes:
```bash
ctc "5 * 2 // 9 + 5.4 * 10"
```
This command takes the second root of 9 (`2 // 9`), multiplies it by 5 and adds 5.4 times 10 to it.

Every part of the calculation must be separated by a space.
The only exceptions are parentheses, which may be directly connected to a number (e.g. `(2 * 5)`).

### Unusual syntax
To keep it simple, CTC does not support functions like `n.pow()` `sqrt()`.
Instead, it has the `**` (power) and `//` (root) operators.
The syntax is inspired by the syntax of the English language.
Therefore, `2 ** 3` means _'2 to the 3'_ and `2 // 3` means _'the 2nd root of 3'_.

Additionally, there is the `%` (modulo) operator, which is used like _'x mod y'_.
This doesn't differ from the way it is implemented in many modern programming languages.

All other operations are the usual ones, as used in programming languages or other calculators.

### Priorities
The applications works from left to right and prioritises operations in the following order:
1. factorials
2. roots and powers
3. multiplication, division and modulo
4. addition and subtraction

The priorities can naturally be changed using parentheses.


## Updating
Currently, there is no way of getting notified by the application if an update is available.
Neither is this package available for any package manager.
<br/>
If you want to get notifications from GitHub, consider watching release activity for this repository.
<br/>
To update the program, repeat steps 1 to 3 from the installation instructions.

## Any problems?
Please [open an issue](https://github.com/Lich-Corals/coral-ctc-terminal-calculator/issues) to get help and to help making this product better!


## Contributing
You are welcome to contribute to this project in any way. Take a look at the [contribution guidelines](https://github.com/Lich-Corals/coral-ctc-terminal-calculator?tab=contributing-ov-file) for more information.
