package api

import (
	"fmt"
	"github.com/Bedrock-Technology/Dsn/app"
	"github.com/Bedrock-Technology/Dsn/app/config"
	"github.com/Bedrock-Technology/Dsn/app/db/postgres"
	"github.com/Bedrock-Technology/Dsn/app/util"
	"github.com/Bedrock-Technology/Dsn/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type DsnApi struct {
}

type LoadRequest struct {
	TomlPath string `json:"toml_path" binding:"required"`
}

type SqlRequest struct {
	FuncName string `json:"func_name" binding:"required"`
	Params   string `json:"params"`
}

// ExecCmd godoc
// @Summary Execute a sql command
// @Description Execute a sql command with the given parameters
// @Tags Dsn Hun Client
// @Accept  json
// @Produce  json
// @Param func_name path string true "Function Name"
// @Param params path object true "Function Parameters (as a map)"
// @Success 200 {object} []map[string]interface{} "Successful operation"
// @Failure 400 {object} proto.ResponseMsg "Failed operation"
// @Router /dsn/exec/{func_name}/{params} [get]
func (api *DsnApi) ExecCmd(c *gin.Context) {
	funcName := c.Param(config.ExecFuncPath)
	if funcName == "" {
		util.ErrorMsg(c, "Function name is required")
		return
	}

	paramsJson := c.Param(config.ExecParamsPath)
	log.Debugf("ExecCmd", "funcName:", funcName, "params:", paramsJson)
	cs := app.GetConfigStore()
	if paramsJson != "" && !cs.IsSafeInput(paramsJson) {
		util.ErrorMsg(c, "Invalid input parameters with sql injection risk")
		return
	}
	content, err := cs.GetDataByKey(funcName)
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	contentMap := content.(map[string]interface{})
	dsn := contentMap[config.SqlDsn].(string)
	cmd := contentMap[config.SqlCmd].(string)
	query, err := cs.BindParams(cmd, paramsJson)
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	psqlProvider := &postgres.Provider{}
	log.Debugf("ExecCmd", "dsn:", dsn, "query:", query)
	queryRes, err := psqlProvider.ExecCmd(dsn, query)
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	util.SuccessMsg(c, http.StatusOK, fmt.Sprintf("Execute sql %s command successfully", funcName), queryRes)
}

// GetFuncList godoc
// @Summary Get func list for dsn hub client
// @Description Retrieve the list of functions that can be executed
// @Tags Dsn Hun Client
// @Produce  json
// @Success 200 {object} []map[string]string "Successful operation"
// @Failure 400 {object} proto.ResponseMsg "Failed operation"
// @Router /dsn/func_list [get]
func (api *DsnApi) GetFuncList(c *gin.Context) {
	cs := app.GetConfigStore()
	funcList, err := cs.GetFucList()
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	util.SuccessMsg(c, http.StatusOK, "Get function list successfully", funcList)
}

func (api *DsnApi) LoadSqlFile(c *gin.Context) {
	var req LoadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	tomlPath := req.TomlPath
	if _, err := os.Stat(tomlPath); os.IsNotExist(err) {
		util.ErrorMsg(c, "File not found")
		return
	}
	cs := app.GetConfigStore()
	err := cs.LoadFile(tomlPath)
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	util.SuccessMsg(c, http.StatusOK, "Load Toml file successfully", nil)
}

func (api *DsnApi) ExecSql(c *gin.Context) {
	var req SqlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	if req.FuncName == "" {
		util.ErrorMsg(c, "Function name is required")
		return
	}
	log.Debugf("ExecCmd", "funcName:", req.FuncName, "params:", req.Params)
	cs := app.GetConfigStore()
	if req.Params != "" && !cs.IsSafeInput(req.Params) {
		util.ErrorMsg(c, "Invalid input parameters with sql injection risk")
		return
	}
	content, err := cs.GetDataByKey(req.FuncName)
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	contentMap := content.(map[string]interface{})
	dsn := contentMap[config.SqlDsn].(string)
	cmd := contentMap[config.SqlCmd].(string)
	query, err := cs.BindParams(cmd, req.Params)
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	psqlProvider := &postgres.Provider{}
	log.Debugf("ExecCmd", "dsn:", dsn, "query:", query)
	queryRes, err := psqlProvider.ExecCmd(dsn, query)
	if err != nil {
		util.ErrorMsg(c, err.Error())
		return
	}
	util.SuccessMsg(c, http.StatusOK, fmt.Sprintf("Execute sql %s command successfully", req.FuncName), queryRes)
}
