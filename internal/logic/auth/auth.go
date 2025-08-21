package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/libJWT"
	"github.com/ciclebyte/template_starter/library/libPassword"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gconv"
)

type sAuth struct{}

func init() {
	service.RegisterAuth(New())
}

func New() service.IAuth {
	return &sAuth{}
}

// Register 用户注册
func (s *sAuth) Register(ctx context.Context, req *service.RegisterReq) (*service.RegisterRes, error) {
	// 检查用户名是否已存在
	count, err := dao.Users.Ctx(ctx).Where("username", req.Username).Count()
	if err != nil {
		g.Log().Error(ctx, "check username exists failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	count, err = dao.Users.Ctx(ctx).Where("email", req.Email).Count()
	if err != nil {
		g.Log().Error(ctx, "check email exists failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// 验证密码强度
	if valid, msg := libPassword.IsValidPassword(req.Password); !valid {
		return nil, errors.New(msg)
	}

	// 加密密码
	passwordHash, err := libPassword.HashPassword(req.Password)
	if err != nil {
		g.Log().Error(ctx, "hash password failed:", err)
		return nil, errors.New("密码加密失败")
	}

	// 开始事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (*service.RegisterRes, error) {
		// 创建用户
		nickname := req.Nickname
		if nickname == "" {
			nickname = req.Username
		}

		userId, err := dao.Users.Ctx(ctx).TX(tx).Data(do.Users{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: passwordHash,
			Nickname:     nickname,
			Status:       1,
		}).InsertAndGetId()

		if err != nil {
			g.Log().Error(ctx, "create user failed:", err)
			return nil, errors.New("创建用户失败")
		}

		// 分配默认角色 (普通用户)
		var role *entity.Roles
		err = dao.Roles.Ctx(ctx).TX(tx).Where("code", "user").Scan(&role)
		if err != nil {
			g.Log().Error(ctx, "find default role failed:", err)
			return nil, errors.New("分配默认角色失败")
		}

		if role != nil {
			_, err = dao.UserRoles.Ctx(ctx).TX(tx).Data(do.UserRoles{
				UserId:     userId,
				RoleId:     role.Id,
				AssignedBy: userId, // 自己分配给自己
			}).Insert()

			if err != nil {
				g.Log().Error(ctx, "assign default role failed:", err)
				return nil, errors.New("分配默认角色失败")
			}
		}

		// 获取用户完整信息
		userInfo, err := s.getUserInfoById(ctx, tx, userId)
		if err != nil {
			return nil, err
		}

		// 生成令牌
		tokenInfo, err := libJWT.GetManager().GenerateTokens(
			userInfo.ID, 
			userInfo.Username, 
			userInfo.Email,
			userInfo.Roles,
			userInfo.Permissions,
		)
		if err != nil {
			g.Log().Error(ctx, "generate tokens failed:", err)
			return nil, errors.New("生成令牌失败")
		}

		// 记录会话
		err = s.createUserSession(ctx, tx, userInfo.ID, tokenInfo.AccessToken)
		if err != nil {
			g.Log().Warning(ctx, "create user session failed:", err)
		}

		return &service.RegisterRes{
			TokenInfo: tokenInfo,
			User:      userInfo,
		}, nil
	})
}

// Login 用户登录
func (s *sAuth) Login(ctx context.Context, req *service.LoginReq) (*service.LoginRes, error) {
	// 查找用户
	var user *entity.Users
	err := dao.Users.Ctx(ctx).Where(g.Map{
		"username": req.Username,
		"status":   1, // 只查找正常状态的用户
	}).Scan(&user)

	if err != nil {
		g.Log().Error(ctx, "find user failed:", err)
		return nil, errors.New("查找用户失败")
	}

	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	valid, err := libPassword.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		g.Log().Error(ctx, "verify password failed:", err)
		return nil, errors.New("密码验证失败")
	}

	if !valid {
		return nil, errors.New("用户名或密码错误")
	}

	// 开始事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (*service.LoginRes, error) {
		// 更新登录信息
		_, err = dao.Users.Ctx(ctx).TX(tx).Data(g.Map{
			"last_login_at": gconv.String(time.Now()),
			"last_login_ip": g.RequestFromCtx(ctx).GetClientIp(),
			"login_count":   gdb.Raw("login_count + 1"),
		}).Where("id", user.Id).Update()

		if err != nil {
			g.Log().Warning(ctx, "update login info failed:", err)
		}

		// 获取用户完整信息
		userInfo, err := s.getUserInfoById(ctx, tx, user.Id)
		if err != nil {
			return nil, err
		}

		// 生成令牌
		tokenInfo, err := libJWT.GetManager().GenerateTokens(
			userInfo.ID, 
			userInfo.Username, 
			userInfo.Email,
			userInfo.Roles,
			userInfo.Permissions,
		)
		if err != nil {
			g.Log().Error(ctx, "generate tokens failed:", err)
			return nil, errors.New("生成令牌失败")
		}

		// 记录会话
		err = s.createUserSession(ctx, tx, userInfo.ID, tokenInfo.AccessToken)
		if err != nil {
			g.Log().Warning(ctx, "create user session failed:", err)
		}

		return &service.LoginRes{
			TokenInfo: tokenInfo,
			User:      userInfo,
		}, nil
	})
}

