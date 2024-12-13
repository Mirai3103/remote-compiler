package executor

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"go.uber.org/zap"
)

type Executor interface {
	Compile(*model.Submission) error
	Execute(*model.Submission, chan<- *model.SubmissionResult) error
}

type executor struct {
	logger *zap.Logger
	cfx    config.ExecutorConfig
}

func (e *executor) run(command string) (*string, error) {
	// rs := exec.Command(*command)
	// err := rs.Run()
	// if err != nil {
	// 	return err
	// }
	// return nil
	randomSleep := rand.IntN(20)
	time.Sleep(time.Duration(randomSleep) * time.Second)
	fmt.Print("Sleeping for ", randomSleep, " seconds\n")
	return nil, nil
}

// Execute implements Executor.
func (e *executor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	testcase := submission.TestCases
	var wg sync.WaitGroup
	for _, testcase := range testcase {
		wg.Add(1)
		go func(testcase *model.TestCase) {
			defer wg.Done()
			result := &model.SubmissionResult{
				SubmissionID: submission.ID,
				TestCaseID:   testcase.ID,
			}
			_, err := e.run("echo ")
			if err != nil {
				var runTimeError = "Run time error"
				var e = err.Error()
				result.Status = &runTimeError
				result.Stderr = &e
			}
			ch <- result
		}(&testcase)
	}
	wg.Wait()
	close(ch)
	return nil
}

// Compile implements Executor.
func (e *executor) Compile(sb *model.Submission) error {
	if sb.Language.CompileCommand == nil {
		return nil
	}
	language := sb.Language
	code := sb.Code
	sourceFilename := e.cfx.CompileDir + "/" + language.GetSourceFileName()
	binaryFilename := e.cfx.CompileDir + "/" + language.GetBinaryFileName()
	err := os.WriteFile(sourceFilename, []byte(*code), 0644)
	if err != nil {
		return err
	}

	command := strings.ReplaceAll(*language.CompileCommand, "$SourceFileName", sourceFilename)
	command = strings.ReplaceAll(command, "$BinaryFileName", binaryFilename)

	execCmd := exec.Command("sh", "-c", command)
	log.Printf("Running command: %s", command)
	execCmd.Dir = e.cfx.CompileDir
	err = execCmd.Run()
	execCmd.Wait()
	if err != nil {
		return err
	}
	defer os.Remove(sourceFilename)
	return nil

}

func NewExecutor(logger *zap.Logger, cfx config.ExecutorConfig) Executor {
	if _, err := os.Stat(cfx.CompileDir); os.IsNotExist(err) {
		os.MkdirAll(cfx.CompileDir, 0755)
	}
	return &executor{
		logger: logger,
		cfx:    cfx,
	}
}
