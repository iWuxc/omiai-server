package c_interact

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"omiai-server/internal/biz"
	biz_omiai "omiai-server/internal/biz/omiai"
	"omiai-server/internal/data"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockClientRepo 模拟 Client 数据层
type MockClientRepo struct {
	interactions []*biz_omiai.ClientInteraction
	clients      map[uint64]*biz_omiai.Client
}

func (m *MockClientRepo) Select(ctx context.Context, clause *biz.WhereClause, fields []string, offset, limit int) ([]*biz_omiai.Client, error) {
	return nil, nil
}
func (m *MockClientRepo) Create(ctx context.Context, client *biz_omiai.Client) error { return nil }
func (m *MockClientRepo) Update(ctx context.Context, client *biz_omiai.Client) error { return nil }
func (m *MockClientRepo) Delete(ctx context.Context, id uint64) error                { return nil }
func (m *MockClientRepo) Get(ctx context.Context, id uint64) (*biz_omiai.Client, error) {
	if c, ok := m.clients[id]; ok {
		return c, nil
	}
	return nil, nil
}
func (m *MockClientRepo) Stats(ctx context.Context) (map[string]int64, error) { return nil, nil }
func (m *MockClientRepo) GetDashboardStats(ctx context.Context) (map[string]int64, error) {
	return nil, nil
}
func (m *MockClientRepo) GetByWxOpenID(ctx context.Context, openID string) (*biz_omiai.Client, error) {
	return nil, nil
}

func (m *MockClientRepo) SaveInteraction(ctx context.Context, interaction *biz_omiai.ClientInteraction) error {
	m.interactions = append(m.interactions, interaction)
	return nil
}

func (m *MockClientRepo) GetInteraction(ctx context.Context, fromID, toID uint64) (*biz_omiai.ClientInteraction, error) {
	for _, v := range m.interactions {
		if v.FromClientID == fromID && v.ToClientID == toID {
			return v, nil
		}
	}
	return nil, nil
}

func (m *MockClientRepo) GetInteractionLeads(ctx context.Context, managerID uint64, offset, limit int) ([]*biz_omiai.ClientInteraction, error) {
	return nil, nil
}

func (m *MockClientRepo) GetClientInteractions(ctx context.Context, clientID uint64, actionType int8, offset, limit int) ([]*biz_omiai.ClientInteraction, error) {
	var res []*biz_omiai.ClientInteraction
	for _, v := range m.interactions {
		if actionType == 2 && v.ToClientID == clientID && v.ActionType == 2 {
			res = append(res, v)
		} else if actionType == 3 && (v.FromClientID == clientID || v.ToClientID == clientID) && v.ActionType == 3 {
			res = append(res, v)
		}
	}
	return res, nil
}

// MockMatchRepo 模拟 Match 数据层
type MockMatchRepo struct {
	createdRecords []*biz_omiai.MatchRecord
}

func (m *MockMatchRepo) Select(ctx context.Context, clause *biz.WhereClause, offset, limit int) ([]*biz_omiai.MatchRecord, error) {
	return nil, nil
}
func (m *MockMatchRepo) Create(ctx context.Context, record *biz_omiai.MatchRecord) error {
	m.createdRecords = append(m.createdRecords, record)
	return nil
}
func (m *MockMatchRepo) Update(ctx context.Context, record *biz_omiai.MatchRecord) error { return nil }
func (m *MockMatchRepo) Get(ctx context.Context, id uint64) (*biz_omiai.MatchRecord, error) {
	return nil, nil
}
func (m *MockMatchRepo) ConfirmMatch(ctx context.Context, clientID, candidateID uint64, adminID, remark string) (*biz_omiai.MatchRecord, error) {
	return nil, nil
}
func (m *MockMatchRepo) Compare(ctx context.Context, maleID, femaleID uint64) (*biz_omiai.Comparison, error) {
	return nil, nil
}
func (m *MockMatchRepo) UpdateStatus(ctx context.Context, recordID uint64, oldStatus, newStatus int8, operator, reason string) error {
	return nil
}
func (m *MockMatchRepo) DissolveMatch(ctx context.Context, clientID uint64, operator, reason string) error {
	return nil
}
func (m *MockMatchRepo) GetStatusHistory(ctx context.Context, recordID uint64) ([]*biz_omiai.MatchStatusHistory, error) {
	return nil, nil
}
func (m *MockMatchRepo) CreateFollowUp(ctx context.Context, record *biz_omiai.FollowUpRecord) error {
	return nil
}
func (m *MockMatchRepo) SelectFollowUps(ctx context.Context, matchRecordID uint64) ([]*biz_omiai.FollowUpRecord, error) {
	return nil, nil
}
func (m *MockMatchRepo) SelectAllFollowUps(ctx context.Context, offset, limit int) ([]*biz_omiai.FollowUpRecord, error) {
	return nil, nil
}
func (m *MockMatchRepo) GetReminders(ctx context.Context) ([]*biz_omiai.MatchRecord, error) {
	return nil, nil
}
func (m *MockMatchRepo) Stats(ctx context.Context) (map[string]interface{}, error) { return nil, nil }
func (m *MockMatchRepo) GetCandidates(ctx context.Context, clientID uint64) ([]*biz_omiai.Candidate, error) {
	return nil, nil
}
func (m *MockMatchRepo) Delete(ctx context.Context, id uint64) error { return nil }

