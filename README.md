**Assignments_Tracker**

**Description:**
This Go program serves as an assignment tracker, by providing a functional command-line application to allow users to manage their assignments efficiently.  Users can add, list,  mark as complete, edit, and delete assignments, among other functionalities. By default any assignment added is sorted sorted by date into the list, however a sort by completion command is also present (functionality 9). The program also employs error handling to ensure to correct input and file handling giving the option to print the assignments to a file for easy reference. (functionality 10)

**Functionalities:**
1. **Add Assignment:** Users can input details such as assignment description, course, and due date to add a new assignment to the tracker.
2. **List Assignments:** Displays a list of all assignments along with their details such as course, name, due date, and completion status.
3. **Mark as Complete:** Allows users to mark an assignment as complete by specifying its index.
4. **Mark All as Complete:** Marks all assignments as complete in one go.
5. **Edit Assignment:** Enables users to modify the details of an existing assignment by providing its index and entering the updated information.
6. **Delete Assignment:** Allows users to delete a specific assignment by specifying its index.
7. **Delete All Assignments:** Clears the entire assignment list.
8. **Delete All Completed Assignments:** Removes all completed assignments, keeping only the incomplete ones.
9. **Sort Assignments by Completion Status:** Sorts the assignments based on their completion status, with incomplete assignments appearing first.
10. **Print Assignments to File:** Writes the list of assignments along with their details to a text file named "assignments.txt".
11. **Exit:** Terminates the program.

