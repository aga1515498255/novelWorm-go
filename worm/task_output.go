package worm

import "fmt"

type InputData struct {
	Id   string
	Json string
}

var OutputTasks []InputData

var Out chan []InputData

func init() {
	Out = make(chan []InputData)
}

func CollectTask() {

	for {
		for _, c := range tasks {
			select {
			case v := <-c.toClient:
				fmt.Println("read from chan:", v.Json)
				updateOutput(v)
			default:
			}

		}

		// fmt.Println("length of OutputTasks is ", len(OutputTasks))

		select {
		case Out <- OutputTasks:
			fmt.Println("write into http service")
		default:
		}

	}

}

func updateOutput(input InputData) {
	for i, v := range OutputTasks {

		if v.Id == input.Id {
			OutputTasks[i] = input
			return
		}
	}

	OutputTasks = append(OutputTasks, input)
}