type MockNotifier struct{}

func (m *MockNotifier) NotifyManager(ctx context.Context, managerID uint64, title, content string) error {
	return nil
}
func (m *MockNotifier) NotifyClient(ctx context.Context, clientID uint64, openID, title, content string) error {
	return nil
}

func setupTestContext(clientID uint64, reqBody interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// 模拟中间件设置的 client_id
	ctx.Set("client_id", clientID)

	if reqBody != nil {
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		ctx.Request = req
	} else {
		ctx.Request, _ = http.NewRequest("GET", "/", nil)
	}

	return ctx, w
}

func TestLike_SingleMatch(t *testing.T) {
	mockClient := &MockClientRepo{
		interactions: make([]*biz_omiai.ClientInteraction, 0),
		clients: map[uint64]*biz_omiai.Client{
			1: {ID: 1, Gender: 1, Name: "Male User"},
			2: {ID: 2, Gender: 2, Name: "Female User"},
		},
	}
	mockMatch := &MockMatchRepo{}
	controller := NewController(&data.DB{}, mockClient, mockMatch, &MockNotifier{})

	// User 1 likes User 2
	ctx, w := setupTestContext(1, LikeRequest{TargetClientID: 2})
	controller.Like(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, float64(0), resp["code"])
	dataMap := resp["data"].(map[string]interface{})
	assert.False(t, dataMap["is_matched"].(bool)) // 仅仅是单向喜欢，不会 matched

	// 验证 interaction 数据
	assert.Equal(t, 1, len(mockClient.interactions))
	assert.Equal(t, uint64(1), mockClient.interactions[0].FromClientID)
	assert.Equal(t, uint64(2), mockClient.interactions[0].ToClientID)
	assert.Equal(t, int8(2), mockClient.interactions[0].ActionType)

	// 验证没有生成 MatchRecord
	assert.Equal(t, 0, len(mockMatch.createdRecords))
}

func TestLike_MutualMatch(t *testing.T) {
	mockClient := &MockClientRepo{
		interactions: []*biz_omiai.ClientInteraction{
			// 预设 User 2 已经喜欢了 User 1
			{FromClientID: 2, ToClientID: 1, ActionType: 2, CreatedAt: time.Now()},
		},
		clients: map[uint64]*biz_omiai.Client{
			1: {ID: 1, Gender: 1, Name: "Male User"},
			2: {ID: 2, Gender: 2, Name: "Female User"},
		},
	}
	mockMatch := &MockMatchRepo{}
	controller := NewController(&data.DB{}, mockClient, mockMatch, &MockNotifier{})

	// User 1 喜欢 User 2 (触发双向奔赴)
	ctx, w := setupTestContext(1, LikeRequest{TargetClientID: 2})
	controller.Like(ctx)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	dataMap := resp["data"].(map[string]interface{})
	assert.True(t, dataMap["is_matched"].(bool)) // 验证双向奔赴成功

	// 验证 interaction 数据 (因为互相心动，会升级旧的记录，同时产生两条保存调用，长度可能会变，所以我们只验证是否有MatchRecord产生即可)

	// 验证生成了相识记录
	assert.Equal(t, 1, len(mockMatch.createdRecords))
	record := mockMatch.createdRecords[0]
	assert.Equal(t, int8(biz_omiai.MatchStatusAcquaintance), record.Status)
	assert.Equal(t, uint64(1), record.MaleClientID)
	assert.Equal(t, uint64(2), record.FemaleClientID)
}

func TestGetReceivedLikes(t *testing.T) {
	mockClient := &MockClientRepo{
		interactions: []*biz_omiai.ClientInteraction{
			{ID: 100, FromClientID: 2, ToClientID: 1, ActionType: 2, CreatedAt: time.Now()},
		},
		clients: map[uint64]*biz_omiai.Client{
			1: {ID: 1, Gender: 1, Name: "User1"},
			2: {ID: 2, Gender: 2, Name: "User2", Avatar: "avatar.jpg"},
		},
	}
	controller := NewController(&data.DB{}, mockClient, &MockMatchRepo{}, &MockNotifier{})

	ctx, w := setupTestContext(1, nil)
	controller.GetReceivedLikes(ctx)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	list := resp["data"].([]interface{})
	assert.Equal(t, 1, len(list))

	item := list[0].(map[string]interface{})
	assert.Equal(t, "U***", item["name"]) // 验证姓名被脱敏
}
