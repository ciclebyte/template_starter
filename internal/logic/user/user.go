package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ciclebyte/template_starter/api/v1/user"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/ciclebyte/template_starter/library/libPassword"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sUser struct{}

func init() {
	service.RegisterUser(New())
}

func New() service.IUser {
	return &sUser{}
}

// ============================================================================
// 用户管理
// ============================================================================

// ListUsers 获取用户列表
func (s *sUser) ListUsers(ctx context.Context, req *user.ListUsersReq) (*user.ListUsersRes, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}

	query := dao.Users.Ctx(ctx)

	// 搜索条件
	if req.Search != "" {
		query = query.WhereLike("username", "%"+req.Search+"%").
			WhereOrLike("email", "%"+req.Search+"%").
			WhereOrLike("nickname", "%"+req.Search+"%")
	}

	// 状态过滤
	statusParam := g.RequestFromCtx(ctx).Get("status").String()
	if statusParam != "" && req.Status != nil {
		query = query.Where("status", *req.Status)
	}

	// 角色过滤
	if req.Role != "" {
		query = query.LeftJoin("user_roles ur", "users.id = ur.user_id").
			LeftJoin("roles r", "ur.role_id = r.id").
			Where("r.code", req.Role)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		g.Log().Error(ctx, "get users count failed:", err)
		return nil, err
	}

	// 获取列表
	var users []entity.Users
	err = query.Page(page, size).OrderDesc("created_at").Scan(&users)
	if err != nil {
		g.Log().Error(ctx, "get users failed:", err)
		return nil, err
	}

	// 转换数据并获取用户角色
	list := make([]user.UserInfo, 0, len(users))
	for _, u := range users {
		// 获取用户角色
		roles, _ := s.getUserRoles(ctx, u.Id)

		list = append(list, user.UserInfo{
			Id:          u.Id,
			Username:    u.Username,
			Email:       u.Email,
			Nickname:    u.Nickname,
			Avatar:      u.Avatar,
			Phone:       u.Phone,
			Status:      u.Status,
			LastLoginAt: u.LastLoginAt,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
			Roles:       roles,
		})
	}

	return &user.ListUsersRes{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// CreateUser 创建用户
func (s *sUser) CreateUser(ctx context.Context, req *user.CreateUserReq) (*user.CreateUserRes, error) {
	// 检查用户名是否已存在
	count, err := dao.Users.Ctx(ctx).Where("username", req.Username).Count()
	if err != nil {
		g.Log().Error(ctx, "check username failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	count, err = dao.Users.Ctx(ctx).Where("email", req.Email).Count()
	if err != nil {
		g.Log().Error(ctx, "check email failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := libPassword.HashPassword(req.Password)
	if err != nil {
		g.Log().Error(ctx, "hash password failed:", err)
		return nil, errors.New("密码加密失败")
	}

	// 开始事务
	var result *user.CreateUserRes
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建用户
		userId, err := dao.Users.Ctx(ctx).TX(tx).Data(do.Users{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hashedPassword,
			Nickname:     req.Nickname,
			Phone:        req.Phone,
			Status:       req.Status,
		}).InsertAndGetId()

		if err != nil {
			g.Log().Error(ctx, "create user failed:", err)
			return errors.New("创建用户失败")
		}

		// 分配角色
		if len(req.Roles) > 0 {
			err = s.assignUserRoles(ctx, tx, userId, req.Roles, "")
			if err != nil {
				return err
			}
		}

		result = &user.CreateUserRes{Id: userId}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateUser 更新用户
func (s *sUser) UpdateUser(ctx context.Context, req *user.UpdateUserReq) (*user.UpdateUserRes, error) {
	// 检查用户是否存在
	exists, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check user exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("用户不存在")
	}

	// 检查用户名是否被其他用户使用
	count, err := dao.Users.Ctx(ctx).Where("username", req.Username).WhereNot("id", req.Id).Count()
	if err != nil {
		g.Log().Error(ctx, "check username failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已被使用")
	}

	// 检查邮箱是否被其他用户使用
	count, err = dao.Users.Ctx(ctx).Where("email", req.Email).WhereNot("id", req.Id).Count()
	if err != nil {
		g.Log().Error(ctx, "check email failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("邮箱已被使用")
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新用户
		_, err = dao.Users.Ctx(ctx).TX(tx).Data(do.Users{
			Username:  req.Username,
			Email:     req.Email,
			Nickname:  req.Nickname,
			Phone:     req.Phone,
			Status:    req.Status,
			UpdatedAt: gtime.Now(),
		}).Where("id", req.Id).Update()

		if err != nil {
			g.Log().Error(ctx, "update user failed:", err)
			return errors.New("更新用户失败")
		}

		// 删除原有角色
		_, err = dao.UserRoles.Ctx(ctx).TX(tx).Where("user_id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete user roles failed:", err)
			return errors.New("更新用户角色失败")
		}

		// 重新分配角色
		if len(req.Roles) > 0 {
			err = s.assignUserRoles(ctx, tx, req.Id, req.Roles, "")
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user.UpdateUserRes{}, nil
}

// DeleteUser 删除用户
func (s *sUser) DeleteUser(ctx context.Context, req *user.DeleteUserReq) (*user.DeleteUserRes, error) {
	// 检查用户是否存在
	exists, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check user exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("用户不存在")
	}

	// 开始事务删除用户及其关联数据
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除用户角色关联
		_, err = dao.UserRoles.Ctx(ctx).TX(tx).Where("user_id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete user roles failed:", err)
			return errors.New("删除用户角色失败")
		}

		// 删除用户会话
		_, err = dao.UserSessions.Ctx(ctx).TX(tx).Where("user_id", req.Id).Delete()
		if err != nil {
			g.Log().Warning(ctx, "delete user sessions failed:", err)
		}

		// 删除用户
		_, err = dao.Users.Ctx(ctx).TX(tx).Where("id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete user failed:", err)
			return errors.New("删除用户失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user.DeleteUserRes{}, nil
}

// GetUser 获取用户详情
func (s *sUser) GetUser(ctx context.Context, req *user.GetUserReq) (*user.GetUserRes, error) {
	var u entity.Users
	err := dao.Users.Ctx(ctx).Where("id", req.Id).Scan(&u)
	if err != nil {
		g.Log().Error(ctx, "get user failed:", err)
		return nil, err
	}

	if u.Id == 0 {
		return nil, errors.New("用户不存在")
	}

	// 获取用户角色
	roles, err := s.getUserRoles(ctx, u.Id)
	if err != nil {
		g.Log().Warning(ctx, "get user roles failed:", err)
		roles = []user.RoleInfo{}
	}

	return &user.GetUserRes{
		UserInfo: &user.UserInfo{
			Id:          u.Id,
			Username:    u.Username,
			Email:       u.Email,
			Nickname:    u.Nickname,
			Avatar:      u.Avatar,
			Phone:       u.Phone,
			Status:      u.Status,
			LastLoginAt: u.LastLoginAt,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
			Roles:       roles,
		},
	}, nil
}

// ResetPassword 重置用户密码
func (s *sUser) ResetPassword(ctx context.Context, req *user.ResetPasswordReq) (*user.ResetPasswordRes, error) {
	// 检查用户是否存在
	exists, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check user exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("用户不存在")
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
	}).Where("id", req.Id).Update()

	if err != nil {
		g.Log().Error(ctx, "reset password failed:", err)
		return nil, errors.New("重置密码失败")
	}

	return &user.ResetPasswordRes{}, nil
}

// UpdateUserStatus 更新用户状态
func (s *sUser) UpdateUserStatus(ctx context.Context, req *user.UpdateUserStatusReq) (*user.UpdateUserStatusRes, error) {
	// 检查用户是否存在
	exists, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check user exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("用户不存在")
	}

	// 更新状态
	_, err = dao.Users.Ctx(ctx).Data(do.Users{
		Status:    req.Status,
		UpdatedAt: gtime.Now(),
	}).Where("id", req.Id).Update()

	if err != nil {
		g.Log().Error(ctx, "update user status failed:", err)
		return nil, errors.New("更新用户状态失败")
	}

	return &user.UpdateUserStatusRes{}, nil
}

// AssignUserRoles 分配用户角色
func (s *sUser) AssignUserRoles(ctx context.Context, req *user.AssignUserRolesReq) (*user.AssignUserRolesRes, error) {
	// 检查用户是否存在
	exists, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check user exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("用户不存在")
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除用户原有角色
		_, err := dao.UserRoles.Ctx(ctx).TX(tx).Where("user_id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete user roles failed:", err)
			return errors.New("删除用户原有角色失败")
		}

		// 分配新角色
		if len(req.Roles) > 0 {
			err = s.assignUserRoles(ctx, tx, req.Id, req.Roles, req.ExpiresAt)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &user.AssignUserRolesRes{}, nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// getUserRoles 获取用户角色列表
func (s *sUser) getUserRoles(ctx context.Context, userId int64) ([]user.RoleInfo, error) {
	results, err := dao.UserRoles.Ctx(ctx).
		LeftJoin("roles r", "user_roles.role_id = r.id").
		Fields("r.id, r.name, r.code, r.description").
		Where("user_roles.user_id", userId).
		Where("r.status", 1). // 只显示启用的角色
		OrderAsc("r.name").
		All()

	if err != nil {
		return nil, err
	}

	roles := make([]user.RoleInfo, 0, len(results))
	for _, record := range results {
		roles = append(roles, user.RoleInfo{
			Id:          record["id"].Int64(),
			Name:        record["name"].String(),
			Code:        record["code"].String(),
			Description: record["description"].String(),
		})
	}

	return roles, nil
}

// assignUserRoles 分配用户角色（事务内使用）
func (s *sUser) assignUserRoles(ctx context.Context, tx gdb.TX, userId int64, roleIds []int64, expiresAt string) error {
	// 获取当前操作用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	assignedBy := int64(1) // 默认系统用户
	if userIdVar != nil {
		assignedBy = gconv.Int64(userIdVar)
	}

	// 解析过期时间
	var expiresTime *gtime.Time
	if expiresAt != "" {
		expTime, err := gtime.StrToTime(expiresAt)
		if err != nil {
			return errors.New("过期时间格式错误")
		}
		expiresTime = expTime
	}

	for _, roleId := range roleIds {
		_, err := dao.UserRoles.Ctx(ctx).TX(tx).Data(do.UserRoles{
			UserId:     userId,
			RoleId:     roleId,
			AssignedBy: assignedBy,
			ExpiresAt:  expiresTime,
		}).Insert()

		if err != nil {
			g.Log().Error(ctx, "assign user role failed:", err)
			return fmt.Errorf("分配角色失败，角色ID: %d", roleId)
		}
	}

	return nil
}