package service

import (
	"fmt"
	"mime/multipart"
	"strings"
	"sync"

	//"github.com/ledongthuc/pdf"
	"github.com/dslipak/pdf"
)

// ExtractTextFromPDF: PDFを並行処理で爆速で読み込む
func ExtractTextFromPDF(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	r, err := pdf.NewReader(f, file.Size)
	if err != nil {
		return "", err
	}

	totalPage := r.NumPage()
	
	// 同時に読むページ数の上限を設定（多すぎるとメモリ不足になる）
	limitPage := 20
	if totalPage > limitPage {
		totalPage = limitPage
	}

	// 各ページのテキストを保存する配列（順番を崩さないため、あらかじめ枠を作る）
	pageTexts := make([]string, totalPage)
	
	var wg sync.WaitGroup // 完了待ち用の同期オブジェクト

	// ページごとに「分身（Goroutine）」を作って走らせる
	for i := 0; i < totalPage; i++ {
		wg.Add(1)
		
		// go func(...) で別スレッド起動！
		go func(index int) {
			defer wg.Done() // 終わったら報告

			// PDFのページ番号は 1 始まり
			p := r.Page(index + 1)
			if p.V.IsNull() {
				return
			}
			
			text, err := p.GetPlainText(nil)
			if err != nil {
				return
			}
			
			// 配列の自分の場所に書き込む（ここは競合しないので安全）
			pageTexts[index] = text
		}(i)
	}

	// 全員の作業が終わるのを待つ
	wg.Wait()

	// バラバラのテキストを結合する
	return strings.Join(pageTexts, "\n"), nil
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