// Logout 用户登出
func (s *sAuth) Logout(ctx context.Context) error {
	// 从上下文获取用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil // 已经是未登录状态
	}

	userId := gconv.Int64(userIdVar)
	if userId == 0 {
		return nil
	}

	// 删除用户会话
	_, err := dao.UserSessions.Ctx(ctx).Where("user_id", userId).Delete()
	if err != nil {
		g.Log().Warning(ctx, "delete user session failed:", err)
	}

	return nil
}

// GetCurrentUser 获取当前用户信息
func (s *sAuth) GetCurrentUser(ctx context.Context) (*service.AuthUserInfo, error) {
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}

	userId := gconv.Int64(userIdVar)
	if userId == 0 {
		return nil, errors.New("用户ID无效")
	}

	return s.getUserInfoById(ctx, nil, userId)
}

// RefreshToken 刷新Token
func (s *sAuth) RefreshToken(ctx context.Context, req *service.RefreshTokenReq) (*service.RefreshTokenRes, error) {
	// 验证刷新令牌并获取用户信息
	claims, err := libJWT.GetManager().ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("刷新令牌无效")
	}

	// 获取用户最新的角色和权限
	roles, err := s.GetUserRoles(ctx, claims.UserID)
	if err != nil {
		g.Log().Warning(ctx, "get user roles failed:", err)
		roles = []string{}
	}

	permissions, err := s.GetUserPermissions(ctx, claims.UserID)
	if err != nil {
		g.Log().Warning(ctx, "get user permissions failed:", err)
		permissions = []string{}
	}

	// 生成新的令牌
	tokenInfo, err := libJWT.GetManager().RefreshAccessToken(req.RefreshToken, roles, permissions)
	if err != nil {
		g.Log().Error(ctx, "refresh access token failed:", err)
		return nil, errors.New("刷新令牌失败")
	}

	return &service.RefreshTokenRes{
		TokenInfo: tokenInfo,
	}, nil
}

// HasPermission 检查用户权限
func (s *sAuth) HasPermission(ctx context.Context, userId int64, permission string) (bool, error) {
	// 检查是否是超级管理员
	isAdmin, err := s.HasRole(ctx, userId, "super_admin")
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}

	// 查询用户权限
	count, err := dao.Users.Ctx(ctx).
		LeftJoin("user_roles ur", "users.id = ur.user_id").
		LeftJoin("roles r", "ur.role_id = r.id").
		LeftJoin("role_permissions rp", "r.id = rp.role_id").
		LeftJoin("permissions p", "rp.permission_id = p.id").
		Where("users.id", userId).
		Where("users.status", 1).
		Where("r.status", 1).
		Where("p.code", permission).
		Where("ur.expires_at IS NULL OR ur.expires_at > ?", "NOW()").
		Count()

	return count > 0, err
}

