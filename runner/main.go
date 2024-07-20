package main

import (
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var dfs func(node string)

// Parse the input file and return dependencies and progress maps.
func parseInput(filename string) (map[string][]string, map[string]float64, map[string][]string, error) {
	dependencies := make(map[string][]string) // skills // [A: [C,E]]
	progress := make(map[string]float64)      // progress for a particular skill // ["A": 0.8]
	mapOfChildWithParents := make(map[string][]string)
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "->") {
			parts := strings.Split(line, "->")
			parent := strings.TrimSpace(parts[0])
			child := strings.TrimSpace(parts[1])
			mapOfChildWithParents[child] = append(mapOfChildWithParents[child], parent)
			dependencies[parent] = append(dependencies[parent], child)
			// setting default value to 0 for both parent and child as we know that Progress may also undefined, in which case it should be treated as 0.
			if _, exists := progress[child]; !exists { // incase no progess are provided
				progress[child] = 0
			}
			if _, exists := progress[parent]; !exists {
				progress[parent] = 0
			}
		} else if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			skill := strings.TrimSpace(parts[0])
			prog, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err != nil {
				return nil, nil, nil, errors.New("error while ParseFloat")
			}
			if _, exists := dependencies[skill]; !exists { // incase progess is given but no edges are provided
				dependencies[skill] = append(dependencies[skill], skill)
				progress[skill] = 0
			}
			progress[skill] = prog
		}
	}
	return dependencies, progress, mapOfChildWithParents, nil
}

// Perform topological sort with progress-based ordering.
func topologicalSort(dependencies map[string][]string, progress map[string]float64, mapOfChildWithParents map[string][]string) []string {
	visited := make(map[string]bool)
	stack := []string{}
	dfs = func(node string) {
		visited[node] = true
		for _, child := range dependencies[node] {
			if !visited[child] {
				dfs(child)
			}
		}
		stack = append(stack, node)
	}

	// Run DFS from all unvisited nodes(parent nodes)
	for parent := range dependencies {
		if !visited[parent] {
			dfs(parent)
		}
	}

	// Reverse the stack to get the topological order
	reverse(stack)
	// If two or more skills are interchangeable in the roadmap
	// Identify nodes with multiple parents
	// Assign levels based on dependencies
	level := make(map[string]int)
	levelMap := make(map[int][]string)
	visitedLevels := make(map[string]bool)

	for _, node := range stack {
		maxLevel := 0
		for _, parent := range mapOfChildWithParents[node] {
			if !visitedLevels[parent] {
				continue
			}
			if level[parent] > maxLevel {
				maxLevel = level[parent]
			}
		}
		level[node] = maxLevel + 1
		visitedLevels[node] = true
		levelMap[level[node]] = append(levelMap[level[node]], node)
	}

	// Sort each level by progress
	var result []string
	for i := 0; i <= len(levelMap); i++ {
		if nodes, ok := levelMap[i]; ok {
			nodes = sortByProgress(nodes, progress)
			result = append(result, nodes...)
		}
	}

	return result
}

// Custom sorting by progress
func sortByProgress(skills []string, progress map[string]float64) []string {
	sort.SliceStable(skills, func(i, j int) bool {
		return progress[skills[i]] > progress[skills[j]]
	})
	return skills
}

// Reverse a slice of strings
func reverse(slice []string) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: go run main.go <input file>")
		return
	}
	filename := os.Args[1]
	dependencies, progress, mapOfChildWithParents, err := parseInput(filename)
	if err != nil {
		log.Println("Error reading input file:", err)
		return
	}

	skills := topologicalSort(dependencies, progress, mapOfChildWithParents)
	for _, skill := range skills {
		log.Println(skill)
	}
}
