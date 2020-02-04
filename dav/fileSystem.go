package dav

import (
	"context"
	"golang.org/x/net/webdav"
	"os"
	"time"
)

type PhantomFileSystem struct {
}

func (pfs *PhantomFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return nil
}

func (pfs *PhantomFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return nil, nil
}

func (pfs *PhantomFileSystem) RemoveAll(ctx context.Context, name string) error {
	return nil
}

func (pfs *PhantomFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return nil
}

func (pfs *PhantomFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	return nil, nil
}

type PhantomLockSystem struct {
}

func (pls *PhantomLockSystem) Confirm(now time.Time, name0, name1 string, conditions ...webdav.Condition) (release func(), err error) {
	return func() {

	}, err
}

func (pls *PhantomLockSystem) Create(now time.Time, details webdav.LockDetails) (token string, err error) {
	return "", nil
}
func (pls *PhantomLockSystem) Refresh(now time.Time, token string, duration time.Duration) (webdav.LockDetails, error) {
	return webdav.LockDetails{}, nil
}
func (pls *PhantomLockSystem) Unlock(now time.Time, token string) error {
	return nil
}
