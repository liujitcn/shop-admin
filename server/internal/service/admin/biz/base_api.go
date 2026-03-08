package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/liujitcn/go-sdk"
	queueData "github.com/liujitcn/go-sdk/queue/data"
	"github.com/liujitcn/shop-admin/server/api/gen/go/admin"
	"github.com/liujitcn/shop-admin/server/cmd/server/assets"
	_const "github.com/liujitcn/shop-admin/server/internal/const"
	"github.com/liujitcn/shop-admin/server/internal/data"
	"github.com/liujitcn/shop-gorm-gen/models"
	"gopkg.in/yaml.v3"
)

type BaseApiCase struct {
	data.BaseApiRepo
}

// NewBaseApiCase new a BaseApi use case.
func NewBaseApiCase(baseApiRepo data.BaseApiRepo) *BaseApiCase {
	return &BaseApiCase{
		BaseApiRepo: baseApiRepo,
	}
}

func (c *BaseApiCase) List(ctx context.Context) (*admin.ListBaseApiResponse, error) {
	condition := &data.BaseApiCondition{}
	all, err := c.FindAll(ctx, condition)
	if err != nil {
		return nil, err
	}
	list := make([]*admin.BaseApi, 0)
	for _, item := range all {
		list = append(list, &admin.BaseApi{
			Id:          item.ID,
			ServiceName: item.ServiceName,
			ServiceDesc: item.ServiceDesc,
			Desc:        item.Desc,
			Operation:   item.Operation,
			Method:      item.Method,
			Path:        item.Path,
		})
	}
	return &admin.ListBaseApiResponse{
		List: list,
	}, nil
}

func (c *BaseApiCase) BatchCreate(ctx context.Context, apis []*models.BaseAPI) error {
	// 查询旧数据
	oldApiList, err := c.FindAll(ctx, &data.BaseApiCondition{})
	if err != nil {
		return err
	}
	// 和id map
	oldApiIdMap := make(map[string]int64)
	for _, oldApi := range oldApiList {
		oldApiIdMap[fmt.Sprintf("%s_%s", oldApi.Method, oldApi.Path)] = oldApi.ID
	}
	apiList := make([]*models.BaseAPI, 0)
	for _, item := range apis {
		key := fmt.Sprintf("%s_%s", item.Method, item.Path)
		if id, ok := oldApiIdMap[key]; ok {
			item.ID = id
			err = c.UpdateByID(ctx, item)
			if err != nil {
				return err
			}
			// 删除map
			delete(oldApiIdMap, key)
		} else {
			apiList = append(apiList, item)
		}
	}
	if len(oldApiIdMap) > 0 {
		oldApiId := make([]int64, 0)
		for _, v := range oldApiIdMap {
			oldApiId = append(oldApiId, v)
		}
		err = c.Delete(ctx, oldApiId)
		if err != nil {
			return err
		}
	}
	if len(apiList) >= 0 {
		return c.BaseApiRepo.BatchCreate(ctx, apiList)
	}
	return nil
}

func (c *BaseApiCase) SaveApi(message queueData.Message) error {
	rb, err := json.Marshal(message.Values)
	if err != nil {
		log.Errorf("json Marshal error, %s", err.Error())
		return err
	}
	var m map[string][]*models.BaseAPI
	err = json.Unmarshal(rb, &m)
	if err != nil {
		log.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	if v, ok := m["data"]; ok {
		err = c.BatchCreate(context.TODO(), v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *BaseApiCase) ApiCheck() error {
	log.Info("Base-API Server ApiCheck...")
	baseApiList, err := c.OpenApiDataToBaseAPI()
	if err != nil {
		return err
	}
	q := sdk.Runtime.GetQueue()
	if q != nil {
		m := make(map[string]interface{})
		m["data"] = baseApiList
		var message queueData.Message
		message, err = sdk.Runtime.GetStreamMessage(_const.ApiCheck, m)
		if err != nil {
			log.Error("GetStreamMessage error, %s ", err.Error())
		} else {
			err = q.Append(_const.ApiCheck, message)
			if err != nil {
				log.Error("Append message error, %s ", err.Error())
			}
		}
	}
	return nil
}

// OpenApiDataToBaseAPI OpenApi 数据转model
func (c *BaseApiCase) OpenApiDataToBaseAPI() ([]*models.BaseAPI, error) {
	// 解析YAML内容到结构体
	var api OpenAPI
	err := yaml.Unmarshal(assets.OpenApiData, &api)
	if err != nil {
		return nil, err
	}
	tagsMap := make(map[string]string)
	for _, item := range api.Tags {
		if strings.HasPrefix(item.Description, "Admin") {
			tagsMap[fmt.Sprintf("admin.%s", item.Name)] = item.Description
		} else if strings.HasPrefix(item.Description, "App") {
			tagsMap[fmt.Sprintf("app.%s", item.Name)] = item.Description
		} else {
			tagsMap[item.Name] = item.Description
		}
	}

	baseApiList := make([]*models.BaseAPI, 0)
	for path, item := range api.Paths {
		getApi := parseOperation(path, "GET", item.Get, tagsMap)
		if getApi != nil {
			baseApiList = append(baseApiList, getApi)
		}
		postApi := parseOperation(path, "POST", item.Post, tagsMap)
		if postApi != nil {
			baseApiList = append(baseApiList, postApi)
		}
		putApi := parseOperation(path, "PUT", item.Put, tagsMap)
		if putApi != nil {
			baseApiList = append(baseApiList, putApi)
		}
		deleteApi := parseOperation(path, "DELETE", item.Delete, tagsMap)
		if deleteApi != nil {
			baseApiList = append(baseApiList, deleteApi)
		}
	}
	//排序
	sort.Slice(baseApiList, func(i, j int) bool {
		return baseApiList[i].Operation < baseApiList[j].Operation
	})
	return baseApiList, nil
}

func parseOperation(path, method string, op *Operation, tagsMap map[string]string) *models.BaseAPI {
	if op == nil {
		return nil
	}
	var pkgName string
	paths := strings.Split(path, "/")
	if len(paths) > 2 {
		pkgName = paths[2]
	}
	var serviceName, serviceDesc string

	tags := op.Tags
	if len(tags) > 0 {
		switch pkgName {
		case "admin":
			serviceName = fmt.Sprintf("admin.%s", tags[0])
		case "app":
			serviceName = fmt.Sprintf("app.%s", tags[0])
		default:
			serviceName = tags[0]

		}
		if v, ok := tagsMap[serviceName]; ok {
			serviceDesc = v
		}
	}
	return &models.BaseAPI{
		ServiceName: serviceName,
		ServiceDesc: serviceDesc,
		Desc:        op.Description,
		Operation:   fmt.Sprintf("/%s.%s", pkgName, strings.ReplaceAll(op.OperationID, "_", "/")),
		Method:      method,
		Path:        path,
	}
}

type OpenAPI struct {
	Paths map[string]PathItem `yaml:"paths"`
	Tags  []TagsItem          `yaml:"tags"`
}

type PathItem struct {
	Get    *Operation `yaml:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
}

type TagsItem struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type Operation struct {
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
	OperationID string   `yaml:"operationId"`
}
