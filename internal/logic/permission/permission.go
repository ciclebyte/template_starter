package permission

import (
	"context"
	"errors"

	"github.com/ciclebyte/template_starter/api/v1/permission"
	"github.com/ciclebyte/template_starter/internal/dao"
	"github.com/ciclebyte/template_starter/internal/model/do"
	"github.com/ciclebyte/template_starter/internal/model/entity"
	"github.com/ciclebyte/template_starter/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sPermission struct{}

func init() {
	service.RegisterPermission(New())
}

func New() service.IPermission {
	return &sPermission{}
}

// ============================================================================
// 权限管理
// ============================================================================

// ListPermissions 获取权限列表
func (s *sPermission) ListPermissions(ctx context.Context, req *permission.ListPermissionsReq) (*permission.ListPermissionsRes, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}

	query := dao.Permissions.Ctx(ctx)

	// 搜索条件
	if req.Search != "" {
		query = query.WhereLike("name", "%"+req.Search+"%").
			WhereOrLike("code", "%"+req.Search+"%").
			WhereOrLike("description", "%"+req.Search+"%")
	}

	// 资源类型过滤
	if req.Resource != "" {
		query = query.Where("resource", req.Resource)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		g.Log().Error(ctx, "get permissions count failed:", err)
		return nil, err
	}

	// 获取列表
	var permissions []entity.Permissions
	err = query.Page(page, size).OrderDesc("created_at").Scan(&permissions)
	if err != nil {
		g.Log().Error(ctx, "get permissions failed:", err)
		return nil, err
	}

	// 转换数据
	list := make([]permission.PermissionInfo, 0, len(permissions))
	for _, p := range permissions {
		list = append(list, permission.PermissionInfo{
			Id:          p.Id,
			Name:        p.Name,
			Code:        p.Code,
			Resource:    p.Resource,
			Action:      p.Action,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.CreatedAt, // 使用CreatedAt作为UpdatedAt的替代
		})
	}

	return &permission.ListPermissionsRes{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// CreatePermission 创建权限
func (s *sPermission) CreatePermission(ctx context.Context, req *permission.CreatePermissionReq) (*permission.CreatePermissionRes, error) {
	// 检查权限代码是否已存在
	count, err := dao.Permissions.Ctx(ctx).Where("code", req.Code).Count()
	if err != nil {
		g.Log().Error(ctx, "check permission code failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("权限代码已存在")
	}

	// 创建权限
	id, err := dao.Permissions.Ctx(ctx).Data(do.Permissions{
		Name:        req.Name,
		Code:        req.Code,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}).InsertAndGetId()

	if err != nil {
		g.Log().Error(ctx, "create permission failed:", err)
		return nil, errors.New("创建权限失败")
	}

	return &permission.CreatePermissionRes{Id: id}, nil
}

// UpdatePermission 更新权限
func (s *sPermission) UpdatePermission(ctx context.Context, req *permission.UpdatePermissionReq) (*permission.UpdatePermissionRes, error) {
	// 检查权限是否存在
	exists, err := dao.Permissions.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check permission exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("权限不存在")
	}

	// 检查权限代码是否被其他权限使用
	count, err := dao.Permissions.Ctx(ctx).Where("code", req.Code).WhereNot("id", req.Id).Count()
	if err != nil {
		g.Log().Error(ctx, "check permission code failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("权限代码已被使用")
	}

	// 更新权限
	_, err = dao.Permissions.Ctx(ctx).Data(do.Permissions{
		Name:        req.Name,
		Code:        req.Code,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}).Where("id", req.Id).Update()

	if err != nil {
		g.Log().Error(ctx, "update permission failed:", err)
		return nil, errors.New("更新权限失败")
	}

	return &permission.UpdatePermissionRes{}, nil
}

// DeletePermission 删除权限
func (s *sPermission) DeletePermission(ctx context.Context, req *permission.DeletePermissionReq) (*permission.DeletePermissionRes, error) {
	// 检查权限是否被角色使用
	count, err := dao.RolePermissions.Ctx(ctx).Where("permission_id", req.Id).Count()
	if err != nil {
		g.Log().Error(ctx, "check permission usage failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("权限正在被角色使用，无法删除")
	}

	// 删除权限
	_, err = dao.Permissions.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		g.Log().Error(ctx, "delete permission failed:", err)
		return nil, errors.New("删除权限失败")
	}

	return &permission.DeletePermissionRes{}, nil
}

// GetPermission 获取权限详情
func (s *sPermission) GetPermission(ctx context.Context, req *permission.GetPermissionReq) (*permission.GetPermissionRes, error) {
	var perm entity.Permissions
	err := dao.Permissions.Ctx(ctx).Where("id", req.Id).Scan(&perm)
	if err != nil {
		g.Log().Error(ctx, "get permission failed:", err)
		return nil, err
	}

	if perm.Id == 0 {
		return nil, errors.New("权限不存在")
	}

	return &permission.GetPermissionRes{
		PermissionInfo: &permission.PermissionInfo{
			Id:          perm.Id,
			Name:        perm.Name,
			Code:        perm.Code,
			Resource:    perm.Resource,
			Action:      perm.Action,
			Description: perm.Description,
			CreatedAt:   perm.CreatedAt,
			UpdatedAt:   perm.CreatedAt, // 使用CreatedAt作为UpdatedAt的替代
		},
	}, nil
}

// ============================================================================
// 角色管理
// ============================================================================

// ListRoles 获取角色列表
func (s *sPermission) ListRoles(ctx context.Context, req *permission.ListRolesReq) (*permission.ListRolesRes, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}

	// 添加调试日志
	g.Log().Debug(ctx, "ListRoles request:", req)

	query := dao.Roles.Ctx(ctx)

	// 搜索条件
	if req.Search != "" {
		g.Log().Debug(ctx, "Adding search condition:", req.Search)
		query = query.WhereLike("name", "%"+req.Search+"%").
			WhereOrLike("code", "%"+req.Search+"%").
			WhereOrLike("description", "%"+req.Search+"%")
	}

	// 状态过滤 - 只有明确指定状态值时才过滤
	// 检查原始请求参数，避免空字符串被转换为0的问题
	statusParam := g.RequestFromCtx(ctx).Get("status").String()
	if statusParam != "" && req.Status != nil {
		g.Log().Debug(ctx, "Adding status condition:", *req.Status)
		query = query.Where("status", *req.Status)
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		g.Log().Error(ctx, "get roles count failed:", err)
		return nil, err
	}
	g.Log().Debug(ctx, "Roles total count:", total)

	// 获取列表
	var roles []entity.Roles
	err = query.Page(page, size).OrderDesc("created_at").Scan(&roles)
	if err != nil {
		g.Log().Error(ctx, "get roles failed:", err)
		return nil, err
	}
	g.Log().Debug(ctx, "Retrieved roles count:", len(roles))

	// 转换数据并获取角色权限
	list := make([]permission.RoleInfo, 0, len(roles))
	for _, r := range roles {
		// 获取角色权限
		permissions, _ := s.getRolePermissions(ctx, r.Id)
		
		// 获取角色用户数量
		userCount, _ := dao.UserRoles.Ctx(ctx).Where("role_id", r.Id).Count()

		list = append(list, permission.RoleInfo{
			Id:          r.Id,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			Status:      r.Status,
			UserCount:   userCount,
			Permissions: permissions,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		})
	}

	return &permission.ListRolesRes{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}

// CreateRole 创建角色
func (s *sPermission) CreateRole(ctx context.Context, req *permission.CreateRoleReq) (*permission.CreateRoleRes, error) {
	// 检查角色代码是否已存在
	count, err := dao.Roles.Ctx(ctx).Where("code", req.Code).Count()
	if err != nil {
		g.Log().Error(ctx, "check role code failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("角色代码已存在")
	}

	// 开始事务
	var result *permission.CreateRoleRes
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 创建角色
		roleId, err := dao.Roles.Ctx(ctx).TX(tx).Data(do.Roles{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Status:      req.Status,
		}).InsertAndGetId()

		if err != nil {
			g.Log().Error(ctx, "create role failed:", err)
			return errors.New("创建角色失败")
		}

		// 分配权限
		if len(req.Permissions) > 0 {
			err = s.assignRolePermissions(ctx, tx, roleId, req.Permissions)
			if err != nil {
				return err
			}
		}

		result = &permission.CreateRoleRes{Id: roleId}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateRole 更新角色
func (s *sPermission) UpdateRole(ctx context.Context, req *permission.UpdateRoleReq) (*permission.UpdateRoleRes, error) {
	// 检查角色是否存在
	exists, err := dao.Roles.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check role exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("角色不存在")
	}

	// 检查角色代码是否被其他角色使用
	count, err := dao.Roles.Ctx(ctx).Where("code", req.Code).WhereNot("id", req.Id).Count()
	if err != nil {
		g.Log().Error(ctx, "check role code failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("角色代码已被使用")
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 更新角色
		_, err = dao.Roles.Ctx(ctx).TX(tx).Data(do.Roles{
			Name:        req.Name,
			Code:        req.Code,
			Description: req.Description,
			Status:      req.Status,
			UpdatedAt:   gtime.Now(),
		}).Where("id", req.Id).Update()

		if err != nil {
			g.Log().Error(ctx, "update role failed:", err)
			return errors.New("更新角色失败")
		}

		// 删除原有权限
		_, err = dao.RolePermissions.Ctx(ctx).TX(tx).Where("role_id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete role permissions failed:", err)
			return errors.New("更新角色权限失败")
		}

		// 重新分配权限
		if len(req.Permissions) > 0 {
			err = s.assignRolePermissions(ctx, tx, req.Id, req.Permissions)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &permission.UpdateRoleRes{}, nil
}

// DeleteRole 删除角色
func (s *sPermission) DeleteRole(ctx context.Context, req *permission.DeleteRoleReq) (*permission.DeleteRoleRes, error) {
	// 检查角色是否被用户使用
	count, err := dao.UserRoles.Ctx(ctx).Where("role_id", req.Id).Count()
	if err != nil {
		g.Log().Error(ctx, "check role usage failed:", err)
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("角色正在被用户使用，无法删除")
	}

	// 开始事务删除角色及其权限
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除角色权限关联
		_, err = dao.RolePermissions.Ctx(ctx).TX(tx).Where("role_id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete role permissions failed:", err)
			return errors.New("删除角色权限失败")
		}

		// 删除角色
		_, err = dao.Roles.Ctx(ctx).TX(tx).Where("id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete role failed:", err)
			return errors.New("删除角色失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &permission.DeleteRoleRes{}, nil
}

// GetRole 获取角色详情
func (s *sPermission) GetRole(ctx context.Context, req *permission.GetRoleReq) (*permission.GetRoleRes, error) {
	var role entity.Roles
	err := dao.Roles.Ctx(ctx).Where("id", req.Id).Scan(&role)
	if err != nil {
		g.Log().Error(ctx, "get role failed:", err)
		return nil, err
	}

	if role.Id == 0 {
		return nil, errors.New("角色不存在")
	}

	// 获取角色权限
	permissions, err := s.getRolePermissions(ctx, role.Id)
	if err != nil {
		g.Log().Warning(ctx, "get role permissions failed:", err)
		permissions = []permission.PermissionInfo{}
	}

	// 获取角色用户数量
	userCount, _ := dao.UserRoles.Ctx(ctx).Where("role_id", role.Id).Count()

	return &permission.GetRoleRes{
		RoleInfo: &permission.RoleInfo{
			Id:          role.Id,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Status:      role.Status,
			UserCount:   userCount,
			Permissions: permissions,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		},
	}, nil
}

// AssignRolePermissions 分配角色权限
func (s *sPermission) AssignRolePermissions(ctx context.Context, req *permission.AssignRolePermissionsReq) (*permission.AssignRolePermissionsRes, error) {
	// 检查角色是否存在
	exists, err := dao.Roles.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		g.Log().Error(ctx, "check role exists failed:", err)
		return nil, err
	}
	if exists.IsEmpty() {
		return nil, errors.New("角色不存在")
	}

	// 开始事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除原有权限
		_, err = dao.RolePermissions.Ctx(ctx).TX(tx).Where("role_id", req.Id).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete role permissions failed:", err)
			return errors.New("删除原有权限失败")
		}

		// 分配新权限
		if len(req.Permissions) > 0 {
			err = s.assignRolePermissions(ctx, tx, req.Id, req.Permissions)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &permission.AssignRolePermissionsRes{}, nil
}

// ============================================================================
// 用户角色管理
// ============================================================================

// ListUserRoles 获取用户角色列表
func (s *sPermission) ListUserRoles(ctx context.Context, req *permission.ListUserRolesReq) (*permission.ListUserRolesRes, error) {
	results, err := dao.UserRoles.Ctx(ctx).
		LeftJoin("roles r", "user_roles.role_id = r.id").
		Fields("user_roles.*, r.name as role_name, r.code as role_code").
		Where("user_roles.user_id", req.UserId).
		Where("r.status", 1). // 只显示启用的角色
		OrderDesc("user_roles.assigned_at").
		All()

	if err != nil {
		g.Log().Error(ctx, "get user roles failed:", err)
		return nil, err
	}

	// 转换数据
	list := make([]permission.UserRoleInfo, 0, len(results))
	for _, record := range results {
		list = append(list, permission.UserRoleInfo{
			Id:         record["id"].Int64(),
			UserId:     record["user_id"].Int64(),
			RoleId:     record["role_id"].Int64(),
			RoleName:   record["role_name"].String(),
			RoleCode:   record["role_code"].String(),
			AssignedBy: record["assigned_by"].Int64(),
			ExpiresAt:  record["expires_at"].GTime(),
			CreatedAt:  record["assigned_at"].GTime(), // 使用assigned_at作为CreatedAt
		})
	}

	return &permission.ListUserRolesRes{List: list}, nil
}

// AssignUserRoles 分配用户角色
func (s *sPermission) AssignUserRoles(ctx context.Context, req *permission.AssignUserRolesReq) (*permission.AssignUserRolesRes, error) {
	// 获取当前操作用户ID
	userIdVar := g.RequestFromCtx(ctx).GetCtxVar("user_id")
	if userIdVar == nil {
		return nil, errors.New("未获取到操作用户信息")
	}
	assignedBy := gconv.Int64(userIdVar)

	// 解析过期时间
	var expiresAt *gtime.Time
	if req.ExpiresAt != "" {
		expTime, err := gtime.StrToTime(req.ExpiresAt)
		if err != nil {
			return nil, errors.New("过期时间格式错误")
		}
		expiresAt = expTime
	}

	// 开始事务
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除用户原有角色
		_, err := dao.UserRoles.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
		if err != nil {
			g.Log().Error(ctx, "delete user roles failed:", err)
			return errors.New("删除用户原有角色失败")
		}

		// 分配新角色
		for _, roleId := range req.Roles {
			_, err = dao.UserRoles.Ctx(ctx).TX(tx).Data(do.UserRoles{
				UserId:     req.UserId,
				RoleId:     roleId,
				AssignedBy: assignedBy,
				ExpiresAt:  expiresAt,
			}).Insert()

			if err != nil {
				g.Log().Error(ctx, "assign user role failed:", err)
				return errors.New("分配用户角色失败")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &permission.AssignUserRolesRes{}, nil
}

// RemoveUserRole 移除用户角色
func (s *sPermission) RemoveUserRole(ctx context.Context, req *permission.RemoveUserRoleReq) (*permission.RemoveUserRoleRes, error) {
	_, err := dao.UserRoles.Ctx(ctx).
		Where("user_id", req.UserId).
		Where("role_id", req.RoleId).
		Delete()

	if err != nil {
		g.Log().Error(ctx, "remove user role failed:", err)
		return nil, errors.New("移除用户角色失败")
	}

	return &permission.RemoveUserRoleRes{}, nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// getRolePermissions 获取角色权限列表
func (s *sPermission) getRolePermissions(ctx context.Context, roleId int64) ([]permission.PermissionInfo, error) {
	var permissions []entity.Permissions
	err := dao.Permissions.Ctx(ctx).
		LeftJoin("role_permissions rp", "permissions.id = rp.permission_id").
		Where("rp.role_id", roleId).
		OrderAsc("permissions.resource, permissions.action").
		Scan(&permissions)

	if err != nil {
		return nil, err
	}

	result := make([]permission.PermissionInfo, 0, len(permissions))
	for _, p := range permissions {
		result = append(result, permission.PermissionInfo{
			Id:          p.Id,
			Name:        p.Name,
			Code:        p.Code,
			Resource:    p.Resource,
			Action:      p.Action,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.CreatedAt, // 使用CreatedAt作为UpdatedAt的替代
		})
	}

	return result, nil
}

// assignRolePermissions 分配角色权限（事务内使用）
func (s *sPermission) assignRolePermissions(ctx context.Context, tx gdb.TX, roleId int64, permissionIds []int64) error {
	for _, permissionId := range permissionIds {
		_, err := dao.RolePermissions.Ctx(ctx).TX(tx).Data(do.RolePermissions{
			RoleId:       roleId,
			PermissionId: permissionId,
		}).Insert()

		if err != nil {
			g.Log().Error(ctx, "assign role permission failed:", err)
			return errors.New("分配角色权限失败")
		}
	}

	return nil
}