// HasRole 检查用户角色
func (s *sAuth) HasRole(ctx context.Context, userId int64, role string) (bool, error) {
	count, err := dao.Users.Ctx(ctx).
		LeftJoin("user_roles ur", "users.id = ur.user_id").
		LeftJoin("roles r", "ur.role_id = r.id").
		Where("users.id", userId).
		Where("users.status", 1).
		Where("r.status", 1).
		Where("r.code", role).
		Where("ur.expires_at IS NULL OR ur.expires_at > ?", "NOW()").
		Count()

	return count > 0, err
}

// GetUserPermissions 获取用户权限列表
func (s *sAuth) GetUserPermissions(ctx context.Context, userId int64) ([]string, error) {
	var permissions []string

	result, err := dao.Users.Ctx(ctx).Fields("DISTINCT p.code").
		LeftJoin("user_roles ur", "users.id = ur.user_id").
		LeftJoin("roles r", "ur.role_id = r.id").
		LeftJoin("role_permissions rp", "r.id = rp.role_id").
		LeftJoin("permissions p", "rp.permission_id = p.id").
		Where("users.id", userId).
		Where("users.status", 1).
		Where("r.status", 1).
		Where("ur.expires_at IS NULL OR ur.expires_at > ?", "NOW()").
		Array()

	if err != nil {
		return nil, err
	}

	for _, perm := range result {
		permissions = append(permissions, gconv.String(perm))
	}

	return permissions, nil
}

// GetUserRoles 获取用户角色列表
func (s *sAuth) GetUserRoles(ctx context.Context, userId int64) ([]string, error) {
	var roles []string

	result, err := dao.Users.Ctx(ctx).Fields("DISTINCT r.code").
		LeftJoin("user_roles ur", "users.id = ur.user_id").
		LeftJoin("roles r", "ur.role_id = r.id").
		Where("users.id", userId).
		Where("users.status", 1).
		Where("r.status", 1).
		Where("ur.expires_at IS NULL OR ur.expires_at > ?", "NOW()").
		Array()

	if err != nil {
		return nil, err
	}

	for _, role := range result {
		roles = append(roles, gconv.String(role))
	}

	return roles, nil
}

// getUserInfoById 根据ID获取用户完整信息
func (s *sAuth) getUserInfoById(ctx context.Context, tx gdb.TX, userId int64) (*service.AuthUserInfo, error) {
	var user *entity.Users
	query := dao.Users.Ctx(ctx)
	if tx != nil {
		query = query.TX(tx)
	}

	err := query.Where("id", userId).Scan(&user)
	if err != nil {
		g.Log().Error(ctx, "find user by id failed:", err)
		return nil, errors.New("获取用户信息失败")
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户角色
	roles, err := s.GetUserRoles(ctx, userId)
	if err != nil {
		g.Log().Warning(ctx, "get user roles failed:", err)
		roles = []string{}
	}

	// 获取用户权限
	permissions, err := s.GetUserPermissions(ctx, userId)
	if err != nil {
		g.Log().Warning(ctx, "get user permissions failed:", err)
		permissions = []string{}
	}

	return &service.AuthUserInfo{
		ID:          user.Id,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      user.Status,
		Roles:       roles,
		Permissions: permissions,
		LastLoginAt: gconv.String(user.LastLoginAt),
	}, nil
}

// createUserSession 创建用户会话记录
func (s *sAuth) createUserSession(ctx context.Context, tx gdb.TX, userId int64, accessToken string) error {
	request := g.RequestFromCtx(ctx)
	sessionId := guid.S()

	_, err := dao.UserSessions.Ctx(ctx).TX(tx).Data(do.UserSessions{
		Id:        sessionId,
		UserId:    userId,
		IpAddress: request.GetClientIp(),
		UserAgent: request.Header.Get("User-Agent"),
		Data:      fmt.Sprintf(`{"access_token_hash": "%s"}`, g.Crypto().Md5String(accessToken)),
		ExpiresAt: gconv.String(time.Now().Add(2 * time.Hour)),
	}).Insert()

	return err
}