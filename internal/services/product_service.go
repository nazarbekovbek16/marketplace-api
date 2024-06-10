// services/product_service.go

package services

import (
	"marketplace-api/internal/models"
	"marketplace-api/internal/repository"
)

type ProductService struct {
	productRepository     *repository.ProductRepository
	distributorRepository *repository.DistributorRepository
}

func NewProductService(productRepository *repository.ProductRepository, distributorRepository *repository.DistributorRepository) *ProductService {
	return &ProductService{productRepository: productRepository, distributorRepository: distributorRepository}
}

func (ps *ProductService) CreateProduct(product *models.Product) error {
	return ps.productRepository.CreateProduct(product)
}

func (ps *ProductService) UpdateProduct(product *models.Product) error {
	return ps.productRepository.UpdateProduct(product)
}

func (ps *ProductService) DeleteProduct(productID int64) error {
	return ps.productRepository.DeleteProduct(productID)
}

func (ps *ProductService) GetProductByID(productID int64) (*models.Product, error) {
	product, err := ps.productRepository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}
	distributor, err := ps.distributorRepository.GetDistributorByID(product.DistributorID)
	if err != nil {
		return nil, err
	}
	product.Distributor = *distributor
	return product, nil
}

func (ps *ProductService) GetProductsByDistributorID(productName string, filters models.Filters, distributorID int64) ([]*models.Product, models.Metadata, error) {
	return ps.productRepository.GetProductsByDistributorID(productName, filters, distributorID)
}
func (ps *ProductService) GetProducts(productName string, filters models.Filters) ([]*models.Product, models.Metadata, error) {
	return ps.productRepository.GetProducts(productName, filters)
}

func (ps *ProductService) CreatReview(review *models.Review) error {
	return ps.productRepository.CreatReview(review)
}
func (ps *ProductService) GetReviewByStoreId(storeId int64) ([]models.Review, error) {
	return ps.productRepository.GetReviews(storeId, "store")
}

func (ps *ProductService) GetReviewsByDistributorId(distributorId int64) ([]models.Review, error) {
	return ps.productRepository.GetReviews(distributorId, "distributor")
}

func (ps *ProductService) GetReviewsByProductId(productID int64) ([]models.Review, error) {
	return ps.productRepository.GetReviews(productID, "product")
}
func (ps *ProductService) DeleteByReviewId(reviewId int64) error {
	return ps.productRepository.DeleteByReviewId(reviewId)
}

// Implement filtering, sorting, and metadata functions as needed
