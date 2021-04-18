package fscat

import (
	"encoding/json"
	"github.com/joeyscat/file-sync/internal/share"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

type Server struct {
	RootPath          string
	UserList          []UserMetadata
	userMetadataMutex sync.Mutex
}

type UserMetadata struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DataDir  string `json:"-"`
}

func (s *Server) Init() error {
	s.userMetadataMutex = sync.Mutex{}

	http.HandleFunc("/v1/register", s.register)
	http.HandleFunc("/v1/login", s.login)
	http.HandleFunc("/v1/tree", s.tree)
	http.HandleFunc("/v1/put", s.put)
	http.HandleFunc("/v1/get", s.get)
	http.HandleFunc("/v1/info", s.info)
	http.HandleFunc("/v1/log", s.log)
	return nil
}

func (s *Server) Start() error {
	err := http.ListenAndServe(":8000", nil)

	return err
}

// create user and data directory
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user UserMetadata
	err := parse(r, w, &user)
	if err != nil {
		log.Println(err)
		return
	}

	if user.Username == "" || user.Password == "" {
		resp := &share.Response{
			Status: http.StatusBadRequest,
			Code:   -1,
			Msg:    "username or password must not empty",
		}
		resp.WriteJson(w)
		return
	}

	s.userMetadataMutex.Lock()
	defer s.userMetadataMutex.Unlock()
	userList := s.UserList
	for _, u := range userList {
		if u.Username == user.Username {
			resp := &share.Response{
				Status: http.StatusConflict,
				Code:   -1,
				Msg:    "user exists",
			}
			resp.WriteJson(w)
			return
		}
	}

	user.DataDir = genUserDataDir(user.Username)
	err = os.MkdirAll(user.DataDir, 0755)
	if err != nil {
		log.Println(err)
		resp := &share.Response{
			Status: http.StatusInternalServerError,
			Code:   -1,
			Msg:    "register error",
		}
		resp.WriteJson(w)
		return
	}

	userList = append(userList, user)
	err = saveUserMetadata(userList)
	if err != nil {
		log.Println(err)
		resp := &share.Response{
			Status: http.StatusInternalServerError,
			Code:   -1,
			Msg:    "register error",
		}
		resp.WriteJson(w)
		return
	}
	s.UserList = userList

	resp := &share.Response{
		Status: http.StatusOK,
		Code:   0,
		Msg:    "OK",
	}
	resp.WriteJson(w)
	return
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var resp *share.Response
	var user UserMetadata

	err := parse(r, w, &user)
	if err != nil {
		log.Println(err)
		return
	}

	if user.Username == "" || user.Password == "" {
		resp = &share.Response{
			Status: http.StatusBadRequest,
			Code:   -1,
			Msg:    "username or password must not empty",
		}
		resp.WriteJson(w)
		return
	}

	s.userMetadataMutex.Lock()
	defer s.userMetadataMutex.Unlock()
	userList := s.UserList
	for _, u := range userList {
		if u.Username == user.Username {
			if u.Password == user.Password {
				resp = &share.Response{
					Status: http.StatusOK,
					Code:   0,
					Msg:    "OK",
				}
				resp.WriteJson(w)
				return
			} else {
				break
			}
		}
	}
	resp = &share.Response{
		Status: http.StatusUnauthorized,
		Code:   -1,
		Msg:    "login failed",
	}
	resp.WriteJson(w)
	return
}

func (s *Server) tree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var resp *share.Response
	var user UserMetadata

	err := parse(r, w, &user)
	if err != nil {
		log.Println(err)
		return
	}

	userList := s.UserList
	user.DataDir = ""
	for _, u := range userList {
		if user.Username == u.Username {
			user.DataDir = u.DataDir
			break
		}
	}
	if user.DataDir != "" {
		// TODO load file tree and write to resp
		panic("unimplemented")
	} else {
		resp = &share.Response{
			Status: http.StatusUnauthorized,
			Code:   -1,
			Msg:    "user not exists", // FIXME
		}
		resp.WriteJson(w)
	}
}

func parse(r *http.Request, w http.ResponseWriter, i interface{}) error {
	var resp *share.Response
	bs, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		resp = &share.Response{
			Status: http.StatusBadRequest,
			Code:   -1,
			Msg:    "read req body error",
		}
		resp.WriteJson(w)
		return err
	}

	err = json.Unmarshal(bs, i)
	if err != nil {
		resp = &share.Response{
			Status: http.StatusBadRequest,
			Code:   -1,
			Msg:    "parse req body error",
		}
		resp.WriteJson(w)
		return err
	}
	return nil
}

func (s *Server) put(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (s *Server) info(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (s *Server) log(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func genUserDataDir(username string) string {
	// TODO
	return path.Join(InitPath, DataDir, username)
}
