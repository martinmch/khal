package main

import (
    "fmt"
    "log"
    "os/exec"
    "strings"
    "time"
    "strconv"
)

func getNextMonth(t time.Time) string {
    // Determine the first of the next month after t.
    firstOfNextMonth := t.AddDate(0,1,-(t.Day()-1))
    // Get the string representation of the month number.
    return strconv.Itoa(int(firstOfNextMonth.Month()))
}

func getMonthAndDay(s string) string {
    if(len(s) < 7) {
        log.Fatal("getMonthAndDay called with too short string!")
    }
    return s[0:6]
}

func main() {
    // We get the output of cal. These will be printed on the left side, showing
    // an overview of the next two months.
    curMonth, err := exec.Command("cal").Output()
    if err != nil {
        log.Fatal(err)
    }

    // We get the output of cal for next month. This will be printed on the left
    // side, below the current month.
    nextMonthI := getNextMonth(time.Now())
    nextMonth, err := exec.Command("ncal", "-bMm " + nextMonthI).Output()
    if err != nil {
        log.Fatal(err)
    }

    // We get the tasks for today and tomorrow. These will be printed on the
    // right side.
    taskOutput, err := exec.Command("calendar", "-A 2").Output()
    if err != nil {
        log.Fatal(err)
    }
    // Now we compose the left hand side into a string array ready to be
    // printed.
    months := string(curMonth) + string(nextMonth)
    monthLines := strings.Split(strings.TrimSuffix(months, "\n"), "\n")
    monthLineCount := len(monthLines)

    tasks := string(taskOutput)
    taskLines := strings.Split(strings.TrimSuffix(tasks, "\n"), "\n")
    // The lines look like "Jan 18  XXXXXXX". We determine what the first
    // letters look like in the first line, "MON DD", and when this string
    // changes, we know that we've progressed to tomorrow.
    today := getMonthAndDay(taskLines[0])
    taskLinesf := []string{"Today:"}
    // We haven't added the header for tomorrow yet.
    missingHeader := true
    for i := range taskLines[1:] {
        md := getMonthAndDay(taskLines[i])
        // We only want to add the heading once, thus we use a boolean flag
        // to indicate, that we've added the header.
        if(md != today && missingHeader) {
            taskLinesf = append(taskLinesf, "")
            taskLinesf = append(taskLinesf, "Tomorrow:")
            missingHeader = false
        }
        // We skip the first 8 characters, since we're not interested in the
        // month and date, when showing the task.
        taskLinesf = append(taskLinesf, taskLines[i][8:])
    }

    // Ensure that there are at least as many task lines, as month lines.
    for len(taskLinesf) < monthLineCount {
        taskLinesf = append(taskLinesf, "")
    }

    // Since there are an equal amount of lines, we can use a single index to
    // print both arrays.
    for i, _ := range monthLines {
        fmt.Printf("%s    %s\n", monthLines[i], taskLinesf[i])
    }
}
