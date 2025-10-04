__Everyone is welcome to contribute!__
If you want to contribute by modifying or adding code, please adhere to these rules:
* Write maintainable code
    * Keep you code understandable
        * Preferably by writing the code itself understandable
        * If that's not possible, add a few comments
    * Add short descriptions to added functions
* Ensure that your code is not affecting the functionality negatively
    * Your code should work correctly (i.e. as expected by the user)
    * Your code should introduce as few bugs and as few room for new bugs as possible
    * Your code has to pass `go test`
* Try to stick to the [Go-styleguide](https://google.github.io/styleguide/go/guide); especially to the naming
* Document new features
    * Document all changes in commit messages and pull-request comments
    * Add a note to the README if necessary
    * If it seems appropriate, please update the tests in `main_test.go`
* Clean up your code before submitting a pull request
    * Remove duplicated code, debug printing, dead variables, etc.
