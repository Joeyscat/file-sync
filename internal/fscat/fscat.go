package fscat

import (
	"encoding/json"
	"fmt"
	"github.com/joeyscat/file-sync/internal/share"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"runtime"
)

var (
	InitPath = ".local/share/file-sync"
)

const (
	UserMetadataDir      = "users"
	UserMetadataFilename = "users.json"
	DataDir              = "data"
)

func Exec() error {
	if runtime.GOOS == "windows" {
		InitPath = ".local/share/file-sync"
	} else if runtime.GOOS == "darwin" ||
		runtime.GOOS == "linux" {
	}

	// TODO cross compile support https://www.w3cschool.cn/cuhkj/cuhkj-azch266d.html
	u, err := user.Current()
	if err != nil {
		return err
	}

	InitPath = path.Join(u.HomeDir, InitPath)
	exists := share.PathExists(InitPath)

	var userList []UserMetadata
	userMetadataPath := path.Join(InitPath, UserMetadataDir, UserMetadataFilename)
	if exists {
		if share.PathExists(userMetadataPath) {
			userList, err = loadUsers(userMetadataPath)
			if err != nil {
				return fmt.Errorf("read user metadata error: %v", err)
			}
		}
	} else {
		err := os.MkdirAll(path.Join(InitPath, UserMetadataDir), 0755)
		if err != nil {
			return fmt.Errorf("create dir for user metadata error: %v", err)
		}
		f, err := os.Create(userMetadataPath)
		if err != nil {
			return fmt.Errorf("create file for user metadata error: %v", err)
		}
		bs, err := json.Marshal(&userList)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(bs)
		return err
	}

	server := Server{
		RootPath: InitPath,
		UserList: userList,
	}
	err = server.Init()
	if err != nil {
		return err
	}

	return server.Start()
}

func loadUsers(userMetadataPath string) ([]UserMetadata, error) {
	usersJson := userMetadataPath
	dataBytes, err := ioutil.ReadFile(usersJson)
	if dataBytes == nil {
		return nil, err
	}

	var userList []UserMetadata
	err = json.Unmarshal(dataBytes, &userList)
	if err != nil {
		return nil, err
	}

	return userList, nil
}

func saveUserMetadata(userList []UserMetadata) error {
	userMetadataPath := path.Join(InitPath, UserMetadataDir, UserMetadataFilename)
	bs, err := json.Marshal(userList)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(userMetadataPath, bs, 0755)
}
