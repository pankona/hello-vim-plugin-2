# Makasero プロジェクト概要

## プロジェクトの目的

Makaseroは、OpenAI GPT-4を活用したコード分析・改善支援のためのCLIツールです。主な目的は：

1. コードの説明生成と理解支援
2. インタラクティブなコードレビューと改善
3. コード改善提案の自動生成と適用
4. 継続的な対話を通じたコード品質の向上

## 主要機能

### 1. コマンドライン機能
- `explain`: コードやファイルの説明生成
  - ファイルパスまたはコードスニペットを直接指定可能
  - 詳細な説明と改善提案の生成

- `chat`: AIとのインタラクティブな対話
  - 常時インタラクティブモード
  - ファイル編集モードでのコード改善
  - セッション管理による会話の継続性確保

### 2. コード提案システム
- AIによる改善提案の生成
  - 標準フォーマットでの提案（PROPOSAL/CODE）
  - 差分の自動検出と表示
  - 変更のプレビューと承認フロー
- コマンド分析と実行
  - 特殊コマンド（/test等）の自動検出
  - コンテキストに応じたコマンド生成
  - 承認フローに基づく安全な実行

### 3. セッション管理
- 会話履歴の永続化
  - `~/.makasero/sessions/`への自動保存
  - セッションの一覧表示と内容確認
  - 過去のセッションの再開機能

### 4. 安全性機能
- 自動バックアップ
  - 変更前のファイルを自動保存
  - タイムスタンプ付きバックアップ
  - カスタマイズ可能なバックアップディレクトリ

## アーキテクチャ

### コア構成
1. コマンドライン処理（cmd/makasero）
   - ユーザー入力の解析
   - コマンドルーティング
   - 対話セッション管理
   - 特殊コマンドの処理

2. API通信（internal/api）
   - OpenAI APIとの通信
   - エラーハンドリング
   - レスポンス処理
   - コマンド提案の生成

3. モデル（internal/models）
   - チャットメッセージの構造
   - セッション管理
   - レスポンス形式の定義
   - コマンドアナライザーとランナー

4. 提案システム（internal/proposal）
   - コード差分の生成
   - 変更の検証
   - バックアップ管理
   - コマンド実行の制御

## 設計方針

1. シンプルさと使いやすさ
   - 直感的なコマンド体系
   - 対話型インターフェース
   - 明確なフィードバック

2. 安全性
   - 自動バックアップ
   - 変更の承認フロー
   - エラーハンドリング

3. 拡張性
   - モジュラー設計
   - 機能の段階的追加
   - プラグイン可能なアーキテクチャ 