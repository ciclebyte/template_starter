package profile

import (
	"context"
	"errors"

	"github.com/ciclebyte/template_starter/api/v1/profile"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/libPassword"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sProfile struct{}

func init() {
	service.RegisterProfile(New())
}

func New() service.IProfile {
	return &sProfile{}
}

// ============================================================================
// 个人资料管理
// ============================================================================

// GetProfile 获取个人资料
func (s *sProfile) GetProfile(ctx context.Context, req *profile.GetProfileReq) (*profile.GetProfileRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 获取用户基本信息
	var user entity.Users
	err := dao.Users.Ctx(ctx).Where("id", userId).Scan(&user)
	if err != nil {
		g.Log().Error(ctx, "get user profile failed:", err)
		return nil, err
	}

	if user.Id == 0 {
		return nil, errors.New("用户不存在")
	}

	// 获取用户角色
	roles, err := s.getUserRoles(ctx, userId)
	if err != nil {
		g.Log().Warning(ctx, "get user roles failed:", err)
		roles = []profile.UserRole{}
	}

	// 获取用户权限
	permissions, err := s.getUserPermissions(ctx, userId)
	if err != nil {
		g.Log().Warning(ctx, "get user permissions failed:", err)
		permissions = []profile.UserPermission{}
	}

	return &profile.GetProfileRes{
		Profile: &profile.UserProfile{
			Id:          user.Id,
			Username:    user.Username,
			Email:       user.Email,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Phone:       user.Phone,
			Status:      user.Status,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		},
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

// UpdateProfile 更新个人资料
func (s *sProfile) UpdateProfile(ctx context.Context, req *profile.UpdateProfileReq) (*profile.UpdateProfileRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 更新用户信息
	_, err := dao.Users.Ctx(ctx).Data(do.Users{
		Nickname:  req.Nickname,
		Avatar:    req.Avatar,
		Phone:     req.Phone,
		UpdatedAt: gtime.Now(),
	}).Where("id", userId).Update()

	if err != nil {
		g.Log().Error(ctx, "update user profile failed:", err)
		return nil, errors.New("更新个人资料失败")
	}

	return &profile.UpdateProfileRes{}, nil
}

// ChangePassword 修改密码
func (s *sProfile) ChangePassword(ctx context.Context, req *profile.ChangePasswordReq) (*profile.ChangePasswordRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 获取用户当前密码哈希
	var user entity.Users
	err := dao.Users.Ctx(ctx).Fields("password_hash").Where("id", userId).Scan(&user)
	if err != nil {
		g.Log().Error(ctx, "get user password failed:", err)
		return nil, err
	}

	if user.Id == 0 {
		return nil, errors.New("用户不存在")
	}

	// 验证原密码
	match, err := libPassword.VerifyPassword(req.OldPassword, user.PasswordHash)
	if err != nil || !match {
		return nil, errors.New("原密码不正确")
	}

	// 加密新密码
	hashedPassword, err := libPassword.HashPassword(req.NewPassword)
	if err != nil {
		g.Log().Error(ctx, "hash password failed:", err)
		return nil, errors.New("密码加密失败")
	}

	// 更新密码
	_, err = dao.Users.Ctx(ctx).Data(do.Users{
		PasswordHash: hashedPassword,
		UpdatedAt:    gtime.Now(),
	}).Where("id", userId).Update()

	if err != nil {
		g.Log().Error(ctx, "update password failed:", err)
		return nil, errors.New("修改密码失败")
	}

	return &profile.ChangePasswordRes{}, nil
}

// UpdateEmail 更新邮箱
func (s *sProfile) UpdateEmail(ctx context.Context, req *profile.UpdateEmailReq) (*profile.UpdateEmailRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 验证当前密码
	var user entity.Users
	err := dao.Users.Ctx(ctx).Fields("password_hash").Where("id", userId).Scan(&user)
	if err != nil {
		g.Log().Error(ctx, "get user password failed:", err)
		return nil, err
	}

	if user.Id == 0 {
		return nil, errors.New("用户不存在")
	}

	match, err := libPassword.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !match {
		return nil, errors.New("密码不正确")
	}

	// 检查新邮箱是否已被使用
	count, err := dao.Users.Ctx(ctx).Where("email", req.Email).WhereNot("id", userId).Count()
	if err != nil {
		g.Log().Error(ctx, "check email failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已被使用")
	}

	// 更新邮箱
	_, err = dao.Users.Ctx(ctx).Data(do.Users{
		Email:         req.Email,
		EmailVerified: 0, // 重置邮箱验证状态
		UpdatedAt:     gtime.Now(),
	}).Where("id", userId).Update()

	if err != nil {
		g.Log().Error(ctx, "update email failed:", err)
		return nil, errors.New("更新邮箱失败")
	}

	return &profile.UpdateEmailRes{}, nil
}

// UploadAvatar 上传头像
func (s *sProfile) UploadAvatar(ctx context.Context, req *profile.UploadAvatarReq) (*profile.UploadAvatarRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 更新头像
	_, err := dao.Users.Ctx(ctx).Data(do.Users{
		Avatar:    req.AvatarUrl,
		UpdatedAt: gtime.Now(),
	}).Where("id", userId).Update()

	if err != nil {
		g.Log().Error(ctx, "update avatar failed:", err)
		return nil, errors.New("更新头像失败")
	}

	return &profile.UploadAvatarRes{
		AvatarUrl: req.AvatarUrl,
	}, nil
}

// GetSecurityInfo 获取账户安全信息
func (s *sProfile) GetSecurityInfo(ctx context.Context, req *profile.GetSecurityInfoReq) (*profile.GetSecurityInfoRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	// 获取用户安全信息
	var user entity.Users
	err := dao.Users.Ctx(ctx).Where("id", userId).Scan(&user)
	if err != nil {
		g.Log().Error(ctx, "get user security info failed:", err)
		return nil, err
	}

	if user.Id == 0 {
		return nil, errors.New("用户不存在")
	}

	return &profile.GetSecurityInfoRes{
		SecurityInfo: &profile.SecurityInfo{
			EmailVerified:      user.EmailVerified == 1,
			PhoneVerified:      false, // 暂时硬编码，后续可添加手机验证功能
			TwoFactorEnabled:   false, // 暂时硬编码，后续可添加双因子认证
			LastLoginAt:        user.LastLoginAt,
			LastLoginIp:        user.LastLoginIp,
			LoginCount:         user.LoginCount,
			PasswordUpdatedAt:  user.UpdatedAt, // 用更新时间代替密码更新时间
		},
	}, nil
}

// GetLoginHistory 获取登录历史
func (s *sProfile) GetLoginHistory(ctx context.Context, req *profile.GetLoginHistoryReq) (*profile.GetLoginHistoryRes, error) {
	// 获取当前用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未登录")
	}
	userId := gconv.Int64(userIdVar)

	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}

	// 获取登录历史总数
	total, err := dao.UserSessions.Ctx(ctx).Where("user_id", userId).Count()
	if err != nil {
		g.Log().Error(ctx, "get login history count failed:", err)
		return nil, err
	}

	// 获取登录历史列表
	var sessions []entity.UserSessions
	err = dao.UserSessions.Ctx(ctx).
		Where("user_id", userId).
		Page(page, size).
		OrderDesc("created_at").
		Scan(&sessions)

	if err != nil {
		g.Log().Error(ctx, "get login history failed:", err)
		return nil, err
	}

	// 转换数据
	list := make([]profile.LoginHistory, 0, len(sessions))
	for _, session := range sessions {
		list = append(list, profile.LoginHistory{
			Id:        session.Id,
			LoginAt:   session.CreatedAt,
			LoginIp:   session.IpAddress,
			UserAgent: session.UserAgent,
			Location:  "", // 暂时为空，后续可以根据IP获取地理位置
			Status:    "成功", // 暂时硬编码，后续可以添加登录状态
		})
	}

	return &profile.GetLoginHistoryRes{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// getUserRoles 获取用户角色列表
func (s *sProfile) getUserRoles(ctx context.Context, userId int64) ([]profile.UserRole, error) {
	results, err := dao.UserRoles.Ctx(ctx).
		LeftJoin("roles r", "user_roles.role_id = r.id").
		Fields("r.id, r.name, r.code, r.description, user_roles.assigned_at, user_roles.expires_at").
		Where("user_roles.user_id", userId).
		Where("r.status", 1). // 只显示启用的角色
		OrderAsc("r.name").
		All()

	if err != nil {
		return nil, err
	}

	roles := make([]profile.UserRole, 0, len(results))
	for _, record := range results {
		roles = append(roles, profile.UserRole{
			Id:          record["id"].Int64(),
			Name:        record["name"].String(),
			Code:        record["code"].String(),
			Description: record["description"].String(),
			AssignedAt:  record["assigned_at"].GTime(),
			ExpiresAt:   record["expires_at"].GTime(),
		})
	}

	return roles, nil
}

// getUserPermissions 获取用户权限列表
func (s *sProfile) getUserPermissions(ctx context.Context, userId int64) ([]profile.UserPermission, error) {
	results, err := dao.Users.Ctx(ctx).
		LeftJoin("user_roles ur", "users.id = ur.user_id").
		LeftJoin("roles r", "ur.role_id = r.id").
		LeftJoin("role_permissions rp", "r.id = rp.role_id").
		LeftJoin("permissions p", "rp.permission_id = p.id").
		Fields("DISTINCT p.id, p.name, p.code, p.resource, p.action, p.description").
		Where("users.id", userId).
		Where("users.status", 1).
		Where("r.status", 1).
		Where("ur.expires_at IS NULL OR ur.expires_at > NOW()").
		OrderAsc("p.resource, p.action").
		All()

	if err != nil {
		return nil, err
	}

	permissions := make([]profile.UserPermission, 0, len(results))
	for _, record := range results {
		// 跳过空记录
		if record["id"].Int64() == 0 {
			continue
		}
		permissions = append(permissions, profile.UserPermission{
			Id:          record["id"].Int64(),
			Name:        record["name"].String(),
			Code:        record["code"].String(),
			Resource:    record["resource"].String(),
			Action:      record["action"].String(),
			Description: record["description"].String(),
		})
	}

	return permissions, nil
}