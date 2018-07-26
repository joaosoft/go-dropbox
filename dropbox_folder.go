package dropbox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	errors "github.com/joaosoft/errors"
	manager "github.com/joaosoft/manager"
)

type Folder struct {
	client manager.IGateway
	config *DropboxConfig
}

type listFolderRequest struct {
	Path                            string `json:"path"`
	Recursive                       bool   `json:"recursive"`
	IncludeMediaInfo                bool   `json:"include_media_info"`
	IncludeDeleted                  bool   `json:"include_deleted"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members"`
	IncludeMountedFolders           bool   `json:"include_mounted_folders"`
}

type listFolderResponse struct {
	Entries []struct {
		Tag            string    `json:".tag"`
		Name           string    `json:"name"`
		ID             string    `json:"id"`
		ClientModified time.Time `json:"client_modified,omitempty"`
		ServerModified time.Time `json:"server_modified,omitempty"`
		Rev            string    `json:"rev,omitempty"`
		Size           int       `json:"size,omitempty"`
		PathLower      string    `json:"path_lower"`
		PathDisplay    string    `json:"path_display"`
		SharingInfo    struct {
			ReadOnly             bool   `json:"read_only"`
			ParentSharedFolderID string `json:"parent_shared_folder_id"`
			ModifiedBy           string `json:"modified_by"`
		} `json:"sharing_info"`
		PropertyGroups []struct {
			TemplateID string `json:"template_id"`
			Fields     []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
		} `json:"property_groups"`
		HasExplicitSharedMembers bool   `json:"has_explicit_shared_members,omitempty"`
		ContentHash              string `json:"content_hash,omitempty"`
	} `json:"entries"`
	Cursor  string `json:"cursor"`
	HasMore bool   `json:"has_more"`
}

func (f *Folder) List(path string) (*listFolderResponse, *errors.ErrorData) {
	if path == "/" {
		path = ""
	}
	body, err := json.Marshal(listFolderRequest{
		Path:                            path,
		Recursive:                       false,
		IncludeMediaInfo:                false,
		IncludeDeleted:                  false,
		IncludeHasExplicitSharedMembers: false,
		IncludeMountedFolders:           true,
	})
	if err != nil {
		newErr := errors.New(err)
		log.Error("error marshal bodyArgs").ToErrorData(newErr)
		return nil, newErr
	}

	headers := manager.Headers{
		"Authorization": {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":  {"application/json"},
	}

	dropboxResponse := &listFolderResponse{}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Api, "/files/list_folder", headers, body); err != nil {
		newErr := errors.New(err)
		log.WithField("response", response).Error("error listing Folder").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, errors.New(err)
	} else if response == nil {
		var err error
		log.Error("error listing Folder").ToError(&err)
		return nil, errors.New(err)
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			newErr := errors.New(err)
			log.Error("error converting Dropbox response data").ToErrorData(newErr)
			return nil, newErr
		}
		return dropboxResponse, nil
	}

	return nil, nil
}

type createFolderRequest struct {
	Path       string `json:"path"`
	AutoRename bool   `json:"autorename"`
}

type createFolderResponse struct {
	Metadata struct {
		Name        string `json:"name"`
		ID          string `json:"id"`
		PathLower   string `json:"path_lower"`
		PathDisplay string `json:"path_display"`
		SharingInfo struct {
			ReadOnly             bool   `json:"read_only"`
			ParentSharedFolderID string `json:"parent_shared_folder_id"`
			TraverseOnly         bool   `json:"traverse_only"`
			NoAccess             bool   `json:"no_access"`
		} `json:"sharing_info"`
		PropertyGroups []struct {
			TemplateID string `json:"template_id"`
			Fields     []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"fields"`
		} `json:"property_groups"`
	} `json:"metadata"`
}

func (f *Folder) Create(path string) (*createFolderResponse, *errors.ErrorData) {
	if path == "/" {
		path = ""
	}
	body, err := json.Marshal(createFolderRequest{
		Path:       path,
		AutoRename: false,
	})
	if err != nil {
		newErr := errors.New(err)
		log.Error("error marshal bodyArgs").ToErrorData(newErr)
		return nil, newErr
	}

	headers := manager.Headers{
		"Authorization": {fmt.Sprintf("%s %s", f.config.Authorization.Access, f.config.Authorization.Token)},
		"Content-Type":  {"application/json"},
	}

	dropboxResponse := &createFolderResponse{}
	if status, response, err := f.client.Request(http.MethodPost, f.config.Hosts.Api, "/files/create_folder_v2", headers, body); err != nil {
		newErr := errors.New(err)
		log.WithField("response", response).Error("error creating Folder").ToErrorData(newErr)
		return nil, newErr
	} else if status != http.StatusOK {
		var err error
		log.WithField("response", response).Errorf("response status %d instead of %d", status, http.StatusOK).ToError(&err)
		return nil, errors.New(err)
	} else if response == nil {
		var err error
		log.Error("error creating Folder").ToError(&err)
		return nil, errors.New(err)
	} else {
		if err := json.Unmarshal(response, dropboxResponse); err != nil {
			newErr := errors.New(err)
			log.Error("error converting Dropbox response data").ToErrorData(newErr)
			return nil, newErr
		}
		return dropboxResponse, nil
	}

	return nil, nil
}

func (f *Folder) DeleteFolder(path string) (*deleteFileResponse, *errors.ErrorData) {
	file := File{
		client: f.client,
		config: f.config,
	}

	return file.Delete(path)
}