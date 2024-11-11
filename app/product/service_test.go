package product_test

import (
	"context"
	"errors"
	"testing"
	"training/app/product"
	"training/persistence"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestServiceDelete(t *testing.T) {
	tests := []struct {
		name          string
		repositoryErr error
		productId     uuid.UUID
		expectedErr   error
	}{
		{
			"1",
			nil,
			uuid.New(),
			nil,
		}, {
			"2",
			errors.New("repository error"),
			uuid.New(),
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := product.NewMockRepository(ctrl)
			repositoryMock.EXPECT().Delete(ctx, gomock.Eq(tt.productId)).Return(tt.repositoryErr).Times(1)

			s := product.NewService(repositoryMock)
			err := s.Delete(ctx, tt.productId)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceGetById(t *testing.T) {
	tests := []struct {
		name            string
		product         *persistence.Product
		repositoryErr   error
		productId       uuid.UUID
		expectedProduct *product.GetProductDto
		expectedErr     error
	}{
		{
			"1",
			&persistence.Product{},
			nil,
			uuid.New(),
			&product.GetProductDto{},
			nil,
		}, {
			"2",
			&persistence.Product{},
			errors.New("repository error"),
			uuid.New(),
			nil,
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := product.NewMockRepository(ctrl)

			repositoryMock.EXPECT().SelectById(ctx, gomock.Eq(tt.productId)).Return(tt.product, tt.repositoryErr).Times(1)

			s := product.NewService(repositoryMock)
			got, err := s.GetById(ctx, tt.productId)

			assert.Equal(t, tt.expectedProduct, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceSave(t *testing.T) {
	tests := []struct {
		name            string
		product         *persistence.Product
		repositoryErr   error
		payload         product.SaveProductPayload
		expectedProduct *product.SaveProductDto
		expectedErr     error
	}{
		{"1",
			&persistence.Product{},
			nil, product.SaveProductPayload{},
			&product.SaveProductDto{},
			nil,
		}, {
			"2",
			nil,
			errors.New("repository error"),
			product.SaveProductPayload{},
			nil,
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := product.NewMockRepository(ctrl)

			repositoryMock.EXPECT().Insert(ctx, gomock.Any()).Return(tt.product, tt.repositoryErr).Times(1)

			s := product.NewService(repositoryMock)
			got, err := s.Save(ctx, tt.payload)

			assert.Equal(t, tt.expectedProduct, got)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	tests := []struct {
		name          string
		repositoryErr error
		payload       product.UpdateProductPayload
		expectedErr   error
	}{
		{
			"1",
			nil,
			product.UpdateProductPayload{},
			nil,
		}, {
			"2",
			errors.New("repository error"),
			product.UpdateProductPayload{},
			errors.New("repository error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			repositoryMock := product.NewMockRepository(ctrl)

			repositoryMock.EXPECT().Update(ctx, gomock.Any()).Return(tt.repositoryErr).Times(1)

			s := product.NewService(repositoryMock)
			err := s.Update(ctx, tt.payload)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
