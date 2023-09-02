package metadata

import (
	"context"
	"errors"
	"testing"

	gen "github.com/ahmad-khatib0/go/microservice/movies/gen/mock/metadata/repository"
	"github.com/ahmad-khatib0/go/microservice/movies/metadata/internal/repository"
	"github.com/ahmad-khatib0/go/microservice/movies/metadata/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	tests := []struct {
		name       string
		expRepoRes *model.Metadata
		expRepoErr error
		wantRes    *model.Metadata
		wantErr    error
	}{
		{
			name:       "not found",
			expRepoErr: repository.ErrNotFound,
			wantErr:    ErrNotFound,
		},
		{
			name:       "unexpected error",
			expRepoErr: errors.New("unexpected error"),
			wantErr:    errors.New("unexpected error"),
		},
		{
			name:       "success",
			expRepoRes: &model.Metadata{},
			wantRes:    &model.Metadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := gen.NewMockmetadataRepository(ctrl)
			c := New(repoMock)
			ctx := context.Background()
			id := "id"
			repoMock.EXPECT().Get(ctx, id).Return(tt.expRepoRes, tt.expRepoErr)
			res, err := c.Get(ctx, id)
			assert.Equal(t, tt.wantRes, res, tt.name)
			assert.Equal(t, tt.wantErr, err, tt.name)
		})
	}
}
