package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/rooveterinaryinc/hello-vim-plugin-2/internal/api"
	"github.com/rooveterinaryinc/hello-vim-plugin-2/internal/chat"
	"github.com/rooveterinaryinc/hello-vim-plugin-2/internal/models"
)

// APIClient はAPIとの対話を行うインターフェースです。
type APIClient interface {
	CreateChatCompletion(messages []models.ChatMessage) (string, error)
}

var (
	app = kingpin.New("roo", "A code improvement assistant")

	// chatコマンド
	chatCmd    = app.Command("chat", "Chat with the assistant")
	chatInput  = chatCmd.Arg("input", "Input text for the chat").Required().String()
	chatFile   = chatCmd.Flag("file", "Target file path for patch proposals").String()
	chatBackup = chatCmd.Flag("backup-dir", "Directory for backup files").String()

	// explainコマンド
	explainCmd  = app.Command("explain", "Get code explanation")
	explainCode = explainCmd.Arg("code", "Code to explain").Required().String()
)

func main() {
	// コマンドライン引数のパース
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	// APIクライアントの初期化
	client, err := api.NewClient()
	if err != nil {
		log.Fatalf("Failed to initialize API client: %v", err)
	}

	// コマンドの実行
	var result string
	switch command {
	case chatCmd.FullCommand():
		result, err = executeChat(client, *chatInput, *chatFile, *chatBackup)
	case explainCmd.FullCommand():
		result, err = executeExplain(client, *explainCode)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println(result)
}

// executeCommandは指定されたコマンドを与えられた入力で実行します
func executeCommand(client APIClient, command, input, targetFile, backupDir string) (string, error) {
	switch command {
	case "explain":
		return executeExplain(client, input)
	case "chat":
		return executeChat(client, input, targetFile, backupDir)
	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
}

// executeExplainはexplainコマンドを処理します
func executeExplain(client APIClient, code string) (string, error) {
	messages := []models.ChatMessage{
		{
			Role:    "system",
			Content: "You are a helpful assistant that explains code.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Please explain this code:\n\n%s", code),
		},
	}

	return client.CreateChatCompletion(messages)
}

// executeChatはchatコマンドを処理します
func executeChat(client APIClient, input, targetFile, backupDir string) (string, error) {
	// バックアップディレクトリの設定
	if backupDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		backupDir = filepath.Join(homeDir, ".roo", "backups")
	}

	// アダプターを使用してチャット実行器を初期化
	adapter := chat.NewAPIClientAdapter(client)
	executor, err := chat.NewExecutor(adapter, backupDir)
	if err != nil {
		return "", fmt.Errorf("failed to initialize chat executor: %w", err)
	}

	// チャットの実行
	return executor.Execute(input, targetFile)
}
