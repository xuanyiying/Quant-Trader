CREATE TABLE sys_tenant
(
    id              VARCHAR(32)  NOT NULL COMMENT '租户ID',
    tenant_code     VARCHAR(50)  NOT NULL COMMENT '租户编码',
    tenant_name     VARCHAR(100) NOT NULL COMMENT '租户名称',
    tenant_status   VARCHAR(20)           DEFAULT NULL COMMENT '租户状态',
    tenant_size     VARCHAR(20)           DEFAULT NULL COMMENT '租户规模',
    created_by      VARCHAR(32)           DEFAULT NULL COMMENT '创建人',
    created_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by      VARCHAR(32)           DEFAULT NULL COMMENT '更新人',
    updated_at      DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_deleted      TINYINT(1)            DEFAULT '0' COMMENT '是否删除：0-否，1-是',
    tenant_desc     TEXT COMMENT '租户描述',
    tenant_logo     VARCHAR(255) COMMENT '租户logo',
    tenant_contact  VARCHAR(255) COMMENT '租户联系方式',
    tenant_address  VARCHAR(255) COMMENT '租户地址',
    tenant_website  VARCHAR(255) COMMENT '租户官网',
    tenant_email    VARCHAR(255) COMMENT '租户邮箱',
    tenant_phone    VARCHAR(20) COMMENT '租户电话',
    tenant_fax      VARCHAR(20) COMMENT '租户传真',
    tenant_postcode VARCHAR(20) COMMENT '租户邮编',
    tenant_country  VARCHAR(20) COMMENT '租户国家',
    tenant_province VARCHAR(20) COMMENT '租户省份',
    tenant_city     VARCHAR(20) COMMENT '租户城市',
    tenant_district VARCHAR(20) COMMENT '租户区县',
    tenant_town     VARCHAR(20) COMMENT '租户乡镇',
    tenant_village  VARCHAR(20) COMMENT '租户村',
    tenant_street   VARCHAR(20) COMMENT '租户街道',
    tenant_no       VARCHAR(20) COMMENT '租户门牌号',
    PRIMARY KEY (id)

) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='系统租户表';
-- 用户表
CREATE TABLE sys_user
(
    id                       VARCHAR(32)  NOT NULL COMMENT '用户ID',
    username                 VARCHAR(50)  NOT NULL COMMENT '用户名',
    password                 VARCHAR(128) NOT NULL COMMENT '密码（加密存储）',
    salt                     VARCHAR(32)  NOT NULL COMMENT '加密盐值',
    real_name                VARCHAR(50)           DEFAULT NULL COMMENT '真实姓名',
    nick_name                VARCHAR(50)           DEFAULT NULL COMMENT '昵称',
    avatar                   VARCHAR(255)          DEFAULT NULL COMMENT '头像URL',
    gender                   VARCHAR(20)           DEFAULT 'unknown' COMMENT '性别：male|female|other|unknown',
    mobile                   VARCHAR(20)           DEFAULT NULL COMMENT '手机号码',
    email                    VARCHAR(100)          DEFAULT NULL COMMENT '电子邮箱',
    id_card                  VARCHAR(20)           DEFAULT NULL COMMENT '身份证号（加密存储）',
    birthday                 DATE                  DEFAULT NULL COMMENT '出生日期',
    user_type                VARCHAR(20)           DEFAULT 'STAFF' COMMENT '用户类型：ADMIN-超管，STAFF-员工，PATIENT-患者',
    last_login_ip            VARCHAR(50)           DEFAULT NULL COMMENT '最后登录IP',
    last_login_time          DATETIME              DEFAULT NULL COMMENT '最后登录时间',
    login_count              INT(11)               DEFAULT '0' COMMENT '登录次数',
    error_count              INT(11)               DEFAULT '0' COMMENT '错误次数',
    error_time               DATETIME              DEFAULT NULL COMMENT '最后错误时间',
    bio                      VARCHAR(255)          DEFAULT NULL COMMENT '个人简介',
    created_by               VARCHAR(32)           DEFAULT NULL COMMENT '创建人',
    created_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by               VARCHAR(32)           DEFAULT NULL COMMENT '更新人',
    updated_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_deleted               TINYINT(1)            DEFAULT '0' COMMENT '是否删除：0-否，1-是',
    tenant_id                VARCHAR(32)  NOT NULL COMMENT '租户ID',
    language_preference      VARCHAR(50) COMMENT '语言偏好',
    communication_preference VARCHAR(50) COMMENT '通信偏好',
    timezone                 VARCHAR(50) COMMENT '时区',
    PRIMARY KEY (id),
    KEY idx_mobile (mobile),
    KEY idx_email (email),
    KEY idx_real_name (real_name)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COLLATE = utf8mb4_unicode_ci COMMENT ='系统用户表';
-- 角色表
CREATE TABLE sys_role
(
    id          VARCHAR(32) NOT NULL COMMENT '角色ID',
    role_name   VARCHAR(50) NOT NULL COMMENT '角色名称',
    role_code   VARCHAR(50) NOT NULL COMMENT '角色编码',
    description VARCHAR(255)         DEFAULT NULL COMMENT '角色描述',
    data_scope  CHAR(1)              DEFAULT '0' COMMENT '数据范围：0-全部，1-本部门，2-本部门及下属部门，3-自定义部门，4-仅本人',
    status      TINYINT(1)  NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
    sort        INT(11)              DEFAULT '0' COMMENT '排序号',
    is_system   TINYINT(1)           DEFAULT '0' COMMENT '是否系统预设：0-否，1-是',
    remark      VARCHAR(255)         DEFAULT NULL COMMENT '备注',
    created_by  VARCHAR(32)          DEFAULT NULL COMMENT '创建人',
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by  VARCHAR(32)          DEFAULT NULL COMMENT '更新人',
    updated_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_deleted  TINYINT(1)           DEFAULT '0' COMMENT '是否删除：0-否，1-是',
    tenant_id   VARCHAR(32) NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_role_code_tenant (role_code, tenant_id),
    KEY idx_role_name (role_name)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='系统角色表';

-- 权限表
CREATE TABLE sys_permission
(
    id                       VARCHAR(32)  NOT NULL COMMENT '权限ID',
    parent_id                VARCHAR(32)           DEFAULT NULL COMMENT '父权限ID',
    permission_name          VARCHAR(50)  NOT NULL COMMENT '权限名称',
    permission_code          VARCHAR(100) NOT NULL COMMENT '权限编码',
    permission_type          CHAR(1)      NOT NULL COMMENT '权限类型：M-菜单，B-按钮，A-API',
    path                     VARCHAR(255)          DEFAULT NULL COMMENT '路由地址',
    component                VARCHAR(255)          DEFAULT NULL COMMENT '前端组件',
    redirect                 VARCHAR(255)          DEFAULT NULL COMMENT '重定向地址',
    icon                     VARCHAR(100)          DEFAULT NULL COMMENT '权限图标',
    is_show                  TINYINT(1)   NOT NULL DEFAULT '1' COMMENT '是否显示：0-隐藏，1-显示',
    is_cache                 TINYINT(1)   NOT NULL DEFAULT '0' COMMENT '是否缓存：0-不缓存，1-缓存',
    is_frame                 TINYINT(1)   NOT NULL DEFAULT '0' COMMENT '是否外链：0-否，1-是',
    api_path                 VARCHAR(255)          DEFAULT NULL COMMENT 'API路径',
    api_method               VARCHAR(10)           DEFAULT NULL COMMENT 'API方法：GET,POST,PUT,DELETE等',
    status                   TINYINT(1)   NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
    sort                     INT(11)               DEFAULT '0' COMMENT '排序号',
    remark                   VARCHAR(255)          DEFAULT NULL COMMENT '备注',
    created_by               VARCHAR(32)           DEFAULT NULL COMMENT '创建人',
    created_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by               VARCHAR(32)           DEFAULT NULL COMMENT '更新人',
    updated_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_deleted               TINYINT(1)            DEFAULT '0' COMMENT '是否删除：0-否，1-是',
    scope_resource           VARCHAR(100) COMMENT '权限范围资源类型',
    action_type              VARCHAR(50) COMMENT '操作类型:C|R|U|D|E',
    data_category            VARCHAR(100) COMMENT '数据类别',
    date_criteria            VARCHAR(255) COMMENT '日期条件',
    tenant_id                VARCHAR(32)  NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_permission_code (permission_code),
    KEY idx_parent_id (parent_id),
    KEY idx_permission_type (permission_type)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='系统权限表';

-- 权限分组表
CREATE TABLE sys_permission_group
(
    id          VARCHAR(32)  NOT NULL COMMENT '分组ID',
    group_name  VARCHAR(100) NOT NULL COMMENT '分组名称',
    group_code  VARCHAR(100) NOT NULL COMMENT '分组编码',
    description VARCHAR(255) COMMENT '分组描述',
    status      TINYINT(1)   NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用',
    created_by  VARCHAR(32) COMMENT '创建人',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by  VARCHAR(32) COMMENT '更新人',
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    tenant_id   VARCHAR(32)  NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_code (group_code, tenant_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='权限分组表';

-- 权限分组关联表
CREATE TABLE sys_permission_group_relation
(
    id            VARCHAR(32) NOT NULL COMMENT '关联ID',
    group_id      VARCHAR(32) NOT NULL COMMENT '分组ID',
    permission_id VARCHAR(32) NOT NULL COMMENT '权限ID',
    created_by    VARCHAR(32) COMMENT '创建人',
    created_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    tenant_id     VARCHAR(32) NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_group_permission (group_id, permission_id, tenant_id),
    KEY idx_permission (permission_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='权限分组关联表';

-- 用户-角色关联表
CREATE TABLE sys_user_role
(
    id         VARCHAR(32) NOT NULL COMMENT '关联ID',
    user_id    VARCHAR(32) NOT NULL COMMENT '用户ID',
    role_id    VARCHAR(32) NOT NULL COMMENT '角色ID',
    created_by VARCHAR(32)          DEFAULT NULL COMMENT '创建人',
    created_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    tenant_id  VARCHAR(32) NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_role (user_id, role_id, tenant_id),
    KEY idx_role_id (role_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='用户角色关联表';

-- 角色-权限关联表
CREATE TABLE sys_role_permission
(
    id            VARCHAR(32) NOT NULL COMMENT '关联ID',
    role_id       VARCHAR(32) NOT NULL COMMENT '角色ID',
    permission_id VARCHAR(32) NOT NULL COMMENT '权限ID',
    created_by    VARCHAR(32)          DEFAULT NULL COMMENT '创建人',
    created_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    tenant_id     VARCHAR(32) NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_role_permission (role_id, permission_id, tenant_id),
    KEY idx_permission_id (permission_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='角色权限关联表';

-- 用户令牌表
CREATE TABLE sys_user_token
(
    id          VARCHAR(32)  NOT NULL COMMENT '令牌ID',
    user_id     VARCHAR(32)  NOT NULL COMMENT '用户ID',
    token       VARCHAR(255) NOT NULL COMMENT '令牌内容',
    token_type  VARCHAR(20)           DEFAULT 'Bearer' COMMENT '令牌类型',
    client_id   VARCHAR(100)          DEFAULT NULL COMMENT '客户端ID',
    ip_address  VARCHAR(50)           DEFAULT NULL COMMENT 'IP地址',
    device_info VARCHAR(255)          DEFAULT NULL COMMENT '设备信息',
    expires_at  DATETIME     NOT NULL COMMENT '过期时间',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_revoked  TINYINT(1)            DEFAULT '0' COMMENT '是否已撤销：0-否，1-是',
    tenant_id   VARCHAR(32)  NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_token (token),
    KEY idx_user_id (user_id),
    KEY idx_expires_at (expires_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='用户令牌表';


-- 医院表
CREATE TABLE org_hospital
(
    id                 VARCHAR(32)  NOT NULL COMMENT '医院ID',
    hospital_code      VARCHAR(50)  NOT NULL COMMENT '医院代码',
    hospital_name      VARCHAR(100) NOT NULL COMMENT '医院名称',
    short_name         VARCHAR(50)           DEFAULT NULL COMMENT '医院简称',
    hospital_type      VARCHAR(20)           DEFAULT NULL COMMENT '医院类型：三甲、二甲等',
    address            VARCHAR(255)          DEFAULT NULL COMMENT '医院地址',
    zip_code           VARCHAR(10)           DEFAULT NULL COMMENT '邮政编码',
    phone              VARCHAR(20)           DEFAULT NULL COMMENT '联系电话',
    email              VARCHAR(100)          DEFAULT NULL COMMENT '电子邮箱',
    website            VARCHAR(255)          DEFAULT NULL COMMENT '官网地址',
    legal_person       VARCHAR(50)           DEFAULT NULL COMMENT '法人代表',
    latitude           DECIMAL(10, 7)        DEFAULT NULL COMMENT '纬度',
    longitude          DECIMAL(10, 7)        DEFAULT NULL COMMENT '经度',
    province           VARCHAR(50)           DEFAULT NULL COMMENT '省份',
    city               VARCHAR(50)           DEFAULT NULL COMMENT '城市',
    district           VARCHAR(50)           DEFAULT NULL COMMENT '区县',
    level              VARCHAR(20)           DEFAULT NULL COMMENT '医院等级',
    intro              TEXT                  DEFAULT NULL COMMENT '医院简介',
    logo               VARCHAR(255)          DEFAULT NULL COMMENT '医院LOGO',
    license_code       VARCHAR(100)          DEFAULT NULL COMMENT '医疗机构执业许可证号',
    establishment_date DATE                  DEFAULT NULL COMMENT '成立日期',
    bed_count          INT(11)               DEFAULT NULL COMMENT '编制床位数',
    status             TINYINT(1)   NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
    remark             VARCHAR(255)          DEFAULT NULL COMMENT '备注',
    created_by         VARCHAR(32)           DEFAULT NULL COMMENT '创建人',
    created_at         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by         VARCHAR(32)           DEFAULT NULL COMMENT '更新人',
    updated_at         DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_deleted         TINYINT(1)            DEFAULT '0' COMMENT '是否删除：0-否，1-是',
    tenant_id          VARCHAR(32)  NOT NULL COMMENT '租户ID',
    qualification      JSON COMMENT '机构资质信息',
    contact            JSON COMMENT '联系人信息',
    alias              TEXT COMMENT '机构别名',
    partOf             VARCHAR(32) COMMENT '上级机构',
    endpoint           JSON COMMENT '服务端点',
    PRIMARY KEY (id),
    UNIQUE KEY uk_hospital_code_tenant (hospital_code, tenant_id),
    KEY idx_hospital_name (hospital_name),
    KEY idx_province_city (province, city)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='医院表';

-- 科室表
CREATE TABLE org_department
(
    id          VARCHAR(32)  NOT NULL COMMENT '科室ID',
    hospital_id VARCHAR(32)  NOT NULL COMMENT '所属医院ID',
    area_id     VARCHAR(32)           DEFAULT NULL COMMENT '所属院区ID',
    parent_id   VARCHAR(32)           DEFAULT NULL COMMENT '父科室ID',
    dept_code   VARCHAR(50)  NOT NULL COMMENT '科室代码',
    dept_name   VARCHAR(100) NOT NULL COMMENT '科室名称',
    short_name  VARCHAR(50)           DEFAULT NULL COMMENT '科室简称',
    dept_type   VARCHAR(20)           DEFAULT NULL COMMENT '科室类型：门诊、住院、医技、行政等',
    dept_level  VARCHAR(20)           DEFAULT NULL COMMENT '科室级别：一级、二级等',
    sort        INT(11)               DEFAULT '0' COMMENT '排序号',
    phone       VARCHAR(20)           DEFAULT NULL COMMENT '联系电话',
    location    VARCHAR(255)          DEFAULT NULL COMMENT '位置信息',
    intro       TEXT                  DEFAULT NULL COMMENT '科室介绍',
    director_id VARCHAR(32)           DEFAULT NULL COMMENT '科室主任ID',
    status      TINYINT(1)   NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
    remark      VARCHAR(255)          DEFAULT NULL COMMENT '备注',
    created_by  VARCHAR(32)           DEFAULT NULL COMMENT '创建人',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by  VARCHAR(32)           DEFAULT NULL COMMENT '更新人',
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_deleted  TINYINT(1)            DEFAULT '0' COMMENT '是否删除：0-否，1-是',
    tenant_id   VARCHAR(32)  NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_dept_code_hospital (dept_code, hospital_id),
    KEY idx_hospital_id (hospital_id),
    KEY idx_area_id (area_id),
    KEY idx_parent_id (parent_id),
    KEY idx_dept_name (dept_name)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='科室表';

-- 医护人员表
CREATE TABLE org_staff
(
    id                VARCHAR(32) NOT NULL COMMENT '医护人员ID',
    user_id           VARCHAR(32) NOT NULL COMMENT '关联用户ID',
    staff_no          VARCHAR(50) NOT NULL COMMENT '工号',
    staff_name        VARCHAR(50) NOT NULL COMMENT '姓名',
    gender            CHAR(1)              DEFAULT NULL COMMENT '性别：M-男，F-女，U-未知',
    birth_date        DATE                 DEFAULT NULL COMMENT '出生日期',
    id_card           VARCHAR(20)          DEFAULT NULL COMMENT '身份证号（加密存储）',
    phone             VARCHAR(20)          DEFAULT NULL COMMENT '联系电话',
    email             VARCHAR(100)         DEFAULT NULL COMMENT '电子邮箱',
    address           VARCHAR(255)         DEFAULT NULL COMMENT '联系地址',
    staff_type        VARCHAR(20)          DEFAULT NULL COMMENT '人员类型：医生，护士，医技，行政等',
    education         VARCHAR(50)          DEFAULT NULL COMMENT '学历',
    degree            VARCHAR(50)          DEFAULT NULL COMMENT '学位',
    graduate_school   VARCHAR(100)         DEFAULT NULL COMMENT '毕业院校',
    major             VARCHAR(100)         DEFAULT NULL COMMENT '专业',
    title             VARCHAR(50)          DEFAULT NULL COMMENT '职称',
    title_level       VARCHAR(20)          DEFAULT NULL COMMENT '职称级别：初级、中级、副高、正高',
    qualification_no  VARCHAR(100)         DEFAULT NULL COMMENT '资格证书编号',
    practicing_no     VARCHAR(100)         DEFAULT NULL COMMENT '执业证书编号',
    specialty         VARCHAR(255)         DEFAULT NULL COMMENT '专长',
    join_date         DATE                 DEFAULT NULL COMMENT '入职日期',
    leave_date        DATE                 DEFAULT NULL COMMENT '离职日期',
    employment_status VARCHAR(20)          DEFAULT 'ACTIVE' COMMENT '在职状态：ACTIVE-在职，LEAVE-离职，RETIRE-退休',
    photo             VARCHAR(255)         DEFAULT NULL COMMENT '照片URL',
    intro             TEXT                 DEFAULT NULL COMMENT '个人简介',
    status            TINYINT(1)  NOT NULL DEFAULT '1' COMMENT '状态：0-禁用，1-启用',
    sort              INT(11)              DEFAULT '0' COMMENT '排序号',
    remark            VARCHAR(255)         DEFAULT NULL COMMENT '备注',
    created_by        VARCHAR(32)          DEFAULT NULL COMMENT '创建人',
    created_at        DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_by        VARCHAR(32)          DEFAULT NULL COMMENT '更新人',
    updated_at        DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_deleted        TINYINT(1)           DEFAULT '0' COMMENT '是否删除：0-否，1-是',
    tenant_id         VARCHAR(32) NOT NULL COMMENT '租户ID',
    communication     JSON COMMENT '通信能力',
    fhir_identifier   JSON COMMENT 'FHIR标识符集合',
    availability      JSON COMMENT '可用性信息',
    specialty_detail  JSON COMMENT '专业详细信息',
    care_team_id      VARCHAR(32) COMMENT '医疗团队ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_id_tenant (user_id, tenant_id),
    UNIQUE KEY uk_staff_no_tenant (staff_no, tenant_id),
    KEY idx_staff_name (staff_name),
    KEY idx_staff_type (staff_type),
    KEY idx_title (title),
    KEY idx_employment_status (employment_status)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='医护人员表';

-- 医疗服务点表
CREATE TABLE org_location
(
    id                       VARCHAR(32)  NOT NULL COMMENT '位置ID',
    status                   VARCHAR(20)  NOT NULL DEFAULT 'active' COMMENT '状态:active|suspended|inactive',
    name                     VARCHAR(100) NOT NULL COMMENT '位置名称',
    description              TEXT COMMENT '位置描述',
    mode                     VARCHAR(20)           DEFAULT 'instance' COMMENT '模式:instance|kind',
    type_code                VARCHAR(50) COMMENT '类型编码',
    type_display             VARCHAR(100) COMMENT '类型显示名称',
    telecom                  VARCHAR(50) COMMENT '联系方式',
    address                  TEXT COMMENT '地址',
    physical_type            VARCHAR(50) COMMENT '物理类型',
    position_longitude       DECIMAL(10, 7) COMMENT '经度',
    position_latitude        DECIMAL(10, 7) COMMENT '纬度',
    position_altitude        DECIMAL(10, 7) COMMENT '海拔',
    managing_organization_id VARCHAR(32) COMMENT '管理组织ID',
    part_of                  VARCHAR(32) COMMENT '所属位置ID',
    operational_status       VARCHAR(50) COMMENT '运营状态',
    created_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    tenant_id                VARCHAR(32)  NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    KEY idx_name (name),
    KEY idx_type (type_code),
    KEY idx_status (status),
    KEY idx_organization (managing_organization_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='医疗服务点表';

-- 医疗团队表
CREATE TABLE org_care_team
(
    id                       VARCHAR(32)  NOT NULL COMMENT '团队ID',
    name                     VARCHAR(100) NOT NULL COMMENT '团队名称',
    status                   VARCHAR(20)  NOT NULL DEFAULT 'active' COMMENT '状态:proposed|active|suspended|inactive|entered-in-error',
    category_code            VARCHAR(50) COMMENT '团队类别编码',
    category_display         VARCHAR(100) COMMENT '团队类别显示名称',
    subject_id               VARCHAR(32) COMMENT '服务对象ID(患者)',
    encounter_id             VARCHAR(32) COMMENT '就诊ID',
    reason_code              VARCHAR(50) COMMENT '原因编码',
    reason_display           VARCHAR(100) COMMENT '原因显示名称',
    managing_organization_id VARCHAR(32) COMMENT '管理组织ID',
    period_start             DATETIME COMMENT '开始时间',
    period_end               DATETIME COMMENT '结束时间',
    description              TEXT COMMENT '描述',
    created_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at               DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    tenant_id                VARCHAR(32)  NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    KEY idx_name (name),
    KEY idx_status (status),
    KEY idx_category (category_code),
    KEY idx_subject (subject_id),
    KEY idx_encounter (encounter_id),
    KEY idx_organization (managing_organization_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='医疗团队表';

-- 医疗团队成员表
CREATE TABLE org_care_team_participant
(
    id           VARCHAR(32) NOT NULL COMMENT '成员ID',
    care_team_id VARCHAR(32) NOT NULL COMMENT '团队ID',
    member_id    VARCHAR(32) NOT NULL COMMENT '成员ID(可能是医护人员ID或其他资源ID)',
    member_type  VARCHAR(50) NOT NULL COMMENT '成员类型:practitioner|patient|related-person|organization',
    role_code    VARCHAR(50) COMMENT '角色编码',
    role_display VARCHAR(100) COMMENT '角色显示名称',
    period_start DATETIME COMMENT '开始时间',
    period_end   DATETIME COMMENT '结束时间',
    is_leader    BOOLEAN              DEFAULT FALSE COMMENT '是否为团队负责人',
    status       VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态:active|inactive',
    created_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    tenant_id    VARCHAR(32) NOT NULL COMMENT '租户ID',
    PRIMARY KEY (id),
    UNIQUE KEY uk_team_member (care_team_id, member_id, member_type),
    KEY idx_care_team (care_team_id),
    KEY idx_member (member_id),
    KEY idx_role (role_code)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT ='医疗团队成员表';