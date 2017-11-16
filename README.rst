=========
flashtext
=========



What's it
------------
Flashtext is a simple and fast keyword extract tool in go.

It was inspired by the paper `Medium freeCodeCamp <https://medium.freecodecamp.org/regex-was-taking-5-days-flashtext-does-it-in-15-minutes-55f04411025f>`_
, here is the python implement:`<https://github.com/vi3k6i5/flashtext>`_



Installation
------------
::

    $ go get github.com/sundy-li/flashtext



Usage
-----
Extract keywords
::
    package main

    import (
        "fmt"

        "github.com/sundy-li/flashtext"
    )

    func main() {
        processor := flashtext.NewKeywordProcessor()
        // set the caseSensitive to false
        processor.SetConfig(false)

        processor.AddKeyword("Big Apple", "New York")
        processor.AddKeywordAndName("java", "Java")
        // set to find the longest keywords
        res := processor.Extracts("I like java, big big apple new york", true)
        fmt.Printf("res => %#v\n", res)
    }

To Remove keywords
::   
    processor.RemoveKeywords("New York")

Test
----
::

    $ go test


