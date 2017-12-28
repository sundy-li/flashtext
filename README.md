## flashtext


#### What's it
`flashtext` is a simple and fast keyword extract tool in go. It was inspired by the paper [Medium freeCodeCamp](https://medium.freecodecamp.org/regex-was-taking-5-days-flashtext-does-it-in-15-minutes-55f04411025f) 

Here is the python implement:  [https://github.com/vi3k6i5/flashtext](https://github.com/vi3k6i5/flashtext)


#### Installation

```
    $ go get github.com/sundy-li/flashtext

```

#### Usage

- Extract keywords
```
    package main

    import (
        "fmt"

        "github.com/sundy-li/flashtext"
    )

    func main() {
        processor := flashtext.NewKeywordProcessor()
        // set the caseSensitive to false
        processor.SetCaseSenstive(false)
        processor.AddKeywords("I love go", "I like python")
        processor.AddKeywordAndName("java", "JavaEE")
        // set to find the longest keywords
        res := processor.ExtractKeywords("Hi, I love go, I like python and java", &flashtext.Option{Longest: true})
        for _, result := range res {
            fmt.Println(result.Keyword, "is met in the start position", result.StartIndex)
        }
    }
    //I love go is met in the start position 4
    //I like python is met in the start position 15
    //JavaEE is met in the start position 33
```
 

