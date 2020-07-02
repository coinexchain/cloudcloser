package cloudlocker

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

type LockerServer struct {
	db     *leveldb.DB
	server *http.Server
}

func NewLockerServer(path, url string) (*LockerServer, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	l := &LockerServer{
		db: db,
	}
	l.server = &http.Server{
		Addr:    url,
		Handler: NewRouter(l),
	}
	return l, nil
}

func (l *LockerServer) Start() {
	if err := l.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func (l *LockerServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := l.server.Shutdown(ctx); err != nil {
		//log
	}
	err := l.db.Close()
	if err != nil {
		//log
	}
}
