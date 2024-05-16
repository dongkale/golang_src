/*
Copyright © 2023 Codoworks
Author:  Dexter Codo
Contact: dexter.codo@gmail.com
*/
package main

import "github.com/codoworks/go-boilerplate/cmd"

var VERSION string = "2.2.1-default"

func main() {
	cmd.Version = VERSION
	cmd.Execute()
}
