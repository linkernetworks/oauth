package osinstorage

import (
	"sync"

	"github.com/linkernetworks/logger"

	"github.com/RangelReale/osin"
)

type MemoryStorage struct {
	clients        map[string]osin.DefaultClient
	clientsMutex   sync.Mutex
	authorize      map[string]osin.AuthorizeData
	authorizeMutex sync.Mutex
	access         map[string]osin.AccessData
	accessMutex    sync.Mutex
	refresh        map[string]string
	refreshMutex   sync.Mutex
}

func NewMemoryStorage(clients ...osin.DefaultClient) *MemoryStorage {
	r := &MemoryStorage{
		clients:   make(map[string]osin.DefaultClient),
		authorize: make(map[string]osin.AuthorizeData),
		access:    make(map[string]osin.AccessData),
		refresh:   make(map[string]string),
	}

	for _, c := range clients {
		r.clients[c.Id] = c
	}

	return r
}

func (s *MemoryStorage) Clone() osin.Storage {
	return s
}

func (s *MemoryStorage) Close() {
}

func (s *MemoryStorage) GetClient(id string) (osin.Client, error) {
	logger.Debugf("GetClient: %s\n", id)

	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if c, ok := s.clients[id]; ok {
		return &osin.DefaultClient{
			Id:          c.GetId(),
			Secret:      c.GetSecret(),
			RedirectUri: c.GetRedirectUri(),
			UserData:    c.GetUserData(),
		}, nil
	}
	return nil, osin.ErrNotFound
}

func (s *MemoryStorage) SetClient(id string, client osin.Client) error {
	logger.Debugf("SetClient: %s\n", id)

	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	s.clients[id] = osin.DefaultClient{
		Id:          client.GetId(),
		Secret:      client.GetSecret(),
		RedirectUri: client.GetRedirectUri(),
		UserData:    client.GetUserData(),
	}
	return nil
}

func (s *MemoryStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	logger.Debugf("SaveAuthorize: %s\n", data.Code)

	s.authorizeMutex.Lock()
	defer s.authorizeMutex.Unlock()

	s.authorize[data.Code] = *data
	return nil
}

func (s *MemoryStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	logger.Debugf("LoadAuthorize: %s\n", code)

	s.authorizeMutex.Lock()
	defer s.authorizeMutex.Unlock()

	if d, ok := s.authorize[code]; ok {
		return &d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *MemoryStorage) RemoveAuthorize(code string) error {
	logger.Debugf("RemoveAuthorize: %s\n", code)

	s.authorizeMutex.Lock()
	defer s.authorizeMutex.Unlock()

	delete(s.authorize, code)
	return nil
}

func (s *MemoryStorage) SaveAccess(data *osin.AccessData) error {
	logger.Debugf("SaveAccess: %s\n", data.AccessToken)

	s.accessMutex.Lock()
	defer s.accessMutex.Unlock()

	s.access[data.AccessToken] = *data
	if data.RefreshToken != "" {
		s.refreshMutex.Lock()
		defer s.refreshMutex.Unlock()
		s.refresh[data.RefreshToken] = data.AccessToken
	}
	return nil
}

func (s *MemoryStorage) LoadAccess(code string) (*osin.AccessData, error) {
	logger.Debugf("LoadAccess: %s\n", code)

	s.accessMutex.Lock()
	defer s.accessMutex.Unlock()

	if d, ok := s.access[code]; ok {
		return &d, nil
	}
	return nil, osin.ErrNotFound
}

func (s *MemoryStorage) RemoveAccess(code string) error {
	logger.Debugf("RemoveAccess: %s\n", code)

	s.accessMutex.Lock()
	defer s.accessMutex.Unlock()

	delete(s.access, code)
	return nil
}

func (s *MemoryStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	logger.Debugf("LoadRefresh: %s\n", code)

	s.refreshMutex.Lock()
	defer s.refreshMutex.Unlock()

	if d, ok := s.refresh[code]; ok {
		return s.LoadAccess(d)
	}
	return nil, osin.ErrNotFound
}

func (s *MemoryStorage) RemoveRefresh(code string) error {
	logger.Debugf("RemoveRefresh: %s\n", code)

	s.refreshMutex.Lock()
	defer s.refreshMutex.Unlock()

	delete(s.refresh, code)
	return nil
}
