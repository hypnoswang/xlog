package xlog

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"testing"
)

func TestRotateNone(t *testing.T) {
	os.Remove("/tmp/file_sinker/testlog")
	conf := NewFileSinkerConf("/tmp/file_sinker/testlog")
	if conf == nil {
		t.Errorf("Create FileSinkerConf failed")
	}

	sinker := NewFileSinker(conf)
	if sinker == nil {
		t.Errorf("Create FileSinker failed")
	}

	sinker.rotate()

	data := []byte("This is for none test data \n")
	for i := 0; i < 1000; i++ {
		sinker.Write(data)
	}

	sinker.Close()

	if _, err := os.Stat("/tmp/file_sinker/testlog"); err != nil {
		t.Errorf("Create none log file failed")
	}

	cmd := exec.Command("wc", "-l", "/tmp/file_sinker/testlog")
	expected := "1000 /tmp/file_sinker/testlog\n"
	if output, err := cmd.Output(); err != nil || string(output) != expected {
		t.Errorf("log file lines is not 1000, %s", string(output))
	}
}

func TestRotateTimeHourly(t *testing.T) {
	tstr := time.Now().Format("2006-01-02-15")
	fname := fmt.Sprintf("/tmp/file_sinker/testlog.%s", tstr)
	t.Logf("The hourly log filename is %s", fname)

	os.Remove(fname)

	conf := NewFileSinkerConf("/tmp/file_sinker/testlog")
	if conf == nil {
		t.Errorf("Create FileSinkerConf failed")
	}

	conf.RotateTime = TimeRotateHourly
	sinker := NewFileSinker(conf)
	if sinker == nil {
		t.Errorf("Create FileSinker failed")
	}

	sinker.rotate()

	data := []byte("This is for hourly test data \n")
	for i := 0; i < 1000; i++ {
		sinker.Write(data)
	}

	sinker.Close()

	if _, err := os.Stat(fname); err != nil {
		t.Errorf("Create hourly log file failed")
	}

	cmd := exec.Command("wc", "-l", fname)
	expected := fmt.Sprintf("1000 %s\n", fname)
	if output, err := cmd.Output(); err != nil || string(output) != expected {
		t.Errorf("log file lines is not 1000, %s", string(output))
	}
}

func TestRotateTimeDaily(t *testing.T) {
	tstr := time.Now().Format("2006-01-02")
	fname := fmt.Sprintf("/tmp/file_sinker/testlog.%s", tstr)
	t.Logf("The daily log filename is %s", fname)

	os.Remove(fname)

	conf := NewFileSinkerConf("/tmp/file_sinker/testlog")
	if conf == nil {
		t.Errorf("Create FileSinkerConf failed")
	}

	conf.RotateTime = TimeRotateDaily
	sinker := NewFileSinker(conf)
	if sinker == nil {
		t.Errorf("Create FileSinker failed")
	}

	sinker.rotate()

	data := []byte("This is for daily test data \n")
	for i := 0; i < 1000; i++ {
		sinker.Write(data)
	}

	sinker.Close()

	if _, err := os.Stat(fname); err != nil {
		t.Errorf("Create daily log file failed")
	}

	cmd := exec.Command("wc", "-l", fname)
	expected := fmt.Sprintf("1000 %s\n", fname)
	if output, err := cmd.Output(); err != nil || string(output) != expected {
		t.Errorf("log file lines is not 1000, %s", string(output))
	}
}

func TestRotateSize1(t *testing.T) {
	fname := fmt.Sprintf("/tmp/file_sinker/testlog.%d", 0)
	t.Logf("The daily log filename is %s", fname)

	os.Remove(fname)

	conf := NewFileSinkerConf("/tmp/file_sinker/testlog")
	if conf == nil {
		t.Errorf("Create FileSinkerConf failed")
	}

	conf.RotateSize = 1024 * 1024 * 100
	sinker := NewFileSinker(conf)
	if sinker == nil {
		t.Errorf("Create FileSinker failed")
	}

	sinker.rotate()

	data := []byte("This is for size test 1 data \n")
	for i := 0; i < 1000; i++ {
		sinker.Write(data)
	}

	sinker.Close()

	if _, err := os.Stat(fname); err != nil {
		t.Errorf("Create daily log file failed")
	}

	cmd := exec.Command("wc", "-l", fname)
	expected := fmt.Sprintf("1000 %s\n", fname)
	if output, err := cmd.Output(); err != nil || string(output) != expected {
		t.Errorf("log file lines is not 1000, %s", string(output))
	}
}

func TestRotateSize2(t *testing.T) {
	for i := 0; i < 11; i++ {
		f := fmt.Sprintf("/tmp/file_sinker/testlog.%d", i)
		if fp, err := os.Create(f); err != nil {
			t.Errorf("Create file %s failed: %s", f, err.Error())
		} else {
			fp.Close()
		}
	}

	fname := fmt.Sprintf("/tmp/file_sinker/testlog.%d", 11)
	t.Logf("The daily log filename is %s", fname)

	os.Remove(fname)

	conf := NewFileSinkerConf("/tmp/file_sinker/testlog")
	if conf == nil {
		t.Errorf("Create FileSinkerConf failed")
	}

	conf.RotateSize = 1024 * 1024 * 100
	sinker := NewFileSinker(conf)
	if sinker == nil {
		t.Errorf("Create FileSinker failed")
	}

	sinker.rotate()

	data := []byte("This is for size test 2 data \n")
	for i := 0; i < 1000; i++ {
		sinker.Write(data)
	}

	sinker.Close()

	if _, err := os.Stat(fname); err != nil {
		t.Errorf("Create daily log file failed")
	}

	cmd := exec.Command("wc", "-l", fname)
	expected := fmt.Sprintf("1000 %s\n", fname)
	if output, err := cmd.Output(); err != nil || string(output) != expected {
		t.Errorf("log file lines is not 1000, %s", string(output))
	}
}

func TestRotateBoth(t *testing.T) {
	tstr := time.Now().Format("2006-01-02-15")
	fname := fmt.Sprintf("/tmp/file_sinker/testlog.%s.0", tstr)
	t.Logf("The rotate-both log filename is %s", fname)

	os.Remove(fname)

	conf := NewFileSinkerConf("/tmp/file_sinker/testlog")
	if conf == nil {
		t.Errorf("Create FileSinkerConf failed")
	}

	conf.RotateTime = TimeRotateHourly
	conf.RotateSize = 1024 * 1024 * 100
	sinker := NewFileSinker(conf)
	if sinker == nil {
		t.Errorf("Create FileSinker failed")
	}

	sinker.rotate()

	data := []byte("This is for both time and size test data \n")
	for i := 0; i < 1000; i++ {
		sinker.Write(data)
	}

	sinker.Close()

	if _, err := os.Stat(fname); err != nil {
		t.Errorf("Create daily log file failed")
	}

	cmd := exec.Command("wc", "-l", fname)
	expected := fmt.Sprintf("1000 %s\n", fname)
	if output, err := cmd.Output(); err != nil || string(output) != expected {
		t.Errorf("log file lines is not 1000, %s", string(output))
	}
}
