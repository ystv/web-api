// Code generated by SQLBoiler 3.6.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxes)
	t.Run("VideoFiles", testVideoFiles)
	t.Run("Videos", testVideos)
}

func TestDelete(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesDelete)
	t.Run("VideoFiles", testVideoFilesDelete)
	t.Run("Videos", testVideosDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesQueryDeleteAll)
	t.Run("VideoFiles", testVideoFilesQueryDeleteAll)
	t.Run("Videos", testVideosQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesSliceDeleteAll)
	t.Run("VideoFiles", testVideoFilesSliceDeleteAll)
	t.Run("Videos", testVideosSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesExists)
	t.Run("VideoFiles", testVideoFilesExists)
	t.Run("Videos", testVideosExists)
}

func TestFind(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesFind)
	t.Run("VideoFiles", testVideoFilesFind)
	t.Run("Videos", testVideosFind)
}

func TestBind(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesBind)
	t.Run("VideoFiles", testVideoFilesBind)
	t.Run("Videos", testVideosBind)
}

func TestOne(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesOne)
	t.Run("VideoFiles", testVideoFilesOne)
	t.Run("Videos", testVideosOne)
}

func TestAll(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesAll)
	t.Run("VideoFiles", testVideoFilesAll)
	t.Run("Videos", testVideosAll)
}

func TestCount(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesCount)
	t.Run("VideoFiles", testVideoFilesCount)
	t.Run("Videos", testVideosCount)
}

func TestHooks(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesHooks)
	t.Run("VideoFiles", testVideoFilesHooks)
	t.Run("Videos", testVideosHooks)
}

func TestInsert(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesInsert)
	t.Run("VideoBoxes", testVideoBoxesInsertWhitelist)
	t.Run("VideoFiles", testVideoFilesInsert)
	t.Run("VideoFiles", testVideoFilesInsertWhitelist)
	t.Run("Videos", testVideosInsert)
	t.Run("Videos", testVideosInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("VideoFileToVideoUsingVideo", testVideoFileToOneVideoUsingVideo)
	t.Run("VideoToVideoUsingRedirect", testVideoToOneVideoUsingRedirect)
	t.Run("VideoToVideoBoxUsingVideoBox", testVideoToOneVideoBoxUsingVideoBox)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("VideoBoxToVideos", testVideoBoxToManyVideos)
	t.Run("VideoToVideoFiles", testVideoToManyVideoFiles)
	t.Run("VideoToRedirectVideos", testVideoToManyRedirectVideos)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("VideoFileToVideoUsingVideoFiles", testVideoFileToOneSetOpVideoUsingVideo)
	t.Run("VideoToVideoUsingRedirectVideos", testVideoToOneSetOpVideoUsingRedirect)
	t.Run("VideoToVideoBoxUsingVideos", testVideoToOneSetOpVideoBoxUsingVideoBox)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("VideoToVideoUsingRedirectVideos", testVideoToOneRemoveOpVideoUsingRedirect)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("VideoBoxToVideos", testVideoBoxToManyAddOpVideos)
	t.Run("VideoToVideoFiles", testVideoToManyAddOpVideoFiles)
	t.Run("VideoToRedirectVideos", testVideoToManyAddOpRedirectVideos)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("VideoToRedirectVideos", testVideoToManySetOpRedirectVideos)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("VideoToRedirectVideos", testVideoToManyRemoveOpRedirectVideos)
}

func TestReload(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesReload)
	t.Run("VideoFiles", testVideoFilesReload)
	t.Run("Videos", testVideosReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesReloadAll)
	t.Run("VideoFiles", testVideoFilesReloadAll)
	t.Run("Videos", testVideosReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesSelect)
	t.Run("VideoFiles", testVideoFilesSelect)
	t.Run("Videos", testVideosSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesUpdate)
	t.Run("VideoFiles", testVideoFilesUpdate)
	t.Run("Videos", testVideosUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("VideoBoxes", testVideoBoxesSliceUpdateAll)
	t.Run("VideoFiles", testVideoFilesSliceUpdateAll)
	t.Run("Videos", testVideosSliceUpdateAll)
}
