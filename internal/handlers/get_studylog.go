package handlers

import (
	"net/http"
	"hacku_2025_meijo/internal/service"
	"hacku_2025_meijo/internal/dtos"
	"github.com/gin-gonic/gin"
)

// ResponseLog と StudyLogsWrapper の定義は、このファイルのどこかにあるか、
// または dtos パッケージにあるものと仮定します。

func GetStudyLogsHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JWTからユーザーIDの取得に失敗しました"})
		return
	}
	
    // サービス層から最新ログを単一オブジェクトとして取得
    // サービス関数が *dtos.StudyLogResponse または *ResponseLog を返すと仮定
	latestLog, err := service.GetLatestStudyLog(userID)
	
	if err != nil {
		/*
        if errors.Is(err, gorm.ErrRecordNotFound) { 
            c.JSON(http.StatusOK, StudyLogsWrapper{Logs: []ResponseLog{}})
            return
        }
        */
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    // 取得した単一のログオブジェクトを要素とするリストを作成
    var logsList []dtos.StudyLogResponse 
    if latestLog != nil {
        // ポインタをデリファレンスしてリストに追加
        logsList = append(logsList, *latestLog) 
    }
    
	// レスポンスラッパーにリストをセット
    response := dtos.StudyLogsWrapper{
        Logs: logsList,
    }

    // JSONレスポンスとして返す: {"logs": [{...}]}
	c.JSON(http.StatusOK, response) 
}