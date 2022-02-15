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