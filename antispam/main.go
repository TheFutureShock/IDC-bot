package antispam

import "time"

type AntiSpam struct {
	Users map[string]struct {
		N int
		M []string
	}

	Timer *time.Timer
}

func (A *AntiSpam) Increase(user string, messageID string) bool {
	A.Users[user] = struct {
		N int
		M []string
	}{N: A.Users[user].N + 1, M: append(A.Users[user].M, messageID)}

	if A.Users[user].N >= 5 {
		return true
	}
	return false
}

func (A *AntiSpam) Init() {
	A.Timer = time.AfterFunc(time.Second*3, A.reset)
}

func (A *AntiSpam) reset() {
	A.Users = make(map[string]struct {
		N int
		M []string
	})
	A.Timer.Reset(time.Second * 3)
}
