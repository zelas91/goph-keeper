package repository

import (
	"fmt"

	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"golang.org/x/net/context"
)

type binaryFile struct {
	tm transactionManager
}

func (b binaryFile) Create(ctx context.Context, bf entities.BinaryFile) error {
	query := `insert into binary_file (path, file_name, user_id, size)
		values (:path,:file_name,:user_id, :size);`
	if _, err := b.tm.getConn(ctx).NamedExecContext(ctx, query, bf); err != nil {
		return fmt.Errorf("repo binary file create err: %v", err)
	}
	return nil
}

func (b binaryFile) FindByIDAndUserID(ctx context.Context, fileID, userID int) (entities.BinaryFile, error) {
	query := `select * from binary_file where id=$1 and user_id=$2`
	var bf entities.BinaryFile
	if err := b.tm.getConn(ctx).GetContext(ctx, &bf, query, fileID, userID); err != nil {
		return bf, fmt.Errorf("repo: binary file get id=%d  err: %v", fileID, err)
	}
	return bf, nil
}

func (b binaryFile) FindAllByUserID(ctx context.Context, userID int) ([]entities.BinaryFile, error) {
	query := `select * from binary_file where user_id=$1`
	var files []entities.BinaryFile
	if err := b.tm.getConn(ctx).SelectContext(ctx, &files, query, userID); err != nil {
		return files, fmt.Errorf("repo: get binary files err %v", err)
	}
	return files, nil
}

func (b binaryFile) Delete(ctx context.Context, userID, fileID int) error {
	query := `delete from binary_file where id=$1 and user_id=$2`
	if _, err := b.tm.getConn(ctx).ExecContext(ctx, query, fileID, userID); err != nil {
		return fmt.Errorf("repo binary file delete err: %v", err)
	}
	return nil
}
