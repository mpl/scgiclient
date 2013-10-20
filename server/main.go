package main

import (
    "github.com/hoisie/web"
)

func hello(val string) string { 
    return "hello " + val 
} 

func main() {
    web.Get("/(.*)", hello)
    web.RunScgi("0.0.0.0:6580")
}
