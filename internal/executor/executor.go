package executor

import (
	"fmt"
	"math/rand/v2"
	"os"
	"sync"
	"time"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"go.uber.org/zap"
)

type Executor interface {
	Compile(*model.Submission) error
	Execute(*model.Submission, chan<- *model.SubmissionResult) error
}

type executor struct {
	logger *zap.Logger
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
	sourceFilename := language.GetSourceFileName()
	err := os.WriteFile(sourceFilename, []byte(*code), 0644)
	if err != nil {
		return err
	}
	defer os.Remove(sourceFilename)
	return nil

}

func NewExecutor(logger *zap.Logger) Executor {
	return &executor{
		logger: logger,
	}
}
