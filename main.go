package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

var tasks []Task

func main() {
	if err := createFile(); err != nil {
		panic(err)
	}

	var t Task
	if os.Args[1] == "add" {
		t.Id = rand.Intn(10000)
		t.Description = strings.Split(os.Args[2], "=")[1]
		t.Status = "todo"
		t.CreatedAt = time.Now().Local().String()
		t.UpdatedAt = time.Now().Local().String()

		err := updateFile(t)
		if err != nil {
			fmt.Print(err.Error())
		}

		fmt.Printf("Task added successfully (ID: %v)", t.Id)
	} else if os.Args[1] == "update" {

	} else if os.Args[1] == "delete" {

	} else if os.Args[1] == "list" {
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

		fmt.Println("List all tasks")
		for _, task := range tasks {
			fmt.Printf("%v\t%v\t%v\t%v\t%v\n", task.Id, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
		}
	} else {
		fmt.Printf(`Command "%v" is not found`, os.Args[1])
	}

}

func createFile() error {
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

func updateFile(t Task) error {
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
