package listeners

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteTorrentCreatesDirectoryAndCompleteFile(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "torrent_files")
	if err := writeTorrentFile(dir, 101, []byte("synthetic test bytes")); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "101.torrent"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "synthetic test bytes" {
		t.Fatalf("file=%q", b)
	}
	info, err := os.Stat(filepath.Join(dir, "101.torrent"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("watch consumer cannot read file: mode=%o", info.Mode().Perm())
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("partial files remain: %v", entries)
	}
}
