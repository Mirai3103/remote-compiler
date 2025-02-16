package executor

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"

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
	cfg    config.ExecutorConfig
}

func (e *executor) Execute(submission *model.Submission, ch chan<- *model.SubmissionResult) error {
	command := e.buildCommand(submission)
	defer e.cleanupFiles(submission)

	var wg sync.WaitGroup
	for _, testcase := range submission.TestCases {
		wg.Add(1)
		go func(tc model.TestCase) {
			defer wg.Done()
			executor := NewTestCaseExecutor(e.logger, e.cfg.IsolateDir, command, submission)
			ch <- executor.Execute(&tc)
		}(testcase)
	}

	wg.Wait()
	close(ch)
	return nil
}

func (e *executor) Compile(submission *model.Submission) error {
	sourceFilename := e.writeSourceFile(submission)
	if sourceFilename == "" {
		return errors.New("failed to write source file")
	}
	if !e.needsCompilation(submission) {
		return nil
	}

	command := e.buildCompileCommand(submission, sourceFilename)
	return e.runCompileCommand(command)
}

func (e *executor) buildCommand(submission *model.Submission) string {
	command := *submission.Language.RunCommand
	command = strings.ReplaceAll(command, "$BinaryFileName", e.cfg.IsolateDir+"/"+submission.Language.GetBinaryFileName())
	return strings.ReplaceAll(command, "$SourceFileName", e.cfg.IsolateDir+"/"+submission.Language.GetSourceFileName())
}

func (e *executor) cleanupFiles(submission *model.Submission) {
	os.Remove(e.cfg.IsolateDir + "/" + submission.Language.GetSourceFileName())
	os.Remove(e.cfg.IsolateDir + "/" + submission.Language.GetBinaryFileName())
}

func (e *executor) needsCompilation(submission *model.Submission) bool {
	return submission.Language.CompileCommand != nil && *submission.Language.CompileCommand != ""
}

func (e *executor) writeSourceFile(submission *model.Submission) string {
	sourceFilename := e.cfg.IsolateDir + "/" + submission.Language.GetSourceFileName()
	e.logger.Info("Writing source file", zap.String("filename", sourceFilename))
	err := os.WriteFile(sourceFilename, []byte(*submission.Code), 0644)
	if err != nil {
		e.logger.Error("Error writing source file", zap.Error(err))
		return ""
	}
	return sourceFilename
}

func (e *executor) buildCompileCommand(submission *model.Submission, sourceFilename string) string {
	binaryFilename := e.cfg.IsolateDir + "/" + submission.Language.GetBinaryFileName()
	command := strings.ReplaceAll(*submission.Language.CompileCommand, "$SourceFileName", sourceFilename)
	return strings.ReplaceAll(command, "$BinaryFileName", binaryFilename)
}

func (e *executor) runCompileCommand(command string) error {
	execCmd := exec.Command("sh", "-c", command)
	execCmd.Dir = e.cfg.IsolateDir
	var stderr strings.Builder
	execCmd.Stderr = &stderr

	if err := execCmd.Run(); err != nil {
		e.logger.Error("Compilation error",
			zap.Error(err),
			zap.String("stderr", stderr.String()))
		return errors.New(stderr.String())
	}
	return nil
}

func NewExecutor(logger *zap.Logger, cfg config.ExecutorConfig) Executor {
	return &executor{
		logger: logger,
		cfg:    cfg,
	}
}
