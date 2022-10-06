package util

import (
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"

	"github.com/nimbolus/terraform-backend/pkg/lock"
	"github.com/nimbolus/terraform-backend/pkg/terraform"
)

func LockTest(t *testing.T, l lock.Locker) {
	viper.AutomaticEnv()

	t.Log(l.GetName())

	s1 := terraform.State{
		ID:      terraform.GetStateID("test", "test"),
		Project: "test",
		Name:    "test",
		Lock:    []byte(uuid.New().String()),
	}
	t.Logf("s1: %s", s1.Lock)

	s2 := terraform.State{
		ID:      terraform.GetStateID("test", "test"),
		Project: "test",
		Name:    "test",
		Lock:    []byte(uuid.New().String()),
	}
	t.Logf("s2: %s", s2.Lock)

	// copy of s2
	s3 := terraform.State{
		ID:      terraform.GetStateID("test", "test"),
		Project: "test",
		Name:    "test",
		Lock:    s2.Lock,
	}

	if locked, err := l.Lock(&s1); err != nil || !locked {
		t.Error(err)
	}

	if locked, err := l.Lock(&s1); err != nil || !locked {
		t.Error("should be able to lock twice from the same process")
	}

	if locked, err := l.Lock(&s2); err != nil || locked {
		t.Error("should not be able to lock twice from different processes")
	}

	if string(s2.Lock) != string(s1.Lock) {
		t.Error("failed Lock() should return the lock information of the current lock")
	}

	if unlocked, err := l.Unlock(&s3); err != nil || unlocked {
		t.Error("should not be able to unlock with wrong lock")
	}

	if unlocked, err := l.Unlock(&s1); err != nil || !unlocked {
		t.Error(err)
	}
}