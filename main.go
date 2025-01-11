package main

import (
    "bufio"
    "context"
    "fmt"
    "log"
    "os"
    "sort"
    "strconv"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Assignment struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `bson:"name"`
    Course   string             `bson:"course"`
    DueDate  time.Time          `bson:"due_date"`
    Complete bool               `bson:"complete"`
}

var client *mongo.Client
var collection *mongo.Collection
var assignments []Assignment

func connectToMongo() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        log.Fatal("Set your 'MONGODB_URI' environment variable.")
    }

    var err error
    clientOptions := options.Client().ApplyURI(uri)
    client, err = mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    collection = client.Database("school").Collection("assignments")
}

func addAssignmentToDB(assignment Assignment) {
    _, err := collection.InsertOne(context.Background(), assignment)
    if err != nil {
        log.Fatal(err)
    }
}

func listAssignmentsFromDB() []Assignment {
    var assignments []Assignment
    cursor, err := collection.Find(context.Background(), bson.D{})
    if err != nil {
        log.Fatal(err)
    }
    if err = cursor.All(context.Background(), &assignments); err != nil {
        log.Fatal(err)
    }
    return assignments
}

func markCompleteInDB(id primitive.ObjectID) {
    filter := bson.D{{"_id", id}}
    update := bson.D{{"$set", bson.D{{"complete", true}}}}
    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        log.Fatal(err)
    }
}

func deleteAssignmentFromDB(id primitive.ObjectID) {
    filter := bson.D{{"_id", id}}
    _, err := collection.DeleteOne(context.Background(), filter)
    if err != nil {
        log.Fatal(err)
    }
}

func addAssignment(name string, course string, dueDate time.Time) {
    newAssignment := Assignment{Name: name, Course: course, DueDate: dueDate, Complete: false}
    addAssignmentToDB(newAssignment)
    assignments = listAssignmentsFromDB()
    sortAssignments()

    fmt.Println("Assignment Added")
}

func listAssignments() {
    assignments = listAssignmentsFromDB()
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
    assignments = listAssignmentsFromDB()
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
        markCompleteInDB(assignments[index-1].ID)
        assignments = listAssignmentsFromDB()
        fmt.Println("Assignment Marked as Complete")
    } else {
        fmt.Println("Invalid Index")
    }
}

func markAllComplete() {
    for _, assignment := range assignments {
        markCompleteInDB(assignment.ID)
    }
    assignments = listAssignmentsFromDB()
    fmt.Println("All Assignments Marked as Complete")
}

func editAssignment(index int, name string, course string, dueDate time.Time) {
    if index >= 1 && index <= len(assignments) {
        assignment := assignments[index-1]
        filter := bson.D{{"_id", assignment.ID}}
        update := bson.D{{"$set", bson.D{{"name", name}, {"course", course}, {"due_date", dueDate}}}}
        _, err := collection.UpdateOne(context.Background(), filter, update)
        if err != nil {
            log.Fatal(err)
        }
        assignments = listAssignmentsFromDB()
        fmt.Println("Assignment Edited Successfully")
    } else {
        fmt.Println("Invalid Index")
    }
}

func deleteAssignment(index int) {
    if index >= 1 && index <= len(assignments) {
        deleteAssignmentFromDB(assignments[index-1].ID)
        assignments = listAssignmentsFromDB()
        fmt.Println("Assignment Deleted")
    } else {
        fmt.Println("Invalid Index")
    }
}

func deleteAllAssignments() {
    for _, assignment := range assignments {
        deleteAssignmentFromDB(assignment.ID)
    }
    assignments = listAssignmentsFromDB()
    fmt.Println("All Assignments Deleted")
}

func deleteAllCompletedAssignments() {
    for _, assignment := range assignments {
        if assignment.Complete {
            deleteAssignmentFromDB(assignment.ID)
        }
    }
    assignments = listAssignmentsFromDB()
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
    connectToMongo()
    assignments = listAssignmentsFromDB()
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


