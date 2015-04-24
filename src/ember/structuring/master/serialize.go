package master

import (
	"regexp"
	"encoding/json"
	"ember/structuring/types"
)

func (p *Master) Dones() (urls []string, err error) {
	for url, _ := range p.dones {
		urls = append(urls, url)
	}
	return
}

func (p *Master) UnDones(urls []string) (err error) {
	for _, v := range urls {
		p.dones[v] = true
	}
	return err
}

func (p *Master) donesSerialize(urls []string) (str string, err error) {
	for _, v := range urls {
		str = str + v + "\n"
	}
	return str, err
}

func (p *Master) donesUnSerialize(str string) (url []string, err error) {
	reg := regexp.MustCompile(`http://[^\n]+`)
	url = reg.FindAllString(str, -1)
	return url, err
}

func (p *Master) Doings() (tasks []types.TaskInfo, err error) {
	for _, v := range p.doings {
		tasks = append(tasks, v)
	}
	return tasks, err
}

func (p *Master) UnDoings(tasks []types.TaskInfo) (err error) {
	//TODO Allow redundancy
	for _, v := range tasks {
		p.LoadPush("load", v)
	}
	return
}

func (p *Master) doingsSerialize(tasks []types.TaskInfo) (str string, err error) {
	for _, v := range tasks {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		str = str + string(data) + "\n"
	}
	return str, err
}

func (p *Master) doingsUnSerialize(str string) (tasks []types.TaskInfo, err error) {
	reg := regexp.MustCompile(`[^\n]+`)
	tasksJson := reg.FindAllString(str, -1)

	task := types.TaskInfo {}
	for _, v := range tasksJson {
		err := json.Unmarshal([]byte(v), &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (p *Master) Tasks() (tasks []types.TaskInfo, err error) {
	return p.tasks, err
}

func (p *Master) UnTasks(tasks []types.TaskInfo) (err error) {
	//p.tasks = tasks
	for _, v := range tasks {
		p.LoadPush("load", v)
	}
	return err
}

func (p *Master) tasksSerialize(tasks []types.TaskInfo) (str string, err error) {
	for _, v := range tasks {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		str = str + string(data) + "\n"
	}
	return str, err
}

func (p *Master) tasksUnSerialize(str string) (tasks []types.TaskInfo, err error) {
	reg := regexp.MustCompile(`[^\n]+`)
	tasksJson := reg.FindAllString(str, -1)

	task := types.TaskInfo {}
	for _, v := range tasksJson {
		err := json.Unmarshal([]byte(v), &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}
