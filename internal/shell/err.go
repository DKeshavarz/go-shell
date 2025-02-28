package shell

import "errors"

var tooManyArgumentERR = errors.New("too many arguments")
var tooFewArgumentERR = errors.New("too few arguments")
