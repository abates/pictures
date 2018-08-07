package pictures

import (
	"io"
	"os"
	"path"
)

type Server struct {
	fs          Filesystem
	db          DB
	disgoDB     *DisgoDB
	ingestChain *ProcessingChain
}

var DefaultPerm = os.FileMode(0750)

func DefaultServer(outputPath string) (server *Server, err error) {
	ingestChain := NewProcessingChain()
	filesDir := path.Join(outputPath, "files")
	err = os.MkdirAll(filesDir, DefaultPerm)

	var fs Filesystem
	if err == nil {
		fs, err = OpenOSFilesystem(filesDir, DefaultPerm)
	}

	dbPath := path.Join(outputPath, "db")
	if err == nil {
		err = os.MkdirAll(dbPath, DefaultPerm)
	}

	var db DB
	if err == nil {
		db, err = OpenBadger(dbPath)
	}

	var disgoDB *DisgoDB
	if err == nil {
		disgoDB, err = OpenDisgoDB(db)
	}

	if err == nil {
		ingestChain.
			Append(&ImageDecoderFilter{}).
			Append(&ExifFilter{}).
			Append(NewDisgoFilter(disgoDB)).
			Append(NewPathFilter(fs)).
			AppendLast(NewOutputFilter(fs))
		server = New(fs, ingestChain, db, disgoDB)
	}

	return
}

func New(fs Filesystem, ingestChain *ProcessingChain, db DB, disgoDB *DisgoDB) *Server {
	return &Server{
		fs:          fs,
		ingestChain: ingestChain,
		db:          db,
		disgoDB:     disgoDB,
	}
}

func (server *Server) SetDebug(debug bool) {
	server.ingestChain.debug = debug
}

func (server *Server) Ingest(buf []byte) error {
	info := &ImageInfo{Buf: buf}
	return server.ingestChain.Process(info)
}

func (server *Server) List(path string) ([]FileInfo, error) {
	return server.fs.List(path)
}

func (server *Server) Open(path string) (io.Reader, error) {
	return server.fs.Open(path)
}

func (server *Server) ReadFile(path string) ([]byte, error) {
	return server.fs.ReadFile(path)
}

func (server *Server) Close() (err error) {
	if closer, ok := server.db.(io.Closer); ok {
		err = closer.Close()
	}

	if err == nil && server.disgoDB != nil {
		err = server.disgoDB.Close()
	}
	return err
}
