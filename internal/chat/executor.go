package chat

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rooveterinaryinc/hello-vim-plugin-2/internal/chat/detector"
	"github.com/rooveterinaryinc/hello-vim-plugin-2/internal/chat/handler"
	"github.com/rooveterinaryinc/hello-vim-plugin-2/internal/chat/prompts"
)

// Message はチャットメッセージを表現する構造体です。
type Message struct {
	Role    string
	Content string
}

// ChatClient はAIとの対話を行うインターフェースです。
type ChatClient interface {
	CreateChatCompletion(messages []Message) (string, error)
}

// Executor はチャット実行を管理する構造体です。
type Executor struct {
	client   ChatClient
	detector *detector.ProposalDetector
	handler  *handler.ProposalHandler
}

// NewExecutor は新しいExecutorインスタンスを作成します。
func NewExecutor(client ChatClient, backupDir string) (*Executor, error) {
	// バックアップディレクトリの設定
	if backupDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("ホームディレクトリの取得に失敗しました: %w", err)
		}
		backupDir = filepath.Join(homeDir, ".roo", "backups")
	}

	// 各コンポーネントの初期化
	fileApplier, err := handler.NewFileApplier(backupDir)
	if err != nil {
		return nil, fmt.Errorf("FileApplierの初期化に失敗しました: %w", err)
	}

	return &Executor{
		client:   client,
		detector: detector.NewProposalDetector(),
		handler: handler.NewProposalHandler(
			handler.NewConsoleApprover(),
			fileApplier,
		),
	}, nil
}

// Execute はチャットを実行し、必要に応じて提案を処理します。
func (e *Executor) Execute(input string) (string, error) {
	// 1. メッセージの準備
	messages := []Message{
		{
			Role:    "system",
			Content: prompts.ProposalSystemPrompt,
		},
		{
			Role:    "user",
			Content: input,
		},
	}

	// 2. AIとの対話
	response, err := e.client.CreateChatCompletion(messages)
	if err != nil {
		return "", fmt.Errorf("チャット実行エラー: %w", err)
	}

	// 3. 提案の検出と処理
	proposal, err := e.detector.Extract(response)
	if err != nil {
		return "", fmt.Errorf("提案検出エラー: %w", err)
	}

	if proposal != nil {
		if err := e.handler.Handle(proposal); err != nil {
			return "", fmt.Errorf("提案処理エラー: %w", err)
		}
	}

	return response, nil
}
