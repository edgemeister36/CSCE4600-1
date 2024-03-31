package main

import (
	"fmt"
	"io"
	"sort"
)

type (
	Process struct {
		ProcessID     string
		ArrivalTime   int64
		BurstDuration int64
		Priority      int64
		RemainingTime int64
		Index         int64
	}
	TimeSlice struct {
		PID   string
		Start int64
		Stop  int64
	}
)

// Additional code for priority queue management

type PriorityQueue []*Process

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].BurstDuration < pq[j].BurstDuration
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Process))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
		}
		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFSchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	// Sort processes by burst duration (Shortest Job First)
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].BurstDuration < processes[j].BurstDuration
	})

	// Process the sorted processes
	for i := range processes {
		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
			if waitingTime < 0 {
				waitingTime = 0
			}
		}

		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
	var (
		serviceTime     int64
		totalWait       float64
		totalTurnaround float64
		lastCompletion  float64
		waitingTime     int64
		schedule        = make([][]string, len(processes))
		gantt           = make([]TimeSlice, 0)
	)

	// Sort processes by burst duration (Shortest Job First)
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].BurstDuration < processes[j].BurstDuration
	})

	// Process the sorted processes
	for i := range processes {
		if processes[i].ArrivalTime > serviceTime {
			serviceTime = processes[i].ArrivalTime
		}

		if processes[i].ArrivalTime > 0 {
			waitingTime = serviceTime - processes[i].ArrivalTime
			if waitingTime < 0 {
				waitingTime = 0
			}
		}

		totalWait += float64(waitingTime)

		start := waitingTime + processes[i].ArrivalTime

		turnaround := processes[i].BurstDuration + waitingTime
		totalTurnaround += float64(turnaround)

		completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
		lastCompletion = float64(completion)

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime),
			fmt.Sprint(turnaround),
			fmt.Sprint(completion),
		}
		serviceTime += processes[i].BurstDuration

		gantt = append(gantt, TimeSlice{
			PID:   processes[i].ProcessID,
			Start: start,
			Stop:  serviceTime,
		})
	}

	count := float64(len(processes))
	aveWait := totalWait / count
	aveTurnaround := totalTurnaround / count
	aveThroughput := count / lastCompletion

	outputTitle(w, title)
	outputGantt(w, gantt)
	outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)

}

func RRSchedule(w io.Writer, title string, processes []Process) {
	var (
		n                     = len(processes)
		remainingTime         = make([]int64, n)
		arrivalTime           = make([]int64, n)
		waitingTime           = make([]int64, n)
		turnaroundTime        = make([]int64, n)
		schedule              = make([][]string, len(processes))
		currentTime     int64 = 0
		totalTurnaround int64 = 0
		totalWaiting    int64 = 0
	)

	// sort processes by burst duration
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].BurstDuration < processes[j].BurstDuration
	})

	// initialize remainingTime and arrivalTime arrays
	for i := 0; i < n; i++ {
		remainingTime[i] = processes[i].BurstDuration
		arrivalTime[i] = processes[i].ArrivalTime
	}

	for {
		done := true

		for i := 0; i < n; i++ {
			if remainingTime[i] > 0 {
				done = false

				// execute process
				currentTime += remainingTime[i]
				turnaroundTime[i] = currentTime - arrivalTime[i]
				remainingTime[i] = 0
			}
		}

		if done {
			break
		}
	}

	for i := 0; i < n; i++ {
		waitingTime[i] = turnaroundTime[i] - processes[i].BurstDuration
		totalTurnaround += turnaroundTime[i]
		totalWaiting += waitingTime[i]

		schedule[i] = []string{
			fmt.Sprint(processes[i].ProcessID),
			fmt.Sprint(processes[i].Priority),
			fmt.Sprint(processes[i].BurstDuration),
			fmt.Sprint(processes[i].ArrivalTime),
			fmt.Sprint(waitingTime[i]),
			fmt.Sprint(turnaroundTime[i]),
		}
	}

	avgTurnaround := float64(totalTurnaround) / float64(n)
	avgWaiting := float64(totalWaiting) / float64(n)
	throughput := float64(n) / float64(currentTime)

	outputTitle(w, title)
	outputSchedule(w, schedule, avgWaiting, avgTurnaround, throughput)

}

//endregion
