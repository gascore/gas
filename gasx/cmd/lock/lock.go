package lock

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/gascore/gas/gasx/cmd/config"
	"github.com/liamylian/jsontime"
)

var json = jsontime.ConfigWithCustomTimeFormat

type Lock struct {
	AtomicStyles map[string]string `json:"atomic_styles"`
	Styles       map[string]string `json:"styles"`
	Date         time.Time         `json:"last" time_format:"2006-01-02 15:04:05"`
}

const gasLock = ".gaslock"

func ParseGasLock(config *config.Config) (*Lock, bool, error) {
	lockFile, err := ioutil.ReadFile(gasLock)
	if err != nil {
		if os.IsNotExist(err) {
			return &Lock{}, true, nil
		}

		return nil, false, err
	}

	var lock Lock
	err = json.Unmarshal(lockFile, &lock)
	if err != nil {
		return nil, false, err
	}

	buildExternal := !config.IgnoreExternal || lock.Date.After(time.Now().Add(-24*time.Hour))

	return &lock, buildExternal, nil
}

func (l *Lock) Save() error {
	l.Date = time.Now()

	lockJSON, err := json.Marshal(&l)
	if err != nil {
		return err
	}

	if exists(gasLock) {
		err := os.Remove(gasLock)
		if err != nil {
			return err
		}
	}

	lockFile, err := os.Create(gasLock)
	if err != nil {
		return err
	}

	_, err = lockFile.Write(lockJSON)
	if err != nil {
		return err
	}

	return nil
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
