package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type Assignment struct {
	Name     string
	Course   string
	DueDate  time.Time
	Complete bool
}

var assignments []Assignment

func addAssignment(name string, course string, dueDate time.Time) {
	newAssignment := Assignment{Name: name, Course: course, DueDate: dueDate, Complete: false}
	assignments = append(assignments, newAssignment)
	sortAssignments()

	fmt.Println("Assignment Added")
}

func listAssignments() {
	fmt.Println("\nAssignments-To-Do\n")
	for i, assignment := range assignments {
		var status string = "Incomplete"
		if assignment.Complete {
			status = "Complete"
		}

		fmt.Printf("%d.\nCourse: %s\n%s\nDue date: %s [%s]\n\n", i+1, assignment.Course, assignment.Name, assignment.DueDate.Format("2006-01-02"), status)
	}
}

func printAssignmentsToFile() {
	file, err := os.OpenFile("assignments.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("Error creating file")
		return
	}
	defer file.Close()

	fmt.Fprintln(file, "Assignments TO-DO") 

	for i, assignment := range assignments {
		var status string = "Incomplete"
		if assignment.Complete {
			status = "Complete"
		}

		fmt.Fprintf(file, "%d.\nCourse: %s\n%s\nDue date: %s [%s]\n\n", i+1, assignment.Course, assignment.Name, assignment.DueDate.Format("2006-01-02"), status)
	}

	fmt.Println("Assignments printed to assignments.txt")
}

func markComplete(index int) {
	if index >= 1 && index <= len(assignments) {
		assignments[index-1].Complete = true
		fmt.Println("Assignment Marked as Complete")
	} else {
		fmt.Println("Invalid Index")
	}
}

func markAllComplete() {
	for i := range assignments {
		assignments[i].Complete = true
	}
	fmt.Println("All Assignments Marked as Complete")
}

func editAssignment(index int, name string, course string, dueDate time.Time) {
	if index >= 1 && index <= len(assignments) {
		assignments[index-1].Name = name
		assignments[index-1].Course = course
		assignments[index-1].DueDate = dueDate
		sortAssignments()
		fmt.Println("Assignment Edited Successfully")
	} else {
		fmt.Println("Invalid Index")
	}
}

func deleteAssignment(index int) {
	if index >= 1 && index <= len(assignments) {
		assignments = append(assignments[:index-1], assignments[index:]...)
		fmt.Println("Assignment Deleted")
	} else {
		fmt.Println("Invalid Index")
	}
}

func deleteAllAssignments() {
	assignments = []Assignment{}
	fmt.Println("All Assignments Deleted")
}

func deleteAllCompletedAssignments() {
	var incompleteAssignments []Assignment
	for _, assignment := range assignments {
		if !assignment.Complete {
			incompleteAssignments = append(incompleteAssignments, assignment)
		}
	}
	assignments = incompleteAssignments
	fmt.Println("All Completed Assignments Deleted")
}

func sortAssignments() {
	sort.Slice(assignments, func(i, j int) bool {
		return assignments[i].DueDate.Before(assignments[j].DueDate)
	})
}

func sortAssignmentsByCompletion() {
	sort.Slice(assignments, func(i, j int) bool {
		return assignments[i].Complete && !assignments[j].Complete
	})
	fmt.Println("Assignments Sorted by Completion Status")
}

func printOptions() {
	fmt.Println("\nOptions")
	fmt.Println("1. Add Assignment")
	fmt.Println("2. List Assignments")
	fmt.Println("3. Mark as Complete")
	fmt.Println("4. Mark All as Complete")
	fmt.Println("5. Edit Assignment")
	fmt.Println("6. Delete Assignment")
	fmt.Println("7. Delete All Assignments")
	fmt.Println("8. Delete All Completed Assignments")
	fmt.Println("9. Sort Assignments by Completion Status")
	fmt.Println("10. Print Assignments to File")
	fmt.Println("11. Exit")
}

func main() {
	var indexInput int
	var nameInput, courseInput, dueDateInput string
	printOptions() 
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter Option Required 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11: ")
		scanner.Scan()
		input := scanner.Text()

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid Choice")
			continue
		}
		switch choice {
		case 1:
			fmt.Print("Enter Assignment Description: ")
			scanner.Scan()
			nameInput = scanner.Text()

			fmt.Print("Enter Course: ")
			scanner.Scan()
			courseInput = scanner.Text()

			fmt.Print("Enter Due Date (YYYY-MM-DD): ")
			scanner.Scan()
			dueDateInput = scanner.Text()
			dueDate, err := time.Parse("2006-01-02", dueDateInput)
			if err != nil {
				fmt.Println("Invalid Date. Please try again.")
				continue
			}

			addAssignment(nameInput, courseInput, dueDate)
		case 2:
			listAssignments()
		case 3:
			fmt.Print("Enter Index: ")
			scanner.Scan()
			indexInput, _ = strconv.Atoi(scanner.Text())
			markComplete(indexInput)
		case 4:
			markAllComplete()
		case 5:
			fmt.Print("Enter Index: ")
			scanner.Scan()
			indexInput, _ = strconv.Atoi(scanner.Text())

			fmt.Print("Enter New Assignment Name: ")
			scanner.Scan()
			nameInput = scanner.Text()

			fmt.Print("Enter New Course: ")
			scanner.Scan()
			courseInput = scanner.Text()

			fmt.Print("Enter New Due Date (YYYY-MM-DD): ")
			scanner.Scan()
			dueDateInput = scanner.Text()
			dueDate, err := time.Parse("2006-01-02", dueDateInput)
			if err != nil {
				fmt.Println("Invalid Date. Please try again.")
				continue
			}

			editAssignment(indexInput, nameInput, courseInput, dueDate)
		case 6:
			fmt.Print("Enter Index: ")
			scanner.Scan()
			indexInput, _ = strconv.Atoi(scanner.Text())
			deleteAssignment(indexInput)
		case 7:
			deleteAllAssignments()
		case 8:
			deleteAllCompletedAssignments()
		case 9:
			sortAssignmentsByCompletion()
		case 10:
			printAssignmentsToFile()
		case 11:
			os.Exit(0)
		default:
			fmt.Println("Invalid Choice")
		}
	}
}
