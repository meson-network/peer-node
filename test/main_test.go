package test

import (
	"log"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	log.Println("run test")

	str := "asdfasdfsfsf_mark"

	r := strings.Split(str, "_mark_")
	log.Println(r)
}
