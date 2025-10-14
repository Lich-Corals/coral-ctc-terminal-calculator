# Coral-CTC
[![GitHub](https://img.shields.io/badge/-GitHub-181717?style=for-the-badge&logo=GitHub&logoColor=white)](https://github.com/Lich-Corals/coral-ctc-terminal-calculator)
[![Sourceforge](https://img.shields.io/badge/-Sourceforge-FF6600?style=for-the-badge&logo=sourceforge&logoColor=white)](https://sourceforge.net/projects/coral-ctc-terminal-calculator/)
[![Coffee Logo](https://img.shields.io/badge/-Buy%20me%20a%20coffee-FFDD00?style=for-the-badge&logo=buymeacoffee&logoColor=black)](https://www.coff.ee/lichcorals)

<p align="center">
  <img src="https://github.com/user-attachments/assets/03021b04-a2d6-4ad6-b470-11343def681a">
</p>

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
- logarithms
- nPr and nCr
- Sine, cosine and tangent
- converting to absolute numbers
- A few constants
- grouping with parentheses
- A variable for the previous answer
- A command history

## Installation
1. Download the binary `ctc_xy` for your architecture and OS from the [latest release](https://github.com/Lich-Corals/coral-ctc-terminal-calculator/releases/latest).
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

> [!CAUTION]   
> MacOS and Windows binaries are included from release 0.4.0.
> These are not tested yet, because I don't have access to those operating systems.
> 
> The setup, especially for Windows will be different from what was said above.
> If you can confirm that CTC is running on Windows or MacOS, you're welcome to let me know in an issue!

> [!NOTE]   
> To compile the code yourself, clone the repository and run `go build` in the directory.

## Usage
### Basic syntax and usage
The application runs with a single argument in quotes:
```bash
ctc "5 * 2 // 9 + 5.4 * 10"
```
This command takes the second root of 9 (`2 // 9`), multiplies it by 5 and adds 5.4 times 10 to it.

CTC also has a continuous mode:
```bash
$ ctc
> 1 + 1
2
> exit
```

You can navigate using the arrow keys; `Up` and `Down` are used to access the command history.
To send a command, press the `Enter` key.
`Delete`-key for deletion of the cursor's character is supported.

The mode can be exited by sending a `exit`, `quit`, `:q`, `;q` and `exit()`.

Every part of the calculation must be separated by a space.
The only exceptions are parentheses, which may be directly connected to a number (e.g. `(2 * 5)`).

### Unusual syntax
To keep it simple, CTC does not support functions like `n.pow()` `sqrt()`.
Instead, it has the `**` (power) and `//` (root) operators.
The syntax is inspired by the syntax of the English language.
Therefore, `2 ** 3` means _'2 to the 3'_ and `2 // 3` means _'the 2nd root of 3'_.

Additionally, there is the `%` (modulo) operator, which is used like _'x mod y'_, and the `log` operator to get the logarithm of _x to the base y_ (`x log y`).
Calculations like `sin x` or `dsin x` are also supported, where `sin` works with radians and `dsin` with degrees.
N.b. the same works with the `cos` and `tan` functions.

> [!NOTE]   
> Promts like `cos tan xy` are currently not supported.   
> Use `cos (tan xy)` instead.

The `nPr` (permutations) and `nCr` (combinations) functions work like on most calculators; for example, `n nPr r` (or `x nPr y`) would be the following:
> 
> n! / (n - r)!
>

All other operations are the usual ones, as used in programming languages or other calculators.

### Priorities
The applications works from left to right and prioritises operations in the following order:
1. factorials
2. roots, powers, absolute-functions, sine, cosine, etc.
3. multiplication, division, nPr, nCr and modulo
4. addition, subtraction and logarithms

The priorities can naturally be changed using parentheses.

### Constants
The following constants are supported
- pi
- tau
- phi
- e (Euler's number)
- g (gravity)
- c (speed of light)

The negative of every constant `x` is available as `-x`. 

> [!NOTE]   
> Constants `g` and `c` use SI units.

### Special values
The "constant" `ans` can be used to insert the previous answer.

## Updating
Currently, there is no way of getting notified by the application if an update is available.
Neither is this package available for any package manager.
<br/>
If you want to get notifications from GitHub, consider watching release activity for this repository.
It is also possible to enable e-mail notifications on SourceForge.
<br/>
To update the program, repeat steps 1 to 3 from the installation instructions.

## Any problems?
Please [open an issue](https://github.com/Lich-Corals/coral-ctc-terminal-calculator/issues) to get help and to help making this product better!


## Contributing
You are welcome to contribute to this project in any way. Take a look at the [contribution guidelines](https://github.com/Lich-Corals/coral-ctc-terminal-calculator?tab=contributing-ov-file) for more information.
