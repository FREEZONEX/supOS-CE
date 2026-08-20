package repo

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const initialNamespaceNodeVersionBase int64 = 10000000

type UnsNamespaceNodeVersion struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Version  int64  `gorm:"column:version;type:bigint;index" json:"version"`
	Etag     string `gorm:"column:etag;size:64;index" json:"etag"`
	Checksum string `gorm:"column:checksum;size:64" json:"checksum"`
	UnsTree  any    `gorm:"column:unstree;type:jsonb;serializer:json" json:"unsTree,omitempty"`
	NoDelTime
}

func (UnsNamespaceNodeVersion) TableName() string { return "uns_namespace_node_version" }

type UnsNamespaceNodeVersionRepo struct{ db *gorm.DB }

func NewUnsNamespaceNodeVersionRepo(in any) *UnsNamespaceNodeVersionRepo {
	return &UnsNamespaceNodeVersionRepo{db: GetCommonConn(in)}
}

func (r *UnsNamespaceNodeVersionRepo) Current(ctx context.Context) (int64, error) {
	var row UnsNamespaceNodeVersion
	err := r.db.WithContext(ctx).Order("version desc").First(&row).Error
	if err == nil {
		return row.Version, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	return 0, nil
}

func (r *UnsNamespaceNodeVersionRepo) Bump(ctx context.Context, trees ...any) (int64, error) {
	var tree any
	if len(trees) > 0 {
		tree = trees[0]
	}
	var next int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var latest UnsNamespaceNodeVersion
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Order("version desc").
			First(&latest).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			next = initialNamespaceNodeVersionBase
		} else if err != nil {
			return err
		} else {
			next = latest.Version + 1
		}
		checksum, etag, err := computeNamespaceNodeVersionMeta(next, tree)
		if err != nil {
			return err
		}
		return tx.Create(&UnsNamespaceNodeVersion{
			Version:  next,
			Etag:     etag,
			Checksum: checksum,
			UnsTree:  tree,
		}).Error
	})
	return next, err
}

func computeNamespaceNodeVersionMeta(nextVersion int64, tree any) (string, string, error) {
	body, err := json.Marshal(tree)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256(body)
	checksum := fmt.Sprintf("%x", sum[:])
	etagSuffix := checksum
	if len(etagSuffix) > 8 {
		etagSuffix = etagSuffix[:8]
	}
	return checksum, fmt.Sprintf("unsrev-%d-%s", nextVersion, etagSuffix), nil
}
