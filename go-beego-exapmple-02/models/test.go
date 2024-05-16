package models

type Test struct {
	Response string 
}

func TestFunction() Test {
	 var test_var Test 
	 test_var.Response = "Hello World"
	 return test_var
}