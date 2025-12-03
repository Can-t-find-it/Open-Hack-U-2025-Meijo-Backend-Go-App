package apns

import (
    "fmt"

    "github.com/sideshow/apns2"
    "github.com/sideshow/apns2/token"
)

type APNsClient struct {
    client *apns2.Client
}

func NewAPNsClient() (*APNsClient, error) {

    authKey, err := token.AuthKeyFromFile("AuthKey_XXXXXXXXXX.p8")
    if err != nil {
        return nil, err
    }

    token := &token.Token{
        AuthKey: authKey,
        KeyID:   "あなたのKeyID",
        TeamID:  "あなたのTeamID",
    }

    client := apns2.NewTokenClient(token).Production()
    // 開発中なら .Development()

    return &APNsClient{client: client}, nil
}

func (a *APNsClient) Send(deviceToken string, title string, body string) error {

    notification := &apns2.Notification{
        DeviceToken: deviceToken,
        Payload: []byte(fmt.Sprintf(`{
            "aps" : {
                "alert" : {
                    "title" : "%s",
                    "body" : "%s"
                },
                "sound" : "default"
            }
        }`, title, body)),
    }

    res, err := a.client.Push(notification)
    if err != nil {
        return err
    }

    if res.Sent() {
        fmt.Println("APNs通知成功")
        return nil
    }

    return fmt.Errorf("APNs Error: %v", res.Reason)
}
