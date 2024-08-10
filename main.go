package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"slices"
	"time"
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func main() {
	// when file is exist, file not create
	if err := createTaskJson(); err != nil {
		panic(err)
	}

	var tasks []Task
	var t Task

	if os.Args[1] == "add" {
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		description := addCmd.String("description", "", "type description")
		addCmd.Parse(os.Args[2:])

		t.Id = rand.Intn(10000)
		t.Description = *description
		t.Status = "todo"
		createdAt, err := time.Parse(time.DateTime, time.Now().GoString())
		if err != nil {
			panic(err)
		}
		t.CreatedAt = createdAt.String()

		updatedAt, err := time.Parse(time.DateTime, time.Now().GoString())
		if err != nil {
			panic(err)
		}
		t.UpdatedAt = updatedAt.String()

		err = addTaskToJson(t)
		if err != nil {
			fmt.Print(err.Error())
		}

		fmt.Printf("Task added successfully (ID: %v)", t.Id)
	} else if os.Args[1] == "update" {
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

		id := updateCmd.Int("id", 0, "type id")
		description := updateCmd.String("description", "", "type description")
		status := updateCmd.String("status", "", "type status")

		updateCmd.Parse(os.Args[2:])

		if *id == 0 {
			panic("id cannot empty")
		}

		tasks := openTaskJson()

		for i := 0; i < len(tasks); i++ {
			if tasks[i].Id == *id {
				if *description != "" {
					tasks[i].Description = *description
				}

				if *status != "" {
					tasks[i].Status = *status
				}

				tasks[i].UpdatedAt = time.Now().Local().String()
				break
			}
		}

		if err := updateTaskJson(tasks); err != nil {
			panic(err)
		}

		fmt.Printf("Task ID: %v is updated.", *id)
	} else if os.Args[1] == "delete" {
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

		id := deleteCmd.Int("id", 0, "type id")

		deleteCmd.Parse(os.Args[2:])

		tasks := openTaskJson()
		if *id != 0 {
			fmt.Println("Masuk")
			for i := 0; i < len(tasks); i++ {
				if tasks[i].Id == *id {
					tasks := slices.Delete(tasks, i, i+1)
					deleteTaskJson(tasks)
					break
				}
			}
		}

		fmt.Printf("Task ID: %v is deleted.", *id)
	} else if os.Args[1] == "list" {
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)

		status := flag.String("status", "", "type status")

		listCmd.Parse(os.Args[2:])

		jsonFile, err := os.Open("task.json")
		if err != nil {
			panic(err)
		}
		defer jsonFile.Close()

		bytes, err := io.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &tasks)
		if err != nil {
			panic(err)
		}

		fmt.Println("List all tasks")
		if *status == "todo" {
			for _, task := range tasks {
				if task.Status == "todo" {
					fmt.Printf("%v\t%v\t%v\t%v\t%v\n", task.Id, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
				}
			}
		} else if *status == "in-progress" {
			for _, task := range tasks {
				if task.Status == "in-progress" {
					fmt.Printf("%v\t%v\t%v\t%v\t%v\n", task.Id, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
				}
			}
		} else if *status == "done" {
			for _, task := range tasks {
				if task.Status == "done" {
					fmt.Printf("%v\t%v\t%v\t%v\t%v\n", task.Id, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
				}
			}
		} else {
			for _, task := range tasks {
				fmt.Printf("%v\t%v\t%v\t%v\t%v\n", task.Id, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
			}
		}
	} else {
		fmt.Printf(`Command "%v" is not found`, os.Args[1])
	}
}

func openTaskJson() []Task {
	var tasks []Task
	jsonFile, err := os.Open("task.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(bytes, &tasks)

	return tasks
}

func createTaskJson() error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filename := "task.json"
	_, err = os.Stat(path.Join(currDir, filename))
	if os.IsNotExist(err) {
		file, err := os.Create(path.Join(currDir, filename))
		if err != nil {
			return err
		}
		defer file.Close()
	}

	return nil
}

func addTaskToJson(t Task) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filename := "task.json"
	file, err := os.OpenFile(path.Join(currDir, filename), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	jsonTask := fmt.Sprintf(`,{
		"id": %v,
		"description": "%v",
		"status": "%v",
		"createdAt": "%v",
		"updatedAt": "%v"
	}
]`, t.Id, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
	_, err = file.WriteAt([]byte(jsonTask), stat.Size()-1)
	if err != nil {
		return err
	}

	return nil
}

func updateTaskJson(tasks []Task) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filename := "task.json"
	file, err := os.Create(path.Join(currDir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	file.Write(bytes)
	return nil
}

func deleteTaskJson(tasks []Task) error {
	err := updateTaskJson(tasks)
	if err != nil {
		return err
	}

	return nil
}
