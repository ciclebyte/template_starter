package apikey

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/ciclebyte/template_starter/api/v1/apikey"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sApiKey struct{}

func init() {
	service.RegisterApiKey(New())
}

func New() service.IApiKey {
	return &sApiKey{}
}

// ============================================================================
// 管理员 API Key 管理
// ============================================================================

// GetApiKeys 获取API Key列表
func (s *sApiKey) GetApiKeys(ctx context.Context, req *apikey.GetApiKeysReq) (*apikey.GetApiKeysRes, error) {
	var (
		model = dao.ApiKeys.Ctx(ctx)
		err   error
	)

	// 搜索条件
	if req.Search != "" {
		model = model.WhereLike("name", "%"+req.Search+"%")
	}
	if req.Status != nil {
		model = model.Where("status", *req.Status)
	}

	// 获取总数
	total, err := model.Count()
	if err != nil {
		return nil, err
	}

	// 分页查询
	var apiKeys []entity.ApiKeys
	err = model.Page(req.Page, req.Size).
		OrderDesc("created_at").
		Scan(&apiKeys)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	list := make([]apikey.ApiKey, 0, len(apiKeys))
	for _, item := range apiKeys {
		list = append(list, s.convertToApiKeyResponse(item))
	}

	return &apikey.GetApiKeysRes{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// CreateApiKey 创建API Key
func (s *sApiKey) CreateApiKey(ctx context.Context, req *apikey.CreateApiKeyReq) (*apikey.CreateApiKeyRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 生成API Key ID和Secret
	keyId, keySecret, hashedSecret, err := s.generateApiKey()
	if err != nil {
		return nil, err
	}

	// 转换权限为JSON
	permissionsJson := ""
	if len(req.Permissions) > 0 {
		permissionsJson = gconv.String(req.Permissions)
	}

	// 插入数据库
	id, err := dao.ApiKeys.Ctx(ctx).Data(do.ApiKeys{
		UserId:      userId,
		Name:        req.Name,
		KeyId:       keyId,
		KeySecret:   hashedSecret,
		Permissions: permissionsJson,
		ExpiresAt:   req.ExpiresAt,
		Status:      1,
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}

	// 获取创建的记录
	var apiKeyEntity entity.ApiKeys
	err = dao.ApiKeys.Ctx(ctx).Where("id", id).Scan(&apiKeyEntity)
	if err != nil {
		return nil, err
	}

	apiKeyResp := s.convertToApiKeyResponse(apiKeyEntity)

	return &apikey.CreateApiKeyRes{
		ApiKey:    apiKeyResp,
		KeySecret: keySecret, // 只在创建时返回明文
	}, nil
}

// UpdateApiKey 更新API Key
func (s *sApiKey) UpdateApiKey(ctx context.Context, req *apikey.UpdateApiKeyReq) (*apikey.UpdateApiKeyRes, error) {
	// 转换权限为JSON
	permissionsJson := ""
	if len(req.Permissions) > 0 {
		permissionsJson = gconv.String(req.Permissions)
	}

	// 更新数据库
	_, err := dao.ApiKeys.Ctx(ctx).Data(do.ApiKeys{
		Name:        req.Name,
		Permissions: permissionsJson,
		ExpiresAt:   req.ExpiresAt,
		Status:      req.Status,
		UpdatedAt:   gtime.Now(),
	}).Where("id", req.Id).Update()

	if err != nil {
		return nil, err
	}

	return &apikey.UpdateApiKeyRes{}, nil
}

// DeleteApiKey 删除API Key
func (s *sApiKey) DeleteApiKey(ctx context.Context, req *apikey.DeleteApiKeyReq) (*apikey.DeleteApiKeyRes, error) {
	_, err := dao.ApiKeys.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}

	return &apikey.DeleteApiKeyRes{}, nil
}

// RegenerateApiKey 重新生成API Key Secret
func (s *sApiKey) RegenerateApiKey(ctx context.Context, req *apikey.RegenerateApiKeyReq) (*apikey.RegenerateApiKeyRes, error) {
	// 生成新的Secret
	_, keySecret, hashedSecret, err := s.generateApiKey()
	if err != nil {
		return nil, err
	}

	// 更新数据库
	_, err = dao.ApiKeys.Ctx(ctx).Data(do.ApiKeys{
		KeySecret: hashedSecret,
		UpdatedAt: gtime.Now(),
	}).Where("id", req.Id).Update()
	if err != nil {
		return nil, err
	}

	return &apikey.RegenerateApiKeyRes{
		KeySecret: keySecret,
	}, nil
}

// ============================================================================
// 个人 API Key 管理
// ============================================================================

// GetMyApiKeys 获取我的API Key列表
func (s *sApiKey) GetMyApiKeys(ctx context.Context, req *apikey.GetMyApiKeysReq) (*apikey.GetMyApiKeysRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	var (
		model = dao.ApiKeys.Ctx(ctx).Where("user_id", userId)
		err   error
	)

	// 搜索条件
	if req.Search != "" {
		model = model.WhereLike("name", "%"+req.Search+"%")
	}
	if req.Status != nil {
		model = model.Where("status", *req.Status)
	}

	// 获取总数
	total, err := model.Count()
	if err != nil {
		return nil, err
	}

	// 分页查询
	var apiKeys []entity.ApiKeys
	err = model.Page(req.Page, req.Size).
		OrderDesc("created_at").
		Scan(&apiKeys)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	list := make([]apikey.ApiKey, 0, len(apiKeys))
	for _, item := range apiKeys {
		list = append(list, s.convertToApiKeyResponse(item))
	}

	return &apikey.GetMyApiKeysRes{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// CreateMyApiKey 创建我的API Key
func (s *sApiKey) CreateMyApiKey(ctx context.Context, req *apikey.CreateMyApiKeyReq) (*apikey.CreateMyApiKeyRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 生成API Key ID和Secret
	keyId, keySecret, hashedSecret, err := s.generateApiKey()
	if err != nil {
		return nil, err
	}

	// 转换权限为JSON
	permissionsJson := ""
	if len(req.Permissions) > 0 {
		permissionsJson = gconv.String(req.Permissions)
	}

	// 插入数据库
	id, err := dao.ApiKeys.Ctx(ctx).Data(do.ApiKeys{
		UserId:      userId,
		Name:        req.Name,
		KeyId:       keyId,
		KeySecret:   hashedSecret,
		Permissions: permissionsJson,
		ExpiresAt:   req.ExpiresAt,
		Status:      1,
		CreatedAt:   gtime.Now(),
		UpdatedAt:   gtime.Now(),
	}).InsertAndGetId()
	if err != nil {
		return nil, err
	}

	// 获取创建的记录
	var apiKeyEntity entity.ApiKeys
	err = dao.ApiKeys.Ctx(ctx).Where("id", id).Scan(&apiKeyEntity)
	if err != nil {
		return nil, err
	}

	apiKeyResp := s.convertToApiKeyResponse(apiKeyEntity)

	return &apikey.CreateMyApiKeyRes{
		ApiKey:    apiKeyResp,
		KeySecret: keySecret, // 只在创建时返回明文
	}, nil
}

// UpdateMyApiKey 更新我的API Key
func (s *sApiKey) UpdateMyApiKey(ctx context.Context, req *apikey.UpdateMyApiKeyReq) (*apikey.UpdateMyApiKeyRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 转换权限为JSON
	permissionsJson := ""
	if len(req.Permissions) > 0 {
		permissionsJson = gconv.String(req.Permissions)
	}

	// 更新数据库（只能更新自己的API Key）
	_, err := dao.ApiKeys.Ctx(ctx).Data(do.ApiKeys{
		Name:        req.Name,
		Permissions: permissionsJson,
		ExpiresAt:   req.ExpiresAt,
		Status:      req.Status,
		UpdatedAt:   gtime.Now(),
	}).Where("id", req.Id).Where("user_id", userId).Update()

	if err != nil {
		return nil, err
	}

	return &apikey.UpdateMyApiKeyRes{}, nil
}

// DeleteMyApiKey 删除我的API Key
func (s *sApiKey) DeleteMyApiKey(ctx context.Context, req *apikey.DeleteMyApiKeyReq) (*apikey.DeleteMyApiKeyRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 删除（只能删除自己的API Key）
	_, err := dao.ApiKeys.Ctx(ctx).Where("id", req.Id).Where("user_id", userId).Delete()
	if err != nil {
		return nil, err
	}

	return &apikey.DeleteMyApiKeyRes{}, nil
}

// RegenerateMyApiKey 重新生成我的API Key Secret
func (s *sApiKey) RegenerateMyApiKey(ctx context.Context, req *apikey.RegenerateMyApiKeyReq) (*apikey.RegenerateMyApiKeyRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 生成新的Secret
	_, keySecret, hashedSecret, err := s.generateApiKey()
	if err != nil {
		return nil, err
	}

	// 更新数据库（只能更新自己的API Key）
	_, err = dao.ApiKeys.Ctx(ctx).Data(do.ApiKeys{
		KeySecret: hashedSecret,
		UpdatedAt: gtime.Now(),
	}).Where("id", req.Id).Where("user_id", userId).Update()
	if err != nil {
		return nil, err
	}

	return &apikey.RegenerateMyApiKeyRes{
		KeySecret: keySecret,
	}, nil
}

// GetApiKeyLogs 获取API Key使用日志
func (s *sApiKey) GetApiKeyLogs(ctx context.Context, req *apikey.GetApiKeyLogsReq) (*apikey.GetApiKeyLogsRes, error) {
	var (
		model = dao.ApiKeyLogs.Ctx(ctx).Where("api_key_id", req.Id)
		err   error
	)

	// 筛选条件
	if req.Method != "" {
		model = model.Where("method", req.Method)
	}
	if req.Path != "" {
		model = model.WhereLike("path", "%"+req.Path+"%")
	}

	// 获取总数
	total, err := model.Count()
	if err != nil {
		return nil, err
	}

	// 分页查询
	var logs []entity.ApiKeyLogs
	err = model.Page(req.Page, req.Size).
		OrderDesc("created_at").
		Scan(&logs)
	if err != nil {
		return nil, err
	}

	// 转换为API响应格式
	list := make([]apikey.ApiKeyLog, 0, len(logs))
	for _, item := range logs {
		list = append(list, apikey.ApiKeyLog{
			Id:           int64(item.Id),
			ApiKeyId:     int64(item.ApiKeyId),
			UserId:       int64(item.UserId),
			Method:       item.Method,
			Path:         item.Path,
			Ip:           item.Ip,
			UserAgent:    item.UserAgent,
			StatusCode:   item.StatusCode,
			ResponseTime: item.ResponseTime,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &apikey.GetApiKeyLogsRes{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// generateApiKey 生成API Key ID和Secret
func (s *sApiKey) generateApiKey() (keyId, keySecret, hashedSecret string, err error) {
	// 生成32字节随机数据作为Key ID
	keyIdBytes := make([]byte, 16)
	if _, err = rand.Read(keyIdBytes); err != nil {
		return
	}
	keyId = "ak_" + hex.EncodeToString(keyIdBytes)

	// 生成64字节随机数据作为Secret
	secretBytes := make([]byte, 32)
	if _, err = rand.Read(secretBytes); err != nil {
		return
	}
	keySecret = "sk_" + hex.EncodeToString(secretBytes)

	// 对Secret进行哈希存储
	hashedSecret = gmd5.MustEncryptString(keySecret)

	return
}

// convertToApiKeyResponse 转换entity为API响应格式
func (s *sApiKey) convertToApiKeyResponse(item entity.ApiKeys) apikey.ApiKey {
	// 解析权限JSON
	var permissions []string
	if item.Permissions != "" {
		gconv.Scan(item.Permissions, &permissions)
	}

	return apikey.ApiKey{
		Id:          int64(item.Id),
		UserId:      int64(item.UserId),
		Name:        item.Name,
		KeyId:       item.KeyId,
		KeySecret:   "", // 不返回Secret
		Permissions: permissions,
		LastUsedAt:  item.LastUsedAt,
		LastUsedIp:  item.LastUsedIp,
		ExpiresAt:   item.ExpiresAt,
		Status:      item.Status,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

// VerifyApiKey 验证API Key
func (s *sApiKey) VerifyApiKey(ctx context.Context, keyId, keySecret string) (*entity.ApiKeys, error) {
	if keyId == "" || keySecret == "" {
		return nil, errors.New("API Key不能为空")
	}

	// 查询API Key
	var apiKey entity.ApiKeys
	err := dao.ApiKeys.Ctx(ctx).Where("key_id", keyId).Where("status", 1).Scan(&apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("无效的API Key")
		}
		return nil, err
	}

	// 验证Secret
	if gmd5.MustEncryptString(keySecret) != apiKey.KeySecret {
		return nil, errors.New("无效的API Key Secret")
	}

	// 检查是否过期
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(gtime.Now()) {
		return nil, errors.New("API Key已过期")
	}

	return &apiKey, nil
}

// UpdateApiKeyUsage 更新API Key使用信息
func (s *sApiKey) UpdateApiKeyUsage(ctx context.Context, apiKeyId int64, ip string) error {
	_, err := dao.ApiKeys.Ctx(ctx).Data(do.ApiKeys{
		LastUsedAt: gtime.Now(),
		LastUsedIp: ip,
		UpdatedAt:  gtime.Now(),
	}).Where("id", apiKeyId).Update()

	return err
}

// LogApiKeyUsage 记录API Key使用日志
func (s *sApiKey) LogApiKeyUsage(ctx context.Context, apiKeyId, userId int64, method, path, ip, userAgent string, statusCode, responseTime int) error {
	_, err := dao.ApiKeyLogs.Ctx(ctx).Data(do.ApiKeyLogs{
		ApiKeyId:     apiKeyId,
		UserId:       userId,
		Method:       method,
		Path:         path,
		Ip:           ip,
		UserAgent:    userAgent,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		CreatedAt:    gtime.Now(),
	}).Insert()

	return err
}
