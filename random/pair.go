package main

import(
)


var Pairs map[string]string

func MakePair(a, b string) {
    Pairs[a] = b
    Pairs[b] = a
}

