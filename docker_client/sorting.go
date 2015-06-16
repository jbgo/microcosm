package docker_client

type ByAge Containers

func (s ByAge) Len() int {
	return len(s)
}

func (s ByAge) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByAge) Less(i, j int) bool {
	a, b := s[i].Original.State, s[j].Original.State
	if a.Running && b.Running {
		return a.StartedAt.Unix() < b.StartedAt.Unix()
	} else if a.Running && !b.Running {
		return true
	} else if !a.Running && b.Running {
		return false
	} else {
		return a.FinishedAt.Unix() > b.FinishedAt.Unix()
	}
}

type ByName Containers

func (s ByName) Len() int {
	return len(s)
}

func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
