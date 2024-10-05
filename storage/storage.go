//go:build !js

package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/m110/kingdoms/save"
)

var marshalOptions = proto.MarshalOptions{
	AllowPartial: true,
}

var unmarshalOptions = proto.UnmarshalOptions{
	AllowPartial:   true,
	DiscardUnknown: true,
}

type Storage struct {
	saveDir string
}

func NewStorage() (*Storage, error) {
	dir, err := storageDirectory()
	if err != nil {
		return nil, err
	}

	return &Storage{
		saveDir: dir,
	}, nil
}

func (s *Storage) SaveGlobalData(global *save.Global) error {
	return s.save(global, s.globalSavePath())
}

func (s *Storage) OccupiedSlots() ([]int, error) {
	dir, err := os.ReadDir(s.gameSavesPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var slots []int
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		if path.Ext(entry.Name()) != ".save" {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".save")
		slot, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		slots = append(slots, slot)
	}

	return slots, nil
}

func (s *Storage) LoadSlot(slot int) (*save.Game, error) {
	game := &save.Game{}
	err := s.load(game, s.gameSavePath(slot))
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (s *Storage) LoadReplay(slot int) (*save.GameReplay, error) {
	replay := &save.GameReplay{}
	err := s.load(replay, s.gameReplayPath(slot))
	if err != nil {
		return nil, err
	}
	return replay, nil
}

func (s *Storage) SaveSlot(slot int, global *save.Global, game *save.Game) error {
	err := s.SaveGlobalData(global)
	if err != nil {
		return err
	}

	return s.save(game, s.gameSavePath(slot))
}

func (s *Storage) DeleteSlot(slot int) error {
	err := s.delete(s.gameSavePath(slot))
	if err != nil {
		return err
	}

	err = s.delete(s.gameReplayPath(slot))
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) SaveReplay(slot int, replay *save.GameReplay) error {
	return s.save(replay, s.gameReplayPath(slot))
}

func (s *Storage) LoadGlobalData() (*save.Global, error) {
	global := &save.Global{}
	err := s.load(global, s.globalSavePath())
	if err != nil {
		return nil, err
	}
	return global, nil
}

func (s *Storage) save(data proto.Message, savePath string) error {
	err := os.MkdirAll(path.Dir(savePath), 0o755)
	if err != nil {
		return err
	}

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	bytes, err := marshalOptions.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) load(data proto.Message, path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = unmarshalOptions.Unmarshal(bytes, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) delete(path string) error {
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func (s *Storage) globalSavePath() string {
	return path.Join(s.saveDir, "global.save")
}

func (s *Storage) gameSavesPath() string {
	return path.Join(s.saveDir, "games")
}

func (s *Storage) gameReplaysPath() string {
	return path.Join(s.saveDir, "replays")
}

func (s *Storage) gameSavePath(slot int) string {
	return path.Join(s.gameSavesPath(), fmt.Sprintf("%d.save", slot))
}

func (s *Storage) gameReplayPath(slot int) string {
	return path.Join(s.gameReplaysPath(), fmt.Sprintf("%d.replay", slot))
}
