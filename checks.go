package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/unprofession-al/objectstore"
)

type Check struct {
	State            int
	Time             time.Time
	ValidityDuration time.Duration
	Body             []byte
}

func (c Check) TimedOut() bool {
	since := time.Since(c.Time)
	if since > c.ValidityDuration {
		return true
	}
	return false
}

func NewCheck(s int, d time.Duration, b []byte) Check {
	return Check{
		State:            s,
		Time:             time.Now(),
		ValidityDuration: d,
		Body:             b,
	}
}

type CheckStates struct {
	path string
}

func NewCheckStates(path string) *CheckStates {
	path = strings.TrimRight(path, "/")
	return &CheckStates{
		path: path,
	}
}

func (cs *CheckStates) Set(key string, state int, valDur string, body []byte) error {
	path := fmt.Sprintf("%s/%s", cs.path, key)
	o, err := objectstore.New(path)
	if err != nil {
		return err
	}

	dur, err := time.ParseDuration(valDur)
	if err != nil {
		return err
	}

	c := NewCheck(state, dur, body)

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = o.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CheckStates) Get(key string) (Check, error) {
	var result Check
	path := fmt.Sprintf("%s/%s", cs.path, key)
	o, err := objectstore.New(path)
	if err != nil {
		return result, err
	}

	body, err := o.Read()
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
