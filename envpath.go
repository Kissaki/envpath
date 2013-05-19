/* Helperscript for managing the PATH environment variable.

License: 3-clause BSD License
*/
package main

import (
    "fmt"
    "os"
    "strings"
    "flag"
)

var printAboutOnlyFlag bool
var countOnlyFlag bool
var listCurrentFlag bool
var listNoInvalidFlag bool

func init() {
    flag.BoolVar(&printAboutOnlyFlag, "about", false, "About this executable")
    flag.BoolVar(&countOnlyFlag, "count", false, "Only count the current number of paths specified in the PATH environment variable and print them")
    flag.BoolVar(&listCurrentFlag, "list", false, "List the current environment paths")
    flag.BoolVar(&listNoInvalidFlag, "do-not-list-invalid", false, "Do not list invalid paths")
}

func main() {
    flag.Parse()
    
    if printAboutOnlyFlag {
        fmt.Println("Helperscript for managing the PATH environment variable.")
        fmt.Println("Original author: Jan Klass - aka Kissaki - http://kcode.de")
        return
    }
    
    if countOnlyFlag {
        paths := getPaths()
        fmt.Printf("%d", len(paths))
        return
    }

    fmt.Println("Checking PATH environment variable ...")

    paths := getPaths()
    fmt.Println(fmt.Sprintf("Found %d paths", len(paths)))

    validPaths := make([]string, 0)
    invalidPaths := make([]string, 0)
    if listCurrentFlag {
        fmt.Println("Current paths:")
    }
    for _, path := range paths {
        if listCurrentFlag {
            fmt.Println(path)
        }
        fileinfo, err := os.Stat(path)
        if os.IsNotExist(err) {
            foundInvalidPath(path)
            invalidPaths = append(invalidPaths, path)
        } else {
            if !fileinfo.IsDir() {
                foundInvalidPath(path)
                invalidPaths = append(invalidPaths, path)
            } else {
                validPaths = append(validPaths, path)
            }
        }
    }
    fmt.Println("Invalid paths:")
    for _, path := range invalidPaths {
        fmt.Println(path)
    }
    if len(invalidPaths) > 0 {
        if askPrintClean() {
            var newpaths string
            for _, path := range validPaths {
                if len(newpaths) != 0 {
                    newpaths += string(os.PathListSeparator)
                }
                newpaths += path
            }
            fmt.Println("Cleaned PATH value:")
            fmt.Println(newpaths)
        }
    }
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

// @return true for user input yes
func askPrintClean() bool {
    fmt.Println("Do you want the path value without the invalid paths? [y]es/[n]o")
    var in string
    fmt.Scanf("%s", &in)
    in = strings.ToLower(in)
    return (in == "y" || in == "yes")
    
}
