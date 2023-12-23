package service

import "fmt"

type Service struct {
	Name     string
	Id       string
	Port     uint
	TargetId string
}

var (
	services map[string]Service
)

func init() {
	services = make(map[string]Service)
}

func Add(s Service) error {
	if _, ok := services[s.Id]; ok {
		return fmt.Errorf("this is service id is existed")
	}
	services[s.Id] = s

	return nil
}

func Remove(id string) {
	delete(services, id)
}

func RemoveByTarget(target_id string) {
	ids := []string{}
	for _, s := range services {
		if s.TargetId == target_id {
			ids = append(ids, s.Id)
		}
	}
	for _, t := range ids {
		delete(services, t)
	}
}

func Update(s Service) error {
	ss, ok := services[s.Id]
	if !ok {
		return fmt.Errorf("service id is not exist for update")
	}

	if ss.TargetId == s.TargetId {
		services[s.Id] = s
	} else {
		return fmt.Errorf("targets not equal for update")
	}

	return nil
}
