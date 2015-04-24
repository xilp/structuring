package master

import (
	"ember/structuring/types"
	"time"
	//"fmt"
)

func (p *Master) save() (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	dones, err:= p.Dones()
	if err != nil {
		return err
	}

	doings, err:= p.Doings()
	if err != nil {
		return err
	}

	tasks, err:= p.Tasks()
	if err != nil {
		return err
	}

	donesStr, err:= p.donesSerialize(dones)
	if err != nil {
		return
	}
	doingsStr, err:= p.doingsSerialize(doings)
	if err != nil {
		return
	}
	tasksStr, err:= p.tasksSerialize(tasks)
	if err != nil {
		return
	}

	p.donesFile.write	(donesStr	, 0)
	p.doingsFile.write	(doingsStr	, 0)
	p.tasksFile.write	(tasksStr	, 0)

	return
}

func (p *Master) load() (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	donesStr, err := p.donesFile.read(0)
	if err != nil {
		return
	}

	doingsStr, err := p.doingsFile.read(0)
	if err != nil {
		return
	}

	tasksStr, err := p.tasksFile.read(0)
	if err != nil {
		return
	}

	dones, err := p.donesUnSerialize(donesStr)
	if err != nil {
		return
	}

	doings, err:= p.doingsUnSerialize(doingsStr)
	if err != nil {
		return
	}

	tasks, err:= p.tasksUnSerialize(tasksStr)
	if err != nil {
		return
	}

	err = p.UnDones(dones)
	if err != nil {
		return err
	}

	err = p.UnDoings(doings)
	if err != nil {
		return err
	}

	err = p.UnTasks(tasks)
	if err != nil {
		return err
	}

	return
}

func (p *Master) LoadPush(slave string, info types.TaskInfo) (err error) {
	//fmt.Printf("appending %v\n", info)
	p.slaves[slave] = time.Now().UnixNano()

	if p.dones[info.Url] {
		return
	}

	if p.todos[info.Url] {
		return
	}

	if _, ok := p.doings[info.Url] ; ok {
		return
	}

	p.tasks = append(p.tasks, info)
	p.todos[info.Url] = true

	return
}
