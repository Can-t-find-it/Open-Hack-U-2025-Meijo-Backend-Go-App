package service

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractTextFromPDF: アップロードされたPDFファイルからテキストを抽出する
func ExtractTextFromPDF(file *multipart.FileHeader) (string, error) {
	// ファイルを開く
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// PDFリーダーを作成 (ファイルサイズが必要)
	r, err := pdf.NewReader(f, file.Size)
	if err != nil {
		return "", err
	}

	var content string
	
	// ページ数分ループしてテキストを取得
	// (全部読むと長すぎる場合は、最初の5ページだけにするなどの制限も可)
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		
		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		content += text + "\n"
	}

	return content, nil
}

// ExtractKeywordsFromText: テキストデータをAIに渡して、重要単語を抽出させる
func ExtractKeywordsFromText(text string) ([]string, error) {
	// トークン制限対策: テキストが長すぎる場合は先頭3000文字くらいに切る
	if len(text) > 3000 {
		text = text[:3000]
	}

	prompt := fmt.Sprintf(`
	以下の文章は、ある学習資料のテキストデータです。
	この文章の中から、学習者が覚えるべき「重要な専門用語・キーワード」を最大10個抜き出してください。

	対象テキスト:
	%s

	【出力条件】
	- 出力は「単語1, 単語2, 単語3」のように、カンマ区切りのリストのみにしてください。
	- 余計な挨拶や説明は一切含めないでください。
	- 日本語または英語の単語のみ出力してください。
	`, text)

	apiResp, err := callOpenAIAPI(prompt)
	if err != nil {
		return nil, err
	}

	// カンマ区切りを配列に変換
	rawContent := strings.TrimSpace(apiResp.Choices[0].Message.Content)
	rawContent = strings.ReplaceAll(rawContent, "、", ",") // 全角カンマ対策
	
	rawList := strings.Split(rawContent, ",")
	var resultList []string
	for _, w := range rawList {
		cleaned := strings.TrimSpace(w)
		if cleaned != "" {
			resultList = append(resultList, cleaned)
		}
	}

	return resultList, nil
}