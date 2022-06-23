package resource

import (
	"errors"
	"testing"

	"github.com/ans-group/cli/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestResourceLocator_Invoke(t *testing.T) {
	t.Run("SingleProperty_RetrievesItem", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		provider := mocks.NewMockResourceLocatorProvider(mockCtrl)

		provider.EXPECT().SupportedProperties().Return([]string{"testproperty1"}).Times(1)
		provider.EXPECT().Locate("testproperty1", "testvalue1").Return([]string{"testlocateresult1"}, nil)

		r := NewResourceLocator(provider)

		result, err := r.Invoke("testvalue1")

		assert.Nil(t, err)
		assert.Equal(t, "testlocateresult1", result.(string))
	})

	t.Run("MultipleProperties_RetrievesItem", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		provider := mocks.NewMockResourceLocatorProvider(mockCtrl)

		provider.EXPECT().SupportedProperties().Return([]string{"testproperty1", "testproperty2", "testproperty3"}).Times(1)
		provider.EXPECT().Locate("testproperty1", "testvalue1").Return(nil, nil)
		provider.EXPECT().Locate("testproperty2", "testvalue1").Return([]string{"testlocateresult1"}, nil)

		r := NewResourceLocator(provider)

		result, err := r.Invoke("testvalue1")

		assert.Nil(t, err)
		assert.Equal(t, "testlocateresult1", result.(string))
	})

	t.Run("ProviderLocateError_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		provider := mocks.NewMockResourceLocatorProvider(mockCtrl)

		provider.EXPECT().SupportedProperties().Return([]string{"testproperty1"}).Times(1)
		provider.EXPECT().Locate("testproperty1", "testvalue1").Return(nil, errors.New("test error 1"))

		r := NewResourceLocator(provider)

		_, err := r.Invoke("testvalue1")

		assert.NotNil(t, err)
		assert.Equal(t, "Error retrieving items: test error 1", err.Error())
	})

	t.Run("MultipleItems_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		provider := mocks.NewMockResourceLocatorProvider(mockCtrl)

		provider.EXPECT().SupportedProperties().Return([]string{"testproperty1"}).Times(1)
		provider.EXPECT().Locate("testproperty1", "testvalue1").Return([]string{"testlocateresult1", "testlocateresult2"}, nil)

		r := NewResourceLocator(provider)

		_, err := r.Invoke("testvalue1")

		assert.NotNil(t, err)
		assert.Equal(t, "More than one item found matching [testvalue1] (testproperty1)", err.Error())
	})

	t.Run("NoItems_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		provider := mocks.NewMockResourceLocatorProvider(mockCtrl)

		provider.EXPECT().SupportedProperties().Return([]string{"testproperty1"}).Times(1)
		provider.EXPECT().Locate("testproperty1", "testvalue1").Return(nil, nil)

		r := NewResourceLocator(provider)

		_, err := r.Invoke("testvalue1")

		assert.NotNil(t, err)
		assert.Equal(t, "No items found matching [testvalue1]", err.Error())
	})

	t.Run("LocateReturnsNoneSlice_ReturnsError", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		provider := mocks.NewMockResourceLocatorProvider(mockCtrl)

		provider.EXPECT().SupportedProperties().Return([]string{"testproperty1"}).Times(1)
		provider.EXPECT().Locate("testproperty1", "testvalue1").Return("non-slice", nil)

		r := NewResourceLocator(provider)

		_, err := r.Invoke("testvalue1")

		assert.NotNil(t, err)
		assert.Equal(t, "Unsupported non-slice type [string]", err.Error())
	})
}
