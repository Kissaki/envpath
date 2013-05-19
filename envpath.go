package main

import (
    "fmt"
    "os"
    "strings"
    "flag"
)

var countOnlyFlag bool
var listCurrentFlag bool
var listNoInvalidFlag bool

func init() {
	flag.BoolVar(&countOnlyFlag, "count", false, "Only count the current number of paths specified in the PATH environment variable and print them")
    flag.BoolVar(&listCurrentFlag, "list", false, "List the current environment paths")
    flag.BoolVar(&listNoInvalidFlag, "do-not-list-invalid", false, "Do not list invalid paths")
}

func main() {
    flag.Parse()
    
    if countOnlyFlag {
        paths := getPaths()
        fmt.Printf("%d", len(paths))
        return
    }

    fmt.Println("Checking PATH environment variable ...")
    paths := getPaths()
    fmt.Println(fmt.Sprintf("Found %d paths", len(paths)))
    // Create string array buffer and set invalidPaths to 
    invalidPaths := make([]string, 0)
    if listCurrentFlag {
        fmt.Println("Current paths:")
    }
    hasInvalidPaths := false
    for _, path := range paths {
        if listCurrentFlag {
            fmt.Println(path)
        }
        fileinfo, err := os.Stat(path)
        if os.IsNotExist(err) {
            hasInvalidPaths = true
            foundInvalidPath(path)
            invalidPaths = append(invalidPaths, path)
        } else {
            if !fileinfo.IsDir() {
                hasInvalidPaths = true
                foundInvalidPath(path)
                invalidPaths = append(invalidPaths, path)
            }
        }
    }
    fmt.Println("Invalid paths:")
    for _, path := range invalidPaths {
        fmt.Println(path)
    }
    if hasInvalidPaths {
        fmt.Println("Do you want to fix the invalid paths? (they will be removed from the PATH environment variable! [y]es/[n]o")
        var in string
        _, err := fmt.Scanf("%s", &in)
        if err != nil {
            fmt.Println("Error: Could not read user input: ", err)
        } else {
            if in == "y" {
                fmt.Println("Cleaning up PATH ...")
                
                var newpaths string
                for _, path := range paths {
                    fileinfo, err := os.Stat(path)
                    if os.IsNotExist(err) || !fileinfo.IsDir() {
                        fmt.Println("Dropping path ", path)
                    } else {
                        if len(newpaths) != 0 {
                            newpaths += string(os.PathListSeparator)
                        }
                        newpaths += path
                    }
                }
                fmt.Println("New PATH value: ", newpaths)
                /*err := os.Setenv("PATH", "abc")
                if err != nil {
                    fmt.Println("Error: Could not set environment variable: ", err)
                }*/
            }
        }
    }
    
    /*fmt.Println()
    dir, err := os.Getwd()
    if err != nil {
        fmt.Println("Could not get current executables active directory: ", err)
    } else {
        fmt.Println("current executables active directory: ", dir)
    }*/
    
}

func getPaths() (paths []string) {
    path := os.Getenv("PATH")
    return strings.Split(path, string(os.PathListSeparator))
}

func foundInvalidPath(path string) {
    if !listNoInvalidFlag {
        fmt.Println("Path does not exist: ", path)
    }
}
