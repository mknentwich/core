package dav

import (
	"context"
	"fmt"
	"golang.org/x/net/webdav"
	"os"
	"time"
)

const separator = "/"

type PhantomFileSystem struct {
}

func (pfs *PhantomFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return nil
}

func (pfs *PhantomFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	node := tree.Subset(name)
	if node == nil {
		return nil, fmt.Errorf("node %s is not available", name)
	}
	file := node.File()
	if file == nil {
		return nil, fmt.Errorf("the node %s in %s has no corresponding file", node.Name(), name)
	}
	return file, nil
}

func (pfs *PhantomFileSystem) RemoveAll(ctx context.Context, name string) error {
	return nil
}

func (pfs *PhantomFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return nil
}

func (pfs *PhantomFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	node := tree.Subset(name)
	if node == nil {
		return BasicStat{}, fmt.Errorf("node %s is not available", name)
	}
	file := node.File()
	if file == nil {
		return BasicStat{}, fmt.Errorf("the node %s in %s has no corresponding file", node.Name(), name)
	}
	return file.Stat()
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

type PhantomInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	dir     bool
}

func (i PhantomInfo) Name() string {
	return i.name
}

func (i PhantomInfo) Size() int64 {
	return i.size
}

func (i PhantomInfo) Mode() os.FileMode {
	return i.mode
}

func (i PhantomInfo) ModTime() time.Time {
	return i.modTime
}

func (i PhantomInfo) IsDir() bool {
	return i.dir
}

func (i PhantomInfo) Sys() interface{} {
	panic("implement me")
}

type PhantomFile struct {
	os.File
}

func (p *PhantomFile) Close() error {
	panic("implement me")
}

func (p *PhantomFile) Read(data []byte) (n int, err error) {
	panic("implement me")
}

func (p *PhantomFile) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (p *PhantomFile) Readdir(count int) ([]os.FileInfo, error) {
	panic("implement me")
}

func (p *PhantomFile) Stat() (os.FileInfo, error) {
	panic("implement me")
}

func (p *PhantomFile) Write(data []byte) (n int, err error) {
	panic("implement me")
}
