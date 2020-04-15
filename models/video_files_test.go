// Code generated by SQLBoiler 3.6.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testVideoFiles(t *testing.T) {
	t.Parallel()

	query := VideoFiles()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testVideoFilesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVideoFilesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := VideoFiles().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVideoFilesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := VideoFileSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testVideoFilesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := VideoFileExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if VideoFile exists: %s", err)
	}
	if !e {
		t.Errorf("Expected VideoFileExists to return true, but got false.")
	}
}

func testVideoFilesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	videoFileFound, err := FindVideoFile(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if videoFileFound == nil {
		t.Error("want a record, got nil")
	}
}

func testVideoFilesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = VideoFiles().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testVideoFilesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := VideoFiles().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testVideoFilesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	videoFileOne := &VideoFile{}
	videoFileTwo := &VideoFile{}
	if err = randomize.Struct(seed, videoFileOne, videoFileDBTypes, false, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}
	if err = randomize.Struct(seed, videoFileTwo, videoFileDBTypes, false, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = videoFileOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = videoFileTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := VideoFiles().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testVideoFilesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	videoFileOne := &VideoFile{}
	videoFileTwo := &VideoFile{}
	if err = randomize.Struct(seed, videoFileOne, videoFileDBTypes, false, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}
	if err = randomize.Struct(seed, videoFileTwo, videoFileDBTypes, false, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = videoFileOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = videoFileTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func videoFileBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func videoFileAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *VideoFile) error {
	*o = VideoFile{}
	return nil
}

func testVideoFilesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &VideoFile{}
	o := &VideoFile{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, videoFileDBTypes, false); err != nil {
		t.Errorf("Unable to randomize VideoFile object: %s", err)
	}

	AddVideoFileHook(boil.BeforeInsertHook, videoFileBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	videoFileBeforeInsertHooks = []VideoFileHook{}

	AddVideoFileHook(boil.AfterInsertHook, videoFileAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	videoFileAfterInsertHooks = []VideoFileHook{}

	AddVideoFileHook(boil.AfterSelectHook, videoFileAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	videoFileAfterSelectHooks = []VideoFileHook{}

	AddVideoFileHook(boil.BeforeUpdateHook, videoFileBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	videoFileBeforeUpdateHooks = []VideoFileHook{}

	AddVideoFileHook(boil.AfterUpdateHook, videoFileAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	videoFileAfterUpdateHooks = []VideoFileHook{}

	AddVideoFileHook(boil.BeforeDeleteHook, videoFileBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	videoFileBeforeDeleteHooks = []VideoFileHook{}

	AddVideoFileHook(boil.AfterDeleteHook, videoFileAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	videoFileAfterDeleteHooks = []VideoFileHook{}

	AddVideoFileHook(boil.BeforeUpsertHook, videoFileBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	videoFileBeforeUpsertHooks = []VideoFileHook{}

	AddVideoFileHook(boil.AfterUpsertHook, videoFileAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	videoFileAfterUpsertHooks = []VideoFileHook{}
}

func testVideoFilesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVideoFilesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(videoFileColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testVideoFileToOneVideoUsingVideo(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local VideoFile
	var foreign Video

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, videoFileDBTypes, false, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, videoDBTypes, false, videoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Video struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.VideoID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Video().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := VideoFileSlice{&local}
	if err = local.L.LoadVideo(ctx, tx, false, (*[]*VideoFile)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Video == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Video = nil
	if err = local.L.LoadVideo(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Video == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testVideoFileToOneSetOpVideoUsingVideo(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a VideoFile
	var b, c Video

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, videoFileDBTypes, false, strmangle.SetComplement(videoFilePrimaryKeyColumns, videoFileColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, videoDBTypes, false, strmangle.SetComplement(videoPrimaryKeyColumns, videoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, videoDBTypes, false, strmangle.SetComplement(videoPrimaryKeyColumns, videoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Video{&b, &c} {
		err = a.SetVideo(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Video != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.VideoFiles[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.VideoID != x.ID {
			t.Error("foreign key was wrong value", a.VideoID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.VideoID))
		reflect.Indirect(reflect.ValueOf(&a.VideoID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.VideoID != x.ID {
			t.Error("foreign key was wrong value", a.VideoID, x.ID)
		}
	}
}

func testVideoFilesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testVideoFilesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := VideoFileSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testVideoFilesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := VideoFiles().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	videoFileDBTypes = map[string]string{`VideoID`: `integer`, `VideoFileTypeName`: `character varying`, `Filename`: `character varying`, `IsEnabled`: `boolean`, `Comments`: `text`, `ID`: `integer`, `Size`: `bigint`}
	_                = bytes.MinRead
)

func testVideoFilesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(videoFilePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(videoFileAllColumns) == len(videoFilePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFilePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testVideoFilesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(videoFileAllColumns) == len(videoFilePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &VideoFile{}
	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFileColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, videoFileDBTypes, true, videoFilePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(videoFileAllColumns, videoFilePrimaryKeyColumns) {
		fields = videoFileAllColumns
	} else {
		fields = strmangle.SetComplement(
			videoFileAllColumns,
			videoFilePrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := VideoFileSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testVideoFilesUpsert(t *testing.T) {
	t.Parallel()

	if len(videoFileAllColumns) == len(videoFilePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := VideoFile{}
	if err = randomize.Struct(seed, &o, videoFileDBTypes, true); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert VideoFile: %s", err)
	}

	count, err := VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, videoFileDBTypes, false, videoFilePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize VideoFile struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert VideoFile: %s", err)
	}

	count, err = VideoFiles().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
