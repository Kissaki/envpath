/* Helperscript for managing the PATH environment variable.

License: 3-clause BSD License
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var printAboutOnlyFlag bool
var onlyCountFlag bool
var onlyListCurrentFlag bool
var onlyListInvalidFlag bool
var countFlag bool
var listCurrentFlag bool
var listInvalidFlag bool

func init() {
	flag.BoolVar(&printAboutOnlyFlag, "about", false, "About this executable")
	flag.BoolVar(&onlyCountFlag, "only-count", false, "Only print number of paths in the PATH environment variable")
	flag.BoolVar(&onlyListCurrentFlag, "only-list-valid", false, "Only print valid paths in the PATH environment variable, one per line")
	flag.BoolVar(&onlyListInvalidFlag, "only-list-invalid", false, "Only print invalid paths in the PATH environment variable, one per line")

	flag.BoolVar(&countFlag, "count", true, "Count and print the number of paths specified in the PATH environment variable")
	flag.BoolVar(&listCurrentFlag, "list", false, "List the paths in the PATH environment variable")
	flag.BoolVar(&listInvalidFlag, "list-invalid", true, "List the invalid paths in the PATH environment variable")
}

func main() {
	flag.Parse()

	if printAboutOnlyFlag {
		fmt.Println("Helperscript for managing the PATH environment variable.")
		fmt.Println("Original author: Jan Klass - aka Kissaki - http://kcode.de")
		return
	}

	if onlyCountFlag {
		fmt.Printf("%d", len(getPaths()))
		return
	}

	if onlyListCurrentFlag {
		validPaths, _ := getSplitPaths()
		printListNoEndl(validPaths)
		return
	}

	if onlyListInvalidFlag {
		_, invalidPaths := getSplitPaths()
		printListNoEndl(invalidPaths)
		return
	}

	if countFlag {
		fmt.Println(fmt.Sprintf("Found %d paths", len(getPaths())))
	}

	validPaths, invalidPaths := getSplitPaths()
	if listCurrentFlag {
		fmt.Println("Current paths:")
		printList(validPaths)
	}
	if listInvalidFlag {
		fmt.Println("Invalid paths:")
		printList(invalidPaths)
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

func getSplitPaths() (validPaths []string, invalidPaths []string) {
	paths := getPaths()
	validPaths = make([]string, 0)
	invalidPaths = make([]string, 0)
	for _, path := range paths {
		fileinfo, err := os.Stat(path)
		if os.IsNotExist(err) || !fileinfo.IsDir() {
			invalidPaths = append(invalidPaths, path)
		} else {
			validPaths = append(validPaths, path)
		}
	}
	return
}

func printList(list []string) {
	for _, path := range list {
		fmt.Println(path)
	}
}

func printListNoEndl(list []string) {
	for i, path := range list {
		if i != 0 {
			fmt.Print("\n")
		}
		fmt.Print(path)
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
