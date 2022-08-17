package enums

// COMPANY_UPDATE_OPTION company update options
type COMPANY_UPDATE_OPTION string

const (
	// APPEND_APPLICATION company update option to append application
	APPEND_APPLICATION = COMPANY_UPDATE_OPTION("APPEND_APPLICATION")
	// APPEND_REPOSITORY company update option to append repository
	APPEND_REPOSITORY = COMPANY_UPDATE_OPTION("APPEND_REPOSITORY")
	// SOFT_DELETE_APPLICATION company update option to soft delete application
	SOFT_DELETE_APPLICATION = COMPANY_UPDATE_OPTION("SOFT_DELETE_APPLICATION")
	// DELETE_APPLICATION company update option to delete application
	DELETE_APPLICATION = COMPANY_UPDATE_OPTION("DELETE_APPLICATION")
	// SOFT_DELETE_REPOSITORY company update option to soft delete repository
	SOFT_DELETE_REPOSITORY = COMPANY_UPDATE_OPTION("SOFT_DELETE_REPOSITORY")
	// DELETE_REPOSITORY company update option to delete repository
	DELETE_REPOSITORY = COMPANY_UPDATE_OPTION("DELETE_REPOSITORY")
)

// COMPANY_STATUS company status options
type COMPANY_STATUS string

const (
	// ACTIVE company status for active company
	ACTIVE = COMPANY_STATUS("ACTIVE")
	// INACTIVE company status for inactive company
	INACTIVE = COMPANY_STATUS("INACTIVE")
)

// REPOSITORY_TYPE repository types[may be any git]
type REPOSITORY_TYPE string

const (
	// GITHUB github as repository
	GITHUB = REPOSITORY_TYPE("GITHUB")
	// BIT_BUCKET bitbucket as repository
	BIT_BUCKET = REPOSITORY_TYPE("BIT_BUCKET")
)

// COMPANY_FLAG company flag types to apply
type COMPANY_FLAG string

const (
	// CREATE_COMPANY create company flag
	CREATE_COMPANY = COMPANY_FLAG("CREATE_COMPANY")
	// UPDATE_REPOSITORIES update company flag
	UPDATE_REPOSITORIES = COMPANY_FLAG("UPDATE_REPOSITORIES")
	// UPDATE_APPLICATIONS update applications flag
	UPDATE_APPLICATIONS = COMPANY_FLAG("UPDATE_APPLICATIONS")
	// GET_COMPANY_BY_ID get company by id flag
	GET_COMPANY_BY_ID = COMPANY_FLAG("GET_COMPANY_BY_ID")
	// GET_COMPANIES get companies flag
	GET_COMPANIES = COMPANY_FLAG("GET_COMPANIES")
	// GET_REPOSITORIES get repositories flag
	GET_REPOSITORIES = COMPANY_FLAG("GET_REPOSITORIES")
)

// REPOSITORY_FLAG repository flag types to apply
type REPOSITORY_FLAG string

const (
	// GET_REPOSITORY get repository flag
	GET_REPOSITORY = REPOSITORY_FLAG("GET_REPOSITORY")
	// GET_APPLICATIONS get applications flag
	GET_APPLICATIONS = REPOSITORY_FLAG("GET_APPLICATIONS")
	// GET_All_APPLICATIONS get all applications flag
	GET_All_APPLICATIONS = REPOSITORY_FLAG("GET_All_APPLICATIONS")
)

// APPLICATION_FLAG application flag types to apply
type APPLICATION_FLAG string

const (
	// GET_APPLICATION get application flag
	GET_APPLICATION = APPLICATION_FLAG("GET_APPLICATION")
)

// USER_FLAG company flag types to apply
type USER_FLAG string

const (
	// CREATE_USER create user flag
	CREATE_USER = USER_FLAG("CREATE_USER")
	// CREATE_ADMIN create admin flag
	CREATE_ADMIN = USER_FLAG("CREATE_ADMIN")
)

// USER_UPDATE_ACTION user update action flag types to apply
type USER_UPDATE_ACTION string

const (
	// RESET_PASSWORD refers to password reset action
	RESET_PASSWORD = USER_UPDATE_ACTION("reset_password")
	// FORGOT_PASSWORD refers to password forgot action
	FORGOT_PASSWORD = USER_UPDATE_ACTION("forgot_password")
	// ATTACH_COMPANY refers to company attachment action
	ATTACH_COMPANY = USER_UPDATE_ACTION("attach_company")
	// UPDATE_STATUS refers to status update action
	UPDATE_STATUS = USER_UPDATE_ACTION("update_status")
)

// AGENT_FLAG agent query action flag types to apply
type AGENT_FLAG string

const (
	// GET_K8SOBJS refers to get k8sobjs action
	GET_K8SOBJS = AGENT_FLAG("get_k8sobjs")
)

// PROCESS_FLAG process query action flag types to apply
type PROCESS_FLAG string

const (
	// GET_PROCESS refers to get process action
	GET_PROCESS = PROCESS_FLAG("get_process")
	// GET_LOGS refers to get logs action
	GET_LOGS = PROCESS_FLAG("get_logs")
)
