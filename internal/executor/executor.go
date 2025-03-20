package executor

import (
	"os"
	"strings"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"github.com/Mirai3103/remote-compiler/pkg/config"
	"go.uber.org/zap"
)

type Executor interface {
	Compile(*model.Submission) error
	Execute(*model.Submission, chan<- *model.SubmissionResult) error
}

type baseExecutor struct {
	logger *zap.Logger
	cfg    config.ExecutorConfig
}

func (e *baseExecutor) writeSourceFile(submission *model.Submission) string {
	sourceFilename := e.cfg.IsolateDir + "/" + submission.Language.GetSourceFileName()
	e.logger.Info("Writing source file", zap.String("filename", sourceFilename))
	err := os.WriteFile(sourceFilename, []byte(*submission.Code), 0644)
	if err != nil {
		e.logger.Error("Error writing source file", zap.Error(err))
		return ""
	}
	return sourceFilename
}

func (e *baseExecutor) buildCompileCommand(submission *model.Submission, sourceFilename string) string {
	binaryFilename := e.cfg.IsolateDir + "/" + submission.Language.GetBinaryFileName()
	command := strings.ReplaceAll(*submission.Language.CompileCommand, "$SourceFileName", sourceFilename)
	return strings.ReplaceAll(command, "$BinaryFileName", binaryFilename)
}

func (e *baseExecutor) cleanupFiles(submission *model.Submission) {
	os.Remove(e.cfg.IsolateDir + "/" + submission.Language.GetSourceFileName())
	os.Remove(e.cfg.IsolateDir + "/" + submission.Language.GetBinaryFileName())
}

func (e *baseExecutor) needsCompilation(submission *model.Submission) bool {
	return submission.Language.CompileCommand != nil && *submission.Language.CompileCommand != ""
}

func (e *baseExecutor) buildCommand(submission *model.Submission) string {
	command := *submission.Language.RunCommand
	command = strings.ReplaceAll(command, "$BinaryFileName", e.cfg.IsolateDir+"/"+submission.Language.GetBinaryFileName())
	return strings.ReplaceAll(command, "$SourceFileName", e.cfg.IsolateDir+"/"+submission.Language.GetSourceFileName())
